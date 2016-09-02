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
// Using formtree, the syntax would be
//
//     tree.Slice("fields").Map(1).Slice("content").Map(1).Values("postcode")
//
package formtree

import (
	"strconv"
	"strings"

	"github.com/crantok/imath"
)

// FormTree is a tree of form values.
type FormTree map[string]interface{}

// Map returns the FormTree corresponding to the given key.
func (f FormTree) Map(key string) FormTree {
	result := f[key]
	if result == nil {
		return nil
	}
	return f[key].(FormTree)
}

// Slice returns the Slice corresponding to the given key.
func (f FormTree) Slice(key string) Slice {
	result := f[key]
	if result == nil {
		return nil
	}
	return f[key].(Slice)
}

// Values returns the form values corresponding to the given key.
func (f FormTree) Values(key string) []string {
	result := f[key]
	if result == nil {
		return nil
	}
	return f[key].([]string)
}

// Slice is one kind of node in a FormTree.
type Slice []interface{}

// Map returns the FormTree at the given index.
func (s Slice) Map(index int) FormTree {
	return s[index].(FormTree)
}

// Slice returns the slice at the given index.
func (s Slice) Slice(index int) Slice {
	return s[index].(Slice)
}

// Values returns the form values at the given index.
func (s Slice) Values(index int) []string {
	return s[index].([]string)
}

func addValuesToTree(m FormTree, keyPath []string, values []string) {

	for _, v := range keyPath[:len(keyPath)-1] {
		if m[v] == nil {
			m[v] = FormTree{}
		}
		m = m[v].(FormTree)
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

func sliceify(m FormTree) interface{} {

	isSlice := true
	var indexes []int

	for k, v := range m {

		// sliceify from the leaves to the root
		if child, isNotLeaf := v.(FormTree); isNotLeaf {
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

	slice := make(Slice, imath.Max(indexes...)+1)
	for _, idx := range indexes {
		slice[idx] = m[strconv.Itoa(idx)]
	}
	return slice
}

// New returns a tree of values derived from an http.Request.PostForm or
// equivalent input.
func New(form map[string][]string) FormTree {

	result := FormTree{}

	for k, v := range form {
		addValuesToTree(result, decomposeKeyPath(k), v)
	}

	// Discarding the final result of sliceify because we are going to return a
	// map.
	// ASSUMPTION: Not all elements of a form will have names beggining with
	// integers.
	sliceify(result)

	return result
}
