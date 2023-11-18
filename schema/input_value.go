package schema

import "fmt"

type InputValue struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         *Type  `json:"type"`
	DefaultValue string `json:"defaultValue"`
}

func (i *InputValue) Compile() string {
	argTypeString := buildArgTypeString(i.Type)

	return fmt.Sprintf("$%s: %s", i.Name, argTypeString)
}
