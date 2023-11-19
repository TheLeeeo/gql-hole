package crawler

import (
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
		if strings.Contains(e.Message, "PermissionDenied") {
			return true
		}
	}

	return false
}
