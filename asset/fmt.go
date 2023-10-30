package asset

import (
	"strings"

	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/formatter"
)

var FmtFields = []*formatter.Field[StockAsset]{
	{
		Name:   "oid",
		Header: "OID",
		Show:   true,
		Getter: func(a StockAsset) any {
			return a.Model().OID
		},
	},
	{
		Name:   "originId",
		Header: "ORIGIN",
		Show:   true,
		Getter: func(a StockAsset) any {
			return a.Model().OriginOID
		},
	},
	{
		Name:   "name",
		Header: "NAME",
		Show:   true,
		Getter: func(a StockAsset) any {
			return a.BaseName()
		},
	},
	{
		Name:   "url",
		Header: "URL",
		Show:   true,
		Getter: func(a StockAsset) any {
			url := a.Model().Url
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
		Getter: func(a StockAsset) any {
			return a.Model().FID.Hash()
		},
	},
	{
		Name:   "created",
		Header: "CREATED",
		Show:   true,
		Getter: func(a StockAsset) any {
			return a.Model().Created
		},
	},
	{
		Name:   "size",
		Header: "SIZE",
		Show:   true,
		Getter: func(a StockAsset) any {
			s, err := a.Content().Size()
			if err != nil {
				return content.IllegalSize
			}
			return s
		},
	},
	{
		Name:   "tags",
		Header: "TAGS",
		Show:   true,
		Getter: func(sa StockAsset) any {
			tags := sa.Model().Tags
			return strings.Join(tags, ", ")
		},
	},
}
