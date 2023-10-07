package asset

import (
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/formatter"
)

var FmtFields = []*formatter.Field[Asset]{
	{
		Name:   "oid",
		Header: "OID",
		Show:   true,
		Getter: func(a Asset) any {
			return a.(*asset).model.OID
		},
	},
	{
		Name:   "originId",
		Header: "ORIGIN",
		Show:   true,
		Getter: func(a Asset) any {
			return a.(*asset).model.OriginOID
		},
	},
	{
		Name:   "name",
		Header: "NAME",
		Show:   true,
		Getter: func(a Asset) any {
			return a.BaseName()
		},
	},
	{
		Name:   "url",
		Header: "URL",
		Show:   true,
		Getter: func(a Asset) any {
			url := a.(*asset).model.Url
			if url == "" {
				return "<none>"
			}
			return url
		},
	},
	{
		Name:   "fhash",
		Header: "FHASH",
		Show:   true,
		Getter: func(a Asset) any {
			return a.(*asset).model.FID.Hash()
		},
	},
	{
		Name:   "created",
		Header: "CREATED",
		Show:   true,
		Getter: func(a Asset) any {
			return a.(*asset).model.Created
		},
	},
	{
		Name:   "size",
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
