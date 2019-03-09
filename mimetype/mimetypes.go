// Package mimetypes defines the constants for all IANA registered MIME Types.
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

const (
	// ApplicationHealthJSON defines the MIME Type "application/health+json".
	ApplicationHealthJSON = "application/health+json"

	// ApplicationJSON defines the MIME Type "application/json".
	ApplicationJSON = "application/json"

	// ApplicationJwt defines the MIME Type "application/jwt".
	ApplicationJwt = "application/jwt"

	// ApplicationXhtmlXml defines the MIME Type "application/xhtml+xml".
	ApplicationXhtmlXml = "application/xhtml+xml"

	// ApplicationXml defines the MIME Type "application/xml".
	ApplicationXml = "application/xml"

	// ImagePng defines the MIME Type "image/png".
	ImagePng = "image/png"

	// ImageSvg defines the MIME Type "image/svg+xml".
	ImageSvg = "image/svg+xml"

	// TextCss defines the MIME Type "text/css".
	TextCss = "text/css"

	// TextHtml defines the MIME Type "text/html".
	TextHtml = "text/html"

	// TextJavascript defines the MIME Type "text/javascript".
	TextJavascript = "text/javascript"

	// TextPlain defines the MIME Type "text/plain".
	TextPlain = "text/plain"
)