// Package formtree derives a tree of form values from an http.Request.PostForm.
// The tree has the kind of structure that json.Unmarshal builds when
// unmarhsalling to a map[string]interface{}. The leaf nodes of the tree all
// have type []string and are the values taken from the PostForm.
//
// PostForm keys are interpreted in the same way that Gorilla Schema interprets
// them when populating a struct, so the form values that had the the key
//
//     "fields.0.content.3.postcode"
//
// would be located at (pseudocode)
//
//     tree["fields"][0]["content"][3]["postcode"]
//
// Each of those indexing operations would require a type assertion to cast from
// interface{} to map[string]interface{} or []interface{} or []string.
package formtree

import (
	"strconv"
	"strings"

	"github.com/crantok/imath"
)

func addValueToMap(m map[string]interface{}, keyPath []string, values []string) {

	for _, v := range keyPath[:len(keyPath)-1] {
		if m[v] == nil {
			m[v] = map[string]interface{}{}
		}
		m = m[v].(map[string]interface{})
	}

	k := keyPath[len(keyPath)-1]

	if m[k] != nil {
		panic("Adding value to map with a key that already exists.")
	}

	m[k] = values
}

func decomposeKeyPath(key string) []string {

	parts := strings.Split(key, ".")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			result = append(result, p)
		}
	}

	return result
}

func sliceify(m map[string]interface{}) interface{} {

	isSlice := true
	var indexes []int

	for k, v := range m {

		// sliceify from the leaves to the root
		if child, isNotLeaf := v.(map[string]interface{}); isNotLeaf {
			m[k] = sliceify(child)
		}

		if isSlice {
			if i, err := strconv.Atoi(k); err != nil {
				isSlice = false
			} else {
				indexes = append(indexes, i)
			}
		}
	}

	if !isSlice || len(indexes) == 0 {
		return m
	}

	slice := make([]interface{}, imath.Max(indexes...)+1)
	for _, idx := range indexes {
		slice[idx] = m[strconv.Itoa(idx)]
	}
	return slice
}

// New returns a tree of values derived from an http.Request.PostForm or
// equivalent input.
func New(form map[string][]string) map[string]interface{} {

	result := map[string]interface{}{}

	for k, v := range form {
		addValueToMap(result, decomposeKeyPath(k), v)
	}

	// Discarding the final result of sliceify because we are going to return a
	// map.
	// ASSUMPTION: Not all elements of a form will have names beggining with
	// integers.
	sliceify(result)

	return result
}
