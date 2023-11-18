package introspection

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/TheLeeeo/gql-test-suite/client"
	"github.com/TheLeeeo/gql-test-suite/client/request"
	"github.com/TheLeeeo/gql-test-suite/schema"
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

type Introspector struct {
	Cfg Config

	gqlClient *client.Client
}

func New(cfg Config) *Introspector {
	gqlClient := client.New(cfg.TargetUrl)

	return &Introspector{
		Cfg:       cfg,
		gqlClient: gqlClient,
	}
}

func (c *Introspector) SetTargetURL(targetURL string) error {
	if targetURL == c.Cfg.TargetUrl {
		return nil
	}

	_, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("error parsing target addr: %v", err)
	}

	c.Cfg.TargetUrl = targetURL

	return nil
}

func (c *Introspector) StartPolling(pollCallback func(*schema.Schema)) {
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
			s, err := c.FetchSchema()
			if err != nil {
				log.Printf("error loading schema: %v", err)
			}
			if pollCallback != nil {
				pollCallback(s)
			}
		}
	}()
}

func (c *Introspector) SetHeaders(headers map[string]string) {
	c.Cfg.Headers = headers
}

func (c *Introspector) FetchType(typeName string) (*schema.Type, error) {
	t, err := c.fetchTypeInternal(typeName, defaultTypeDepth)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// The internal function for fetching a type. Deals with incomplete types
func (c *Introspector) fetchTypeInternal(typeName string, typeDepth int) (*schema.Type, error) {
	ofTypeField := buildRecursiveOfTypeField(typeDepth)
	reqString := fmt.Sprintf(typeIntrospectionQuery, typeName, ofTypeField)

	req := request.New(reqString, nil)
	req.Headers = c.Cfg.Headers

	resp, err := c.gqlClient.Execute(req)
	if err != nil {
		return nil, err
	}

	dataMap := resp.Data["__type"].(map[string]any)

	t := schema.Type{}
	err = utils.ParseMap(dataMap, &t)
	if err != nil {
		return nil, err
	}

	if isCompleteType(t) {
		return &t, nil
	}

	if typeDepth == 0 {
		typeDepth = 1
	}

	return c.fetchTypeInternal(typeName, typeDepth*2)
}

// Checks if a type is comleted.
// It is considered complete if it knows the base typ of all of its fields
func isCompleteType(t schema.Type) bool { //TODO: check inputFIelds and args
	baseType := t.GetBaseType()
	if baseType.Kind == schema.NonNullTypeKind || baseType.Kind == schema.ListTypeKind {
		return false
	}

	for _, f := range t.Fields {
		if !isCompleteField(f) {
			return false
		}
	}

	return true
}

// Checks if a field has all of its types complted
func isCompleteField(f schema.Field) bool {
	if f.Type == nil {
		return false
	}

	if !isCompleteType(*f.Type) {
		return false
	}

	for _, a := range f.Args {
		if a.Type == nil {
			return false
		}

		if !isCompleteType(*a.Type) {
			return false
		}
	}

	return true
}

func (c *Introspector) FetchSchema() (*schema.Schema, error) {
	if c.Cfg.TargetUrl == "" {
		return nil, ErrNoTargetAddr
	}

	q := fmt.Sprintf(schemaIntrospectionQuery, buildRecursiveOfTypeField(defaultTypeDepth))
	req := request.New(q, nil)
	req.Headers = c.Cfg.Headers

	resp, err := c.gqlClient.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}

	dataMap, ok := resp.Data["__schema"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("error parsing request, no valid __schema field found")
	}

	sch := &schema.Schema{}
	err = utils.ParseMap(dataMap, sch)
	if err != nil {
		return nil, fmt.Errorf("error parsing schema: %v", err)
	}

	for i := 0; i < len(sch.Types); i++ {
		t := sch.Types[i]

		if isCompleteType(t) {
			sch.Types[i] = t
		} else {
			t, err := c.FetchType(t.Name)
			if err != nil {
				return nil, fmt.Errorf("error fetching type %s: %v", t.Name, err)
			}
			sch.Types[i] = *t
		}
	}

	return sch, nil
}
