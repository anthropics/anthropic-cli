package jsonflag

import (
	"testing"

	"github.com/urfave/cli/v3"
)

func TestApply(t *testing.T) {
	Clear()

	globalRegistry.Register(Body, "name", "test")
	globalRegistry.Register(Query, "page", 1)
	globalRegistry.Register(Header, "authorization", "Bearer token")

	body, query, header, err := globalRegistry.ApplyMutations(
		[]byte(`{}`),
		[]byte(`{}`),
		[]byte(`{}`),
	)

	if err != nil {
		t.Fatalf("Failed to apply mutations: %v", err)
	}

	expectedBody := `{"name":"test"}`
	expectedQuery := `{"page":1}`
	expectedHeader := `{"authorization":"Bearer token"}`

	if string(body) != expectedBody {
		t.Errorf("Body mismatch. Expected: %s, Got: %s", expectedBody, string(body))
	}
	if string(query) != expectedQuery {
		t.Errorf("Query mismatch. Expected: %s, Got: %s", expectedQuery, string(query))
	}
	if string(header) != expectedHeader {
		t.Errorf("Header mismatch. Expected: %s, Got: %s", expectedHeader, string(header))
	}
}
