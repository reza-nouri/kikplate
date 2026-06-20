package generate

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func marshalJSON(v any) (*bytes.Reader, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
