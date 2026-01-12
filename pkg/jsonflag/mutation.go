package jsonflag

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type MutationKind string

const (
	Body   MutationKind = "body"
	Query  MutationKind = "query"
	Header MutationKind = "header"
)

type Mutation struct {
	Kind  MutationKind
	Path  string
	Value any
}

type registry struct {
	mutations []Mutation
}

var globalRegistry = &registry{}

func (r *registry) Mutate(kind MutationKind, path string, value any) {
	r.mutations = append(r.mutations, Mutation{
		Kind:  kind,
		Path:  path,
		Value: value,
	})
}

func (r *registry) Apply(body, query, header []byte) ([]byte, []byte, []byte, error) {
	var err error

	for _, mutation := range r.mutations {
		switch mutation.Kind {
		case Body:
			body, err = jsonSet(body, mutation.Path, mutation.Value)
		case Query:
			query, err = jsonSet(query, mutation.Path, mutation.Value)
		case Header:
			header, err = jsonSet(header, mutation.Path, mutation.Value)
		}
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to apply mutation %s.%s: %w", mutation.Kind, mutation.Path, err)
		}
	}

	return body, query, header, nil
}

func (r *registry) Clear() {
	r.mutations = nil
}

func (r *registry) List() []Mutation {
	result := make([]Mutation, len(r.mutations))
	copy(result, r.mutations)
	return result
}

// Mutate adds a mutation that will be applied to the specified kind of data
func Mutate(kind MutationKind, path string, value any) {
	globalRegistry.Mutate(kind, path, value)
}

// ApplyMutations applies all registered mutations to the provided JSON data
func ApplyMutations(body, query, header []byte) ([]byte, []byte, []byte, error) {
	return globalRegistry.Apply(body, query, header)
}

// ClearMutations removes all registered mutations from the global registry
func ClearMutations() {
	globalRegistry.Clear()
}

// ListMutations returns a copy of all currently registered mutations
func ListMutations() []Mutation {
	return globalRegistry.List()
}

func jsonSet(json []byte, path string, value any) ([]byte, error) {
	keys := strings.Split(path, ".")
	path = ""
	for _, key := range keys {
		if key == "#" {
			key = strconv.Itoa(len(gjson.GetBytes(json, path).Array()) - 1)
		}

		if len(path) > 0 {
			path += "."
		}
		path += key
	}
	return sjson.SetBytes(json, path, value)
}
