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
			func(ft formtree.FormTree) string { return ft.Map("a").Slice("b").Slice(1).Map(2).Value("c") },
		},
		{
			"1.2.a.b.3",
			[]string{"val2", "val3"},
			1,
			func(ft formtree.FormTree) string { return ft.Slice("1").Map(2).Map("a").Slice("b").Values(3)[1] },
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
