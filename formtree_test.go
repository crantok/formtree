package formtree_test

import (
	"net/url"
	"testing"

	"github.com/crantok/formtree"
)

func Test(t *testing.T) {

	/* Test data derived from these form elements:

	   <div class="document-id">
	       <input type="hidden" name="document-id" value="476128394763523">
	   </div>

	   <div class="field.0">
	       <input type="text" name="fields.0.type" value="location">
	       <input type="text" name="fields.0.label" value="HQ">
	       <input type="text" name="fields.0.content.address" value="BlAh">
	       <input type="text" name="fields.0.content.postcode" value="814h">
	   </div>

	   <div class="field.1">
	       <input type="text" name="fields.1.type" value="location">
	       <input type="text" name="fields.1.label" value="outlets">
	       <input type="text" name="fields.1.content.0.address" value="addr 1">
	       <input type="text" name="fields.1.content.0.postcode" value="pc 1">
	       <input type="text" name="fields.1.content.1.address" value="addr 2">
	       <input type="text" name="fields.1.content.1.postcode" value="pc 2">
	   </div>
	*/
	form := url.Values{
		"fields.1.content.1.postcode": []string{"pc 2"},
		"fields.0.content.postcode":   []string{"814h"},
		"fields.1.content.1.address":  []string{"addr 2"},
		"fields.0.label":              []string{"HQ"},
		"fields.0.content.address":    []string{"BlAh"},
		"fields.1.type":               []string{"location"},
		"fields.1.label":              []string{"outlets"},
		"fields.1.content.0.address":  []string{"addr 1"},
		"fields.1.content.0.postcode": []string{"pc 1"},
		"document-id":                 []string{"476128394763523"},
		"fields.0.type":               []string{"location"},
	}

	tree := formtree.New(form)

	treeVal := tree["fields"].([]interface{})[1].(map[string]interface{})["content"].([]interface{})[1].(map[string]interface{})["postcode"].([]string)[0]
	formVal := form["fields.1.content.1.postcode"][0]
	if treeVal != formVal {
		t.Errorf(`Examined form["fields.1.content.1.postcode"][0], expected %v, got %v`, treeVal, formVal)
	}
}
