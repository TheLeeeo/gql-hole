package crawler

import (
	"log"

	"github.com/fatih/color"
)

type CrawlOperation struct {
	// The name of the operation
	Name string `json:"name"`
	// The operation request made
	Operation string `json:"operation"`
	// The type of the operation

	// Was the operation considered denied
	Denied bool `json:"success"`
	// Failed to perform the operation
	Failed bool `json:"failed"`

	// The response of the operation
	Response string `json:"response"`

	// The error message if the operation failed
	Error error `json:"error"`
}

func NewOperation(name, operation string, vars map[string]any) CrawlOperation {
	resp := CrawlOperation{
		Name:      name,
		Operation: operation,
	}

	return resp
}

func (o *CrawlOperation) SetResponse(resp []byte) {
	o.Response = string(resp)

	o.Denied = is403Error(resp) || is401Error(resp)
	o.Failed = isFetchFailed(resp)
}

func (o *CrawlOperation) PrintResult() {
	var resultString string

	if o.Failed {
		resultString = color.YellowString("FAILED TO FETCH")
	} else if o.Denied {
		resultString = color.GreenString("DENIED")
	} else {
		resultString = color.RedString("ALLOWED")
	}

	log.Printf("\"%s\": %s\n", o.Name, resultString)
	// Got allowed
	if !o.Failed && !o.Denied {
		log.Println("	Response: ", o.Response)
	}
}
