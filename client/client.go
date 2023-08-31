package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/TheLeeeo/gql-test-suite/models"
	"github.com/TheLeeeo/gql-test-suite/request"
	"github.com/TheLeeeo/gql-test-suite/utils"
)

// The levels of type nesting to begin fetching
const defaultTypeDepth = 4

// Builds the recursive "ofType" field for the type introspection query
// Used to fetch the entire type tree
func buildRecursiveOfTypeField(depth int) string {
	if depth == 0 {
		return ""
	}

	return fmt.Sprintf(recursiveOfTypeField, buildRecursiveOfTypeField(depth-1))
}

type Client struct {
	Addr      string
	Types     map[string]*models.Type
	Queries   []*models.Field
	Mutations []*models.Field
}

func New(addr string) *Client {
	return &Client{
		Addr:  addr,
		Types: make(map[string]*models.Type),
	}
}

func (c *Client) GetType(name string) *models.Type {
	t, ok := c.Types[name]
	if !ok {
		// Type not found in client, fetching
		t, err := c.FetchType(name)
		if err != nil {
			log.Printf("error fetching type %s: %v", name, err)
		}
		return t
	}

	return t
}

// Executes the graphql request and returns the response
func (c *Client) Execute(request []byte) ([]byte, error) {
	requestBody := bytes.NewBuffer(request)

	req, err := http.NewRequest("POST", c.Addr, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)

	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "test")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 { //not neccesarly an error
		log.Printf("statuscode (%d) making request:\n\t%s\nresponse:\n\t%s", resp.StatusCode, utils.PrettyRequest(string(request)), string(responseBody))
	}

	return responseBody, nil
}

func (c *Client) FetchType(typeName string) (*models.Type, error) {
	if typeName == "_Service" || typeName == "_Any" || typeName == "_FieldSet" {
		return nil, fmt.Errorf("type %s is not supported, skipping", typeName)
	}

	t, err := c.fetchTypeInternal(typeName, defaultTypeDepth)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// The internal function for fetching a type. Deals with incomplete types
func (c *Client) fetchTypeInternal(typeName string, typeDepth int) (*models.Type, error) {
	ofTypeField := buildRecursiveOfTypeField(typeDepth)
	reqString := fmt.Sprintf(typeIntrospectionQuery, typeName, ofTypeField)

	q := request.BuildFromString(reqString, nil)
	resp, err := c.Execute([]byte(q))
	if err != nil {
		return nil, err
	}

	respType, err := utils.ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	dataMap := respType.Data["__type"].(map[string]any)

	t := &models.Type{}
	err = utils.ParseMap(dataMap, t)
	if err != nil {
		return nil, err
	}

	if isCompleteType(t) {
		return t, nil
	}

	if typeDepth == 0 {
		typeDepth = 1
	}

	return c.fetchTypeInternal(typeName, typeDepth*2)
}

// Checks if a type is comleted.
// It is considered complete if it knows the base typ of all of its fields
func isCompleteType(t *models.Type) bool { //TODO: check inputFIelds and args
	baseType := t.GetBaseType()
	if baseType.Kind == models.NonNullTypeKind || baseType.Kind == models.ListTypeKind {
		return false
	}

	for _, f := range t.Fields {
		if !isCompleteField(&f) {
			return false
		}
	}

	return true
}

// Checks if a field has all of its types complted
func isCompleteField(f *models.Field) bool {
	if !isCompleteType(f.Type) {
		return false
	}

	for _, a := range f.Args {
		if !isCompleteType(a.Type) {
			return false
		}
	}

	return true
}

func (c *Client) LoadSchema() error {
	q := fmt.Sprintf(schemaIntrospectionQuery, buildRecursiveOfTypeField(defaultTypeDepth))
	req := request.BuildFromString(q, nil)
	resp, err := c.Execute([]byte(req))
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}

	respType, err := utils.ParseResponse(resp)
	if err != nil {
		return err
	}
	dataMap := respType.Data["__schema"].(map[string]any)

	sch := &models.Schema{}
	err = utils.ParseMap(dataMap, sch)
	if err != nil {
		return fmt.Errorf("error parsing schema: %v", err)
	}

	for _, t := range sch.Types {
		if isCompleteType(t) {
			c.Types[t.Name] = t
		} else {
			t, err := c.FetchType(t.Name)
			if err != nil {
				return fmt.Errorf("error fetching type %s: %v", t.Name, err)
			}
			c.Types[t.Name] = t
		}
	}

	queries, ok := c.Types["Query"]
	if !ok || len(queries.Fields) == 0 {
		log.Println("Schema does not contain any queries")
	} else {
		for _, f := range queries.Fields {
			f := f
			c.Queries = append(c.Queries, &f)
		}
	}

	mutations, ok := c.Types["Mutation"]
	if !ok || len(mutations.Fields) == 0 {
		log.Println("Schema does not contain any mutations")
	} else {
		for _, f := range mutations.Fields {
			f := f
			c.Mutations = append(c.Mutations, &f)
		}
	}

	_, ok = c.Types["Subscription"]
	if ok {
		log.Println("Schema contains subscriptions, these are not supported")
	}

	return nil
}

func (c *Client) ExecuteFile(filename string) (string, error) {
	q := utils.LoadQuery(filename)
	req := request.BuildFromString(q, nil)
	resp, err := c.Execute([]byte(req))
	if err != nil {
		return "", fmt.Errorf("error executing request: %v", err)
	}

	return string(resp), nil
}

func (c *Client) Build(requestField *models.Field, variables map[string]any, t request.RequestType) string {
	if t != request.Query && t != request.Mutation {
		panic(fmt.Sprintf("invalid request type: %s", t))
	}

	var input string
	if len(requestField.Args) > 0 {
		input = fmt.Sprintf(" (%s)", requestField.Args[0].Compile())
	}

	requestString := fmt.Sprintf("%s%s{\n%s\n}", t, input, c.CompileField(requestField))

	return string(request.BuildFromString(requestString, variables))
}
