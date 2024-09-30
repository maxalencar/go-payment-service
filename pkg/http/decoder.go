package http

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

func Decode(r io.Reader, mimeType string, v any) error {
	if mimeType == MIMETypeXML {
		return xml.NewDecoder(r).Decode(v)
	}

	return json.NewDecoder(r).Decode(v)
}
