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

// For fetching the mutations
const mutationIntrospectionQuery = `
query{
    __schema {
		mutationType {
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
        }
    }
}

`

// For fetching the queries
const queryIntrospectionQuery = `
query{
    __schema {
        queryType {
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
        }
    }
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
