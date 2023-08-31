package client

// For fetching a type
const typeIntrospectionQuery = `
query TypeQuery{
    __type(name: "%s"){
		name
		kind
		fields {
			name
			type {
				name
				kind
				%s
			}
			args {
				name
				type {
					name
					kind
					%s
				}
			}
		}
        interfaces {
            name
            kind
            %s
        }
        possibleTypes {
            name
            kind
            %s
        }
		enumValues {
			name
		}
		inputFields {
			name
			type {
				name
				kind
				%s
			}
			defaultValue
		}
    }
}
`
const schemaIntrospectionQuery = `
query IntrospectionQuery {
    __schema {    
        queryType { name }
        mutationType { name }
        subscriptionType { name }
        types {
            ...FullType
        }
        directives {
            name
            description
            locations
            args {
                ...InputValue
            }
        }
    }
}

fragment FullType on __Type {
    kind
    name
    description
    fields(includeDeprecated: true) {
        name
        description
        args {
            ...InputValue
        }
        type {
            ...TypeRef
        }
        isDeprecated
        deprecationReason
    }
    inputFields {
        ...InputValue
    }
    interfaces {
        ...TypeRef
    }
    enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
    }
    possibleTypes {
        ...TypeRef
    }
}

fragment InputValue on __InputValue {
    name
    description
    type { ...TypeRef }
    defaultValue 
}

fragment TypeRef on __Type {
    kind
    name
    %s
}
`

// For fetching the names of all types so they can be fetched individually.
// This is to avoid refethcing every type if one was incomplete
const typeNamesIntrospectionQuery = `
query{
    __schema {
		types {
			name
        }
    }
}
`

// To be added to a type request in a variable depth to fetch all nested types
// Example: A non-null (1) list (2) of non-null (3) strings (4) would require a depth of 4
const recursiveOfTypeField = `
ofType {
	name
	kind
	%s
}
`
