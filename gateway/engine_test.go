package gateway

import (
	"testing"
)

func TestParseResponse(t *testing.T) {
	data := `{"statusCode":200,"body":"{\"message\":\"Hello World\"}"}`
	resp, err := parseResponse([]byte(data))
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Status code: expected 200, got %v", resp.StatusCode)
		return
	}
	if resp.Body != `{"message":"Hello World"}` {
		t.Errorf("body: {\"message\":\"Hello World\"} got %v", resp.Body)
		return
	}
}
