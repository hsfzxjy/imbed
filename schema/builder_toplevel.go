package schema

func New[S any](builder *structBuilder[S]) Schema[*S] {
	return new_Toplevel(builder.buildStruct())
}
