package http

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

func Encode(w io.Writer, mimeType string, v any) error {
	if mimeType == MIMETypeXML {
		return xml.NewEncoder(w).Encode(v)
	}

	return json.NewEncoder(w).Encode(v)
}
