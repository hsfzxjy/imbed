package schema

type builder[T any] interface {
	buildSchema() schemaTyped[T]
	builderUntyped
}

type builderUntyped interface {
	buildSchemaUntyped() schema
}
