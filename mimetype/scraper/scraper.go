package main

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/iancoleman/strcase"
	"os"
	"regexp"
)

var format = `
	// %s is the constant for the mime type "%s".
	%s = "%s"
`

var re = regexp.MustCompile("[/+-.]")

func main() {
	url := "https://www.iana.org/assignments/media-types/media-types.xhtml"
	if doc, err := xmlquery.LoadURL(url); err != nil {
		panic(err)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, `// Package mimetype defines the constants for all IANA registered MIME Types.
// See https://www.iana.org/assignments/media-types for the official list.
// The advantage of using the constants from this package over string literals is correctness.
// The compiler can spot errors that the runtime couldn't.
// Example:
//     // Spelling mistakes, not caught by the compiler
//     w.Header().Add("Content-Tpye", "applicaiton/xml")
//
//     // Spelling mistakes, caught by the compiler
//     w.Header().Add(header.ContentTpye, mimetype.ApplicatoinXml)
//
// Note: Currently, only a subset is implemented.
// The goal is to implement all registered MIME Types by scraping.
package mimetype
`)
		superTypes := []string{"application", "audio", "font", "example", "image", "message", "model", "multipart", "text", "video"}
		for _, superType := range superTypes {
			scrapeTypes(doc, superType)
		}
	}
}

func scrapeTypes(doc *xmlquery.Node, superType string) {
	_, _ = fmt.Fprintf(os.Stdout, "\nconst (")
	for _, row := range xmlquery.Find(doc, fmt.Sprintf("//table[@id='table-%s']/tbody/tr", superType)) {
		typeName := superType + "/" + xmlquery.FindOne(row, "td[1]").InnerText()
		constantName := constantNameForTypeName(typeName)
		_, _ = fmt.Fprintf(os.Stdout, format, constantName, typeName, constantName, typeName)
	}
	_, _ = fmt.Fprintf(os.Stdout, ")\n")
}

func constantNameForTypeName(t string) string {
	return strcase.ToCamel(re.ReplaceAllString(t, "_"))
}
