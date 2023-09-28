package schema

type builder[T any] interface {
	buildSchema() schema[T]
	builderUntyped
}

type builderUntyped interface {
	buildSchemaUntyped() genericSchema
}
