package crawler

import (
	"fmt"
	"log"
	"time"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/models"
	"github.com/TheLeeeo/gql-test-suite/request"
	"golang.org/x/exp/slices"
)

type Crawler struct {
	cfg *Config

	// The introspection client
	client *client.Client
	// Queries and mutations to ignored
}

var defaultUnsupportedQueries = []string{
	"_entities",
	"_service",
}

// New creates a new crawler
func New(cfg *Config) *Crawler {
	cfg.Ignored = append(cfg.Ignored, defaultUnsupportedQueries...)

	c := client.New(cfg.ClientConfig)

	return &Crawler{
		client: c,
		cfg:    cfg,
	}
}

func (c *Crawler) IsReady() bool {
	return c.client.HasSchema()
}

func (c *Crawler) SetTargetURL(targetURL string) error {
	return c.client.SetTargetURL(targetURL)
}

func (c *Crawler) SetIgnored(ignored []string) {
	c.cfg.Ignored = ignored
}

func (c *Crawler) GetIgnored() []string {
	return c.cfg.Ignored
}

func (c *Crawler) StartPolling() {
	c.client.StartPolling()
}

func (c *Crawler) Crawl() ([]*CrawlOperation, error) {
	if !c.IsReady() {
		err := c.client.LoadSchema()
		if err != nil {
			return nil, err
		}
	}

	return c.testAllOperations(), nil
}

func (c *Crawler) Do(op *CrawlOperation) error {
	resp, err := c.client.Execute([]byte(op.Operation))
	if err != nil {
		op.Error = err
	}

	op.SetResponse(resp)

	return err
}

func (c *Crawler) TestQuery(queryName string) *CrawlOperation {
	var query *models.Field
	for _, q := range c.client.Queries {
		if q.Name == queryName {
			query = q
			break
		}
	}
	if query == nil {
		log.Println("Query not found")
		return nil
	}

	vars := c.GenerateMinimalTestDataForRequest(query)
	r := c.client.Build(query, vars, request.Query)

	operation := NewOperation(queryName, r, vars)

	c.Do(operation)

	return operation
}

func (c *Crawler) TestMutation(mutationName string) *CrawlOperation {
	var mutation *models.Field
	for _, m := range c.client.Mutations {
		if m.Name == mutationName {
			mutation = m
			break
		}
	}
	if mutation == nil {
		log.Println("Mutation not found")
		return nil
	}

	vars := c.GenerateMinimalTestDataForRequest(mutation)
	r := c.client.Build(mutation, vars, request.Mutation)

	operation := NewOperation(mutationName, r, vars)

	c.Do(operation)

	return operation
}

func (c *Crawler) testAllOperations() []*CrawlOperation {
	var allOperations []*CrawlOperation

	for _, q := range c.client.Queries {
		if slices.Contains(c.cfg.Ignored, q.Name) {
			continue
		}

		vars := c.GenerateMinimalTestDataForRequest(q)
		r := c.client.Build(q, vars, request.Query)

		crawlResp := NewOperation(q.Name, r, vars)

		c.Do(crawlResp)

		allOperations = append(allOperations, crawlResp)
	}

	for _, m := range c.client.Mutations {
		if slices.Contains(c.cfg.Ignored, m.Name) {
			continue
		}
		vars := c.GenerateMinimalTestDataForRequest(m)
		r := c.client.Build(m, vars, request.Mutation)

		crawlResp := NewOperation(m.Name, r, vars)

		c.Do(crawlResp)

		allOperations = append(allOperations, crawlResp)
	}

	return allOperations
}

func (c *Crawler) GenerateMinimalTestDataForRequest(f *models.Field) map[string]any {
	if len(f.Args) == 0 {
		return nil
	}

	vars := make(map[string]any)

	arg := f.Args[0]

	if arg.Type.Kind != models.NonNullTypeKind {
		//Optional args :)
		return nil
	}

	baseType := arg.Type.GetBaseType()

	var value any
	switch baseType.Kind {
	case models.EnumTypeKind:
		value = f.Type.EnumValues
	case models.ScalarTypeKind:
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
	case models.InputObjectTypeKind:
		completeBaseType := c.client.GetType(baseType.Name)
		value = c.GenerateMinimalTestDataForType(completeBaseType)
	default:
		panic(fmt.Sprintf("Unimplemented variable kind %s", f.Type.Kind))
	}

	vars[arg.Name] = value

	return vars
}

func (c *Crawler) GenerateMinimalTestDataForType(t *models.Type) map[string]any {
	vars := make(map[string]any)

	for _, f := range t.InputFields {
		if f.Type.Kind != models.NonNullTypeKind {
			continue
		}

		baseType := f.Type.GetBaseType()

		var value any
		switch baseType.Kind {
		case models.EnumTypeKind:
			value = c.client.GetType(baseType.Name).EnumValues[0].Name
		case models.ScalarTypeKind:
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
		case models.InputObjectTypeKind:
			completeBaseType := c.client.GetType(baseType.Name)
			value = c.GenerateMinimalTestDataForType(completeBaseType)
		default:
			panic(fmt.Sprintf("Unimplemented variable kind %s", f.Type.Kind))
		}

		vars[f.Name] = value
	}

	return vars
}
