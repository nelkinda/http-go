package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node
	return d.DecodeElement((*node)(n), &start)
}

func main() {
	url := "https://www.iana.org/assignments/media-types/media-types.xhtml"
	fmt.Fprintf(os.Stderr, "%s\n", url)
}
