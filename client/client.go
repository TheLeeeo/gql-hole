package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

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
	Cfg Config

	hasSchema bool

	Types     map[string]*models.Type
	Queries   []models.Field
	Mutations []models.Field
}

func New(cfg Config) *Client {
	return &Client{
		Cfg:   cfg,
		Types: make(map[string]*models.Type),
	}
}

func (c *Client) HasSchema() bool {
	return c.hasSchema
}

// Sets the target url of the client and tries to fetch a schema.
// Returns true if the target url could be set
// Can return both true and an error if the target url was set but the schema could not be fetched
func (c *Client) SetTargetURL(targetURL string) (bool, error) {
	if targetURL == c.Cfg.TargetUrl && c.hasSchema {
		return true, nil
	}

	_, err := url.Parse(targetURL)
	if err != nil {
		return false, fmt.Errorf("error parsing target addr: %v", err)
	}

	c.Cfg.TargetUrl = targetURL
	c.hasSchema = false

	//Load the new schema
	err = c.LoadSchema()
	if err != nil {
		return true, fmt.Errorf("error loading schema: %v", err)
	}

	return true, nil
}

func (c *Client) SetHeaders(headers map[string]string) {
	c.Cfg.Headers = headers
}

func (c *Client) StartPolling() {
	if !c.Cfg.PollingConfig.Enabled {
		return
	}

	go func() {
		for {
			time.Sleep(time.Duration(c.Cfg.PollingConfig.Interval) * time.Minute)

			if c.Cfg.TargetUrl == "" {
				log.Println("No target addr specified, skipping polling")
				continue
			}

			log.Println("Polling for changes to the schema")
			err := c.LoadSchema()
			if err != nil {
				log.Printf("error loading schema: %v", err)
			}
		}
	}()
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

	req, err := http.NewRequest("POST", c.Cfg.TargetUrl, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)

	}
	req.Header.Add("Content-Type", "application/json")

	for k, v := range c.Cfg.Headers {
		req.Header.Add(k, v)
	}

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
	if c.Cfg.TargetUrl == "" {
		return ErrNoTargetAddr
	}

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
	dataMap, ok := respType.Data["__schema"].(map[string]any)
	if !ok {
		return fmt.Errorf("error parsing request, no valid __schema field found")
	}

	sch := &models.Schema{}
	err = utils.ParseMap(dataMap, sch)
	if err != nil {
		return fmt.Errorf("error parsing schema: %v", err)
	}

	for i := 0; i < len(sch.Types); i++ {
		t := sch.Types[i]

		if isCompleteType(&t) {
			c.Types[t.Name] = &t
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
			c.Queries = append(c.Queries, f)
		}
	}

	mutations, ok := c.Types["Mutation"]
	if !ok || len(mutations.Fields) == 0 {
		log.Println("Schema does not contain any mutations")
	} else {
		for _, f := range mutations.Fields {
			f := f
			c.Mutations = append(c.Mutations, f)
		}
	}

	_, ok = c.Types["Subscription"]
	if ok {
		log.Println("Schema contains subscriptions, these are not supported")
	}

	c.hasSchema = true
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
