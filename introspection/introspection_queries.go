package introspection

// For fetching a single type
// The first %s will be replaced with the typename and the second %s will be the recursiveOfTypeField
const typeIntrospectionQuery = `
query TypeQuery{
    __type(name: "%s"){
		...FullType
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

// For fetching the entire schema
// The %s will be replaced with the recursiveOfTypeField
const schemaIntrospectionQuery = `
query IntrospectionQuery {
    __schema {
        description
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

// To be added to a type request in a variable depth to fetch all nested types
// Example: A non-null (1) list (2) of non-null (3) strings (4) would require a depth of 4
const recursiveOfTypeField = `
ofType {
	name
	kind
	%s
}
`
