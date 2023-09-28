package schema

func New[S any](builder *structBuilder[S]) Schema[S] {
	schema := builder.buildSchema()
	return &_TopLevel[S]{
		sig:     sigFor(schema),
		_Struct: schema.(*_Struct[S]),
	}
}
