package uri

import "strings"

// Schema returns a first part of the URI string.
// If the ':' sign is not present in the input string, the input string will be returned.
// Examples:
// "file:./inbox" -> "file"
// "http://server.com" -> "http"
// "Some string" -> "Some string"
// "sql?q=select * from table" -> "sql"
func Schema(uri string) string {
	colonPos := strings.Index(uri, ":")
	if colonPos < 0 {
		questionPos := strings.Index(uri, "?")

		if questionPos > 0 {
			return uri[0:questionPos]
		}

		return uri
	}

	return uri[0:colonPos]
}
