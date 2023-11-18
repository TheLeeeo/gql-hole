package schema

type Directive struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Locations    []DirectiveLocation `json:"locations"` // should really be directiveLocation
	Args         []InputValue        `json:"args"`
	IsRepeatable bool                `json:"isRepeatable"`
}

type DirectiveLocation string

const (
	QueryDirectiveLocation                DirectiveLocation = "QUERY"
	MutationDirectiveLocation             DirectiveLocation = "MUTATION"
	SubscriptionDirectiveLocation         DirectiveLocation = "SUBSCRIPTION"
	FieldDirectiveLocation                DirectiveLocation = "FIELD"
	FragmentDefinitionDirectiveLocation   DirectiveLocation = "FRAGMENT_DEFINITION"
	FragmentSpreadDirectiveLocation       DirectiveLocation = "FRAGMENT_SPREAD"
	InlineFragmentDirectiveLocation       DirectiveLocation = "INLINE_FRAGMENT"
	VariableDefinitionDirectiveLocation   DirectiveLocation = "VARIABLE_DEFINITION"
	SchemaDirectiveLocation               DirectiveLocation = "SCHEMA"
	ScalarDirectiveLocation               DirectiveLocation = "SCALAR"
	ObjectDirectiveLocation               DirectiveLocation = "OBJECT"
	FieldDefinitionDirectiveLocation      DirectiveLocation = "FIELD_DEFINITION"
	ArgumentDefinitionDirectiveLocation   DirectiveLocation = "ARGUMENT_DEFINITION"
	InterfaceDirectiveLocation            DirectiveLocation = "INTERFACE"
	UnionDirectiveLocation                DirectiveLocation = "UNION"
	EnumDirectiveLocation                 DirectiveLocation = "ENUM"
	EnumValueDirectiveLocation            DirectiveLocation = "ENUM_VALUE"
	InputObjectDirectiveLocation          DirectiveLocation = "INPUT_OBJECT"
	InputFieldDefinitionDirectiveLocation DirectiveLocation = "INPUT_FIELD_DEFINITION"
)
