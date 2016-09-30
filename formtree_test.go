package formtree_test

import (
	"net/url"
	"testing"

	"github.com/crantok/formtree"
)

func Test(t *testing.T) {

	var tests = []struct {
		key    string
		values []string
		index  int
		fn     func(formtree.FormTree) string
	}{
		{
			"a.b.1.2.c",
			[]string{"val1"},
			0,
			func(ft formtree.FormTree) string {
				return ft.MapAt("a").SliceAt("b").SliceAt(1).MapAt(2).ValueAt("c")
			},
		},
		{
			"a.b.1.2.c",
			[]string{"val1", "val4"},
			0,
			func(ft formtree.FormTree) string {
				return ft.MapAt("a").SliceAt("b").SliceAt(1).MapAt(2).ValueAt("c")
			},
		},
		{
			"1.2.a.b.3",
			[]string{"val2", "val3"},
			1,
			func(ft formtree.FormTree) string {
				return ft.SliceAt("1").MapAt(2).MapAt("a").SliceAt("b").AllValuesAt(3)[1]
			},
		},
		{
			"a.b.1.2.c",
			[]string{"val5"},
			0,
			func(ft formtree.FormTree) string {
				return ft.MapAt("a").SliceAt("b").SliceAt(1).MapAt(2).AllValuesAt("c")[0]
			},
		},
	}

	for _, x := range tests {
		form := url.Values{x.key: x.values}
		tree := formtree.New(form)
		if x.fn(tree) != x.values[x.index] {
			t.Errorf(`Examined %v (value at index %v), expected %q, got %q`, form, x.index, x.values[x.index], x.fn(tree))
		}
	}
}
