package schema

func New[pS *S, S any](builder *structBuilder[S]) Schema[pS] {
	schema := builder.buildSchema()
	return &_TopLevel[pS, S]{
		sig:     sigFor(schema),
		_Struct: schema.(*_Struct[S]),
	}
}
