// The client is the base of the service, responsible for making requests to the server and parsing the responses.
package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/TheLeeeo/gql-test-suite/utils"
)

type Client struct {
	Endpoint string
}

func New(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) Execute(req *Request) (*Response, error) {
	// Create the http request
	requestBody := bytes.NewBuffer(req.Build())
	httpRequest, err := http.NewRequest("POST", c.Endpoint, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	httpRequest.Header.Add("Content-Type", "application/json")

	for k, v := range req.Headers {
		httpRequest.Header.Add(k, v)
	}

	// Send the request
	client := &http.Client{}
	// fmt.Println(httpRequest)
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	// Read the response body
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Try to parse the response into a gql response
	resp, err := Parse(responseBody)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	resp.StatusCode = httpResponse.StatusCode

	return resp, nil
}

func (c *Client) ExecuteFile(filename string) (*Response, error) {
	q := utils.LoadQuery(filename)
	req := NewRequest(q, map[string]any{"params": map[string]any{}})
	resp, err := c.Execute(req)
	fmt.Printf("%+v\n", req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}

	return resp, nil
}
