package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/client/request"
	"github.com/TheLeeeo/gql-test-suite/introspection"
	"github.com/TheLeeeo/gql-test-suite/schema"
	"github.com/TheLeeeo/gql-test-suite/schema/manager"
	"golang.org/x/exp/slices"
)

type Crawler struct {
	cfg Config

	// The introspection client
	intrClient *introspection.Introspector

	gqlClient *client.Client

	schemaManager *manager.Manager
}

var defaultUnsupportedQueries = []string{
	"_entities",
	"_service",
}

// New creates a new crawler
func New(cfg Config) *Crawler {
	cfg.Ignore = append(cfg.Ignore, defaultUnsupportedQueries...)

	ic := introspection.New(cfg.ClientConfig)

	gqlC := client.New(cfg.ClientConfig.TargetUrl)

	return &Crawler{
		intrClient: ic,
		gqlClient:  gqlC,
		cfg:        cfg,
	}
}

func (c *Crawler) IsReady() bool {
	return c.schemaManager != nil
}

func (c *Crawler) SetTargetURL(targetURL string) error {
	return c.intrClient.SetTargetURL(targetURL)
}

func (c *Crawler) GetTargetURL() string {
	return c.intrClient.Cfg.TargetUrl
}

func (c *Crawler) SetIgnore(ignore []string) {
	c.cfg.Ignore = append(ignore, defaultUnsupportedQueries...)
}

func (c *Crawler) GetIgnore() []string {
	return c.cfg.Ignore
}

func (c *Crawler) StartPolling() {
	c.intrClient.StartPolling(nil)
}

func (c *Crawler) Crawl() ([]CrawlOperation, error) {
	if !c.IsReady() {
		s, err := c.intrClient.FetchSchema()
		if err != nil {
			return nil, err
		}

		c.schemaManager = manager.New(s)
	}

	return c.testAllOperations(), nil
}

func (c *Crawler) Do(op *CrawlOperation) error {
	resp, err := c.gqlClient.Execute(&op.Request)
	if err != nil {
		op.Error = err
	}

	r, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	op.SetResponse(r)

	return err
}

func (c *Crawler) TestQuery(queryName string) *CrawlOperation {
	var query *schema.Field
	for _, q := range c.schemaManager.Queries {
		if q.Name == queryName {
			query = &q
			break
		}
	}
	if query == nil {
		log.Println("Query not found")
		return nil
	}

	vars := c.GenerateMinimalTestDataForRequest(query)
	r := c.schemaManager.Build(*query, request.Query)

	request := request.New(r, vars)
	operation := NewOperation(queryName, *request)

	c.Do(&operation)

	return &operation
}

func (c *Crawler) TestMutation(mutationName string) *CrawlOperation {
	var mutation *schema.Field
	for _, m := range c.schemaManager.Mutations {
		if m.Name == mutationName {
			mutation = &m
			break
		}
	}
	if mutation == nil {
		log.Println("Mutation not found")
		return nil
	}

	vars := c.GenerateMinimalTestDataForRequest(mutation)
	r := c.schemaManager.Build(*mutation, request.Mutation)
	req := request.New(r, vars)

	operation := NewOperation(mutationName, *req)

	c.Do(&operation)

	return &operation
}

func (c *Crawler) testAllOperations() []CrawlOperation {
	var allOperations []CrawlOperation

	for _, q := range c.schemaManager.Queries {
		if slices.Contains(c.cfg.Ignore, q.Name) {
			continue
		}

		vars := c.GenerateMinimalTestDataForRequest(&q)
		r := c.schemaManager.Build(q, request.Query)

		req := request.New(r, vars)
		crawlResp := NewOperation(q.Name, *req)

		fmt.Printf("%+v\n", req)

		c.Do(&crawlResp)

		allOperations = append(allOperations, crawlResp)
	}

	for _, m := range c.schemaManager.Mutations {
		if slices.Contains(c.cfg.Ignore, m.Name) {
			continue
		}
		vars := c.GenerateMinimalTestDataForRequest(&m)
		r := c.schemaManager.Build(m, request.Mutation)
		req := request.New(r, vars)
		crawlResp := NewOperation(m.Name, *req)

		c.Do(&crawlResp)

		allOperations = append(allOperations, crawlResp)
	}

	return allOperations
}

func (c *Crawler) GenerateMinimalTestDataForRequest(f *schema.Field) map[string]any {
	if len(f.Args) == 0 {
		return nil
	}

	vars := make(map[string]any)

	arg := f.Args[0]

	if arg.Type.Kind != schema.NonNullTypeKind {
		//Optional args :)
		return nil
	}

	baseType := arg.Type.GetBaseType()

	var value any
	switch baseType.Kind {
	case schema.EnumTypeKind:
		value = f.Type.EnumValues
	case schema.ScalarTypeKind:
		if baseType.Name == "Boolean" {
			value = true
		} else if baseType.Name == "String" {
			value = "0"
		} else if baseType.Name == "Int" {
			value = 0
		} else if baseType.Name == "Float" {
			value = 0.0
		} else if baseType.Name == "ID" {
			value = "0"
		} else if baseType.Name == "Time" {
			value = time.Now()
		} else {
			panic(fmt.Sprintf("Unhandled scalar type %s", baseType.Name))
		}
	case schema.InputObjectTypeKind:
		completeBaseType := c.schemaManager.Types[baseType.Name]
		value = c.GenerateMinimalTestDataForType(&completeBaseType)
	default:
		panic(fmt.Sprintf("Unimplemented variable kind %s", f.Type.Kind))
	}

	vars[arg.Name] = value

	return vars
}

func (c *Crawler) GenerateMinimalTestDataForType(t *schema.Type) map[string]any {
	vars := make(map[string]any)

	for _, f := range t.InputFields {
		if f.Type.Kind != schema.NonNullTypeKind {
			continue
		}

		baseType := f.Type.GetBaseType()

		var value any
		switch baseType.Kind {
		case schema.EnumTypeKind:
			value = c.schemaManager.Types[baseType.Name].EnumValues[0].Name
		case schema.ScalarTypeKind:
			if baseType.Name == "Boolean" {
				value = true
			} else if baseType.Name == "String" {
				value = "0"
			} else if baseType.Name == "Int" {
				value = 0
			} else if baseType.Name == "Float" {
				value = 0.0
			} else if baseType.Name == "ID" {
				value = "0"
			} else if baseType.Name == "Time" {
				value = time.Now()
			} else {
				panic(fmt.Sprintf("Unhandled scalar type %s", baseType.Name))
			}
		case schema.InputObjectTypeKind:
			completeBaseType := c.schemaManager.Types[baseType.Name]
			value = c.GenerateMinimalTestDataForType(&completeBaseType)
		default:
			panic(fmt.Sprintf("Unimplemented variable kind %s", f.Type.Kind))
		}

		vars[f.Name] = value
	}

	return vars
}
