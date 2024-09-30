package http

import (
	"encoding/json"
	"encoding/xml"
)

func Marshal(mimeType string, v any) ([]byte, error) {
	if mimeType == MIMETypeXML {
		return xml.Marshal(v)
	}

	return json.Marshal(v)
}
