package crawler

import (
	"encoding/json"
	"strings"

	"github.com/TheLeeeo/gql-test-suite/client"
)

func isFetchFailed(resp []byte) bool {
	respType, err := client.Parse(resp)
	if err != nil {
		panic(err)
	}

	if len(respType.Errors) == 0 {
		return false
	}

	for _, e := range respType.Errors {
		if strings.Contains(e.Message, "HTTP fetch failed") {
			return true
		}
	}

	return false
}

func is401Error(resp []byte) bool {
	respType, err := client.Parse(resp)
	if err != nil {
		panic(err)
	}

	if len(respType.Errors) == 0 {
		return false
	}

	for _, e := range respType.Errors {
		if strings.Contains(strings.ToLower(e.Message), "unauthenticated") {
			return true
		}

		for _, ext := range e.Extensions {
			extString, ok := ext.(string)
			if !ok {
				b, err := json.Marshal(ext)
				if err != nil {
					panic(err)
				}
				extString = string(b)
			}

			if strings.Contains(strings.ToLower(extString), "unauthenticated") {
				return true
			}
		}
	}

	return false
}

func is403Error(resp []byte) bool {
	respType, err := client.Parse(resp)
	if err != nil {
		panic(err)
	}

	if len(respType.Errors) == 0 {
		return false
	}

	for _, e := range respType.Errors {
		if strings.Contains(strings.ToLower(e.Message), "permissiondenied") {
			return true
		}

		for _, ext := range e.Extensions {
			extString, ok := ext.(string)
			if !ok {
				b, err := json.Marshal(ext)
				if err != nil {
					panic(err)
				}
				extString = string(b)
			}

			if strings.Contains(strings.ToLower(extString), "unauthenticated") {
				return true
			}
		}
	}

	return false
}
