package asset

import "github.com/hsfzxjy/imbed/formatter"

var FmtFields = []*formatter.Field[Asset]{
	{
		Name:   "Basename",
		Header: "BASENAME",
		Show:   true,
		Getter: func(a Asset) any {
			return a.BaseName()
		},
	},
	{
		Name:   "Size",
		Header: "SIZE",
		Show:   true,
		Getter: func(a Asset) any {
			return a.Content().BytesReader().Len()
		},
	},
	{
		Name:   "Url",
		Header: "URL",
		Show:   true,
		Getter: func(a Asset) any {
			return a.(*asset).model.Url
		},
	},
}
