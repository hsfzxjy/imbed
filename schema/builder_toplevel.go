package schema

func New[S any](builder *StructBuilder[S]) Schema[*S] {
	return new_Toplevel(builder.buildStruct())
}
