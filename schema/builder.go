package schema

type builder[T any] interface {
	buildSchema() schema[T]
	genericBuilder
}

type genericBuilder interface {
	buildGenericSchema() genericSchema
}
