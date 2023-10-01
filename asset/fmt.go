package asset

import (
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/formatter"
)

var FmtFields = []*formatter.Field[Asset]{
	{
		Name:   "Oid",
		Header: "OID",
		Show:   true,
		Getter: func(a Asset) any { return a.(*asset).model.OID },
	},
	{
		Name:   "Name",
		Header: "NAME",
		Show:   true,
		Getter: func(a Asset) any {
			return a.BaseName()
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
	{
		Name:   "Size",
		Header: "SIZE",
		Show:   true,
		Getter: func(a Asset) any {
			s, err := a.Content().Size()
			if err != nil {
				return content.IllegalSize
			}
			return s
		},
	},
}
