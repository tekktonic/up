package main;

import (
	"github.com/antchfx/xquery/xml"
)

func getAttr(n *xmlquery.Node, namespace string, local string) string {
	for _,attr := range n.Attr {
		if (attr.Name.Space == namespace &&
			attr.Name.Local == local) {
			return attr.Value
		}
	}
	return ""
}
