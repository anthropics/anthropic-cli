package jsonflag

import (
	"testing"
)

func TestApply(t *testing.T) {
	ClearMutations()

	Mutate(Body, "name", "test")
	Mutate(Query, "page", 1)
	Mutate(Header, "authorization", "Bearer token")

	body, query, header, err := ApplyMutations(
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
