// Package formtree creates a tree of form values from a url.Values or any other
// type derived from map[string][]string . (One use of the url.Values type is
// http.Request.PostForm). The created tree has the kind of structure that
// json.Unmarshal builds when unmarhsalling to an empty interface. The leaf
// nodes of the tree all have type []string and are the values taken from the
// original map.
//
// Keys (e.g. html form field names used as the keys in http.Request.PostForm)
// are interpreted in the same way that Gorilla Schema interprets them when
// populating a struct, so the form values that had the the key
//
//     "fields.0.content.3.postcode"
//
// would be located at (pseudocode)
//
//     tree["fields"][0]["content"][3]["postcode"]
//
// Using formtree, the syntax would be
//
//     tree.Slice("fields").Map(0).Slice("content").Map(3).Values("postcode")
//
// or
//
//     tree.Slice("fields").Map(0).Slice("content").Map(3).Value("postcode")
//                                                         -----
//
// if you knew that there would be only one value.
package formtree

import (
	"strconv"
	"strings"

	"github.com/crantok/imath"
)

// FormTree is a tree of form values and, like Slice, may be a node within a
// FormTree.
type FormTree map[string]interface{}

// Map returns the FormTree corresponding to the given key.
func (f FormTree) Map(key string) FormTree {
	result := f[key]
	if result == nil {
		return nil
	}
	return result.(FormTree)
}

// Slice returns the Slice corresponding to the given key.
func (f FormTree) Slice(key string) Slice {
	result := f[key]
	if result == nil {
		return nil
	}
	return result.(Slice)
}

// Values returns the form values corresponding to the given key.
func (f FormTree) Values(key string) []string {
	values := f[key]
	if values == nil {
		return nil
	}
	if slice, isSlice := values.([]string); isSlice {
		return slice
	}
	return []string{values.(string)}
}

// Value returns the first form value corresponding to the given key.
func (f FormTree) Value(key string) string {
	value := f[key]
	if value == nil {
		return ""
	}
	if slice, isSlice := value.([]string); isSlice {
		if len(slice) == 0 {
			return ""
		}
		return slice[0]
	}
	return value.(string)
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
	values := s[index]
	if slice, isSlice := values.([]string); isSlice {
		return slice
	}
	return []string{values.(string)}
}

// Value returns the first form value at the given index.
func (s Slice) Value(index int) string {
	value := s[index]
	if slice, isSlice := value.([]string); isSlice {
		if len(slice) == 0 {
			return ""
		}
		return slice[0]
	}
	return value.(string)
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

	// Add the leaf node.
	//
	// If the leaf only contains one value then add a single string, not the
	// slice. This increases similarity with the output of json.Unmarshal and
	// friends.
	//
	// We would not normally expect to see an empty slice in a url.Values
	// instance. I've decided to use nil rather than "" or []string{}  in this
	// case.
	switch len(values) {
	case 0:
		// do nothing - no leaf node
	case 1:
		m[k] = values[0]
	default:
		m[k] = values
	}
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

// New returns a new FormTree whose structure is derived from the stucture of
// the keys in the input map.
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
