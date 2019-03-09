// Package mimetypes defines the constants for all IANA registered MIME Types.
// See https://www.iana.org/assignments/media-types for the official list.
// The advantage of using the constants from this package over string literals is correctness.
// The compiler can spot errors that the runtime couldn't.
// Example:
//     // Spelling mistake, not caught by the compiler
//     w.Header().Add("Content-Type", "applicaiton/xml")
//
//     // Spelling mistake, caught by the compiler
//     w.Header().Add("Content-Type", ApplicatoinXml)
package mimetypes

const (
)