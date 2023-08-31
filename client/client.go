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

// The number of extra levels of nesting to add to the recursive type field
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
	Addr       string
	Types      map[string]*models.Type
	InputTypes map[string]*models.Type
	EnumTypes  map[string][]models.EnumValue
	Unions     map[string]*models.Type
	Queries    []*models.Field
	Mutations  []*models.Field
}

func New(addr string) *Client {
	return &Client{
		Addr:       addr,
		Types:      make(map[string]*models.Type),
		InputTypes: make(map[string]*models.Type),
		EnumTypes:  make(map[string][]models.EnumValue),
		Unions:     make(map[string]*models.Type),
	}
}

// Gets the type from the client
func (c *Client) GetType(name string) *models.Type {
	t, ok := c.Types[name]
	if !ok {
		// fmt.Printf("type (%s) not found in client, fetching\n", name)
		err := c.FetchType(name)
		if err != nil {
			panic(fmt.Sprintf("error fetching type %s: %v", name, err))
		}

		t = c.Types[name]
	}

	return t
}

// Gets the inputtype from the client
func (c *Client) GetInputType(name string) *models.Type {
	t, ok := c.InputTypes[name]
	if !ok {
		// fmt.Printf("type (%s) not found in client, fetching\n", name)
		err := c.FetchType(name)
		if err != nil {
			panic(fmt.Sprintf("error fetching type %s: %v", name, err))
		}

		t = c.InputTypes[name]
	}

	return t
}

// Gets the enumvalues from the client
func (c *Client) GetEnumValues(name string) []models.EnumValue {
	t, ok := c.EnumTypes[name]
	if !ok {
		// fmt.Printf("type (%s) not found in client, fetching\n", name)
		err := c.FetchType(name)
		if err != nil {
			panic(fmt.Sprintf("error fetching type %s: %v", name, err))
		}

		t = c.EnumTypes[name]
	}

	return t
}

func (c *Client) GetUnionType(name string) *models.Type {
	t, ok := c.Unions[name]
	if !ok {
		// fmt.Printf("type (%s) not found in client, fetching\n", name)
		err := c.FetchType(name)
		if err != nil {
			panic(fmt.Sprintf("error fetching type %s: %v", name, err))
		}

		t = c.Unions[name]
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

// Fethes the type specified by typeName and saves it to the client
func (c *Client) FetchType(typeName string) error {
	if _, ok := c.Types[typeName]; ok { //TODO Handle input, enum and union types
		// fmt.Printf("type %s already fetched, skipping\n", typeName)
		return nil
	}

	// fmt.Println("Fetching type: ", typeName)

	if typeName == "_Service" || typeName == "_Any" || typeName == "_FieldSet" {
		// fmt.Printf("type %s is not supported, skipping\n", typeName)
		return nil
	}

	t, err := c.fetchTypeInternal(typeName, defaultTypeDepth)
	if err != nil {
		return err
	}

	if t.Name == "Query" {
		for _, f := range t.Fields {
			f := f // Reused loop closure variable fix [https://github.com/golang/go/wiki/CommonMistakes]
			if f.Name == "_service" {
				continue
			}

			c.Queries = append(c.Queries, &f)
		}
		return nil
	}

	if t.Name == "Mutation" {
		for _, f := range t.Fields {
			f := f
			c.Mutations = append(c.Mutations, &f)
		}
		return nil
	}

	switch t.Kind {
	case models.ObjectTypeKind:
		c.Types[typeName] = t
	case models.InputObjectTypeKind:
		c.InputTypes[typeName] = t
	case models.EnumTypeKind:
		c.EnumTypes[typeName] = t.EnumValues
	case models.UnionTypeKind:
		c.Unions[t.Name] = t
	}

	return nil
}

// The internal function for fetching a type. Deals with incomplete types
func (c *Client) fetchTypeInternal(typeName string, typeDepth int) (*models.Type, error) {
	ofTypeField := buildRecursiveOfTypeField(typeDepth)
	reqString := fmt.Sprintf(typeIntrospectionQuery, typeName, ofTypeField, ofTypeField, ofTypeField, ofTypeField, ofTypeField)

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

// Fetches all the types in the schema and adds them to the client
func (c *Client) FetchSchema() error {
	typeNames, err := c.fetchTypeNames()
	if err != nil {
		return fmt.Errorf("error fetching type names: %v", err)
	}

	for _, typeName := range typeNames {
		err = c.FetchType(typeName)
		if err != nil {
			return fmt.Errorf("error fetching type %s: %v", typeName, err)
		}
	}

	return nil
}

// The queries have incomplete types consisting of only the name
// Replace the types with their complete variants
func (c *Client) populateTypesInQueries() {
	for _, q := range c.Queries {
		c.populateBaseType(q.Type)
	}

	for _, m := range c.Mutations {
		c.populateBaseType(m.Type)
	}
}

func (c *Client) populateBaseType(t *models.Type) {
	baseType := t.GetBaseType()

	if baseType.Kind == models.ScalarTypeKind {
		return
	}

	if baseType.Kind == models.ObjectTypeKind {
		*baseType = *c.GetType(baseType.Name)
	} else if baseType.Kind == models.InputObjectTypeKind {
		*baseType = *c.GetInputType(baseType.Name)
	} else if baseType.Kind == models.EnumTypeKind {
		baseType.EnumValues = c.GetEnumValues(baseType.Name)
	} else if baseType.Kind == models.UnionTypeKind {
		*baseType = *c.GetUnionType(baseType.Name)
		for _, pt := range baseType.PossibleTypes {
			c.populateBaseType(&pt)
		}
	} else {
		panic(fmt.Sprintf("Unimplemented base type kind %s", baseType.Kind))
	}
}

// Fetches the names of all of the types in the schema
func (c *Client) fetchTypeNames() ([]string, error) {
	request := request.BuildFromString(typeNamesIntrospectionQuery, nil)
	respBytes, err := c.Execute([]byte(request))
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}

	resp, err := utils.ParseResponse(respBytes)
	if err != nil {
		return nil, err
	}
	dataMap := resp.Data["__schema"].(map[string]any)

	sch := &models.Schema{}
	err = utils.ParseMap(dataMap, sch)
	if err != nil {
		return nil, fmt.Errorf("error parsing schema: %v", err)
	}

	var names []string
	for _, t := range sch.Types {
		names = append(names, t.Name)
	}

	return names, nil
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
