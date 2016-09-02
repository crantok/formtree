# Package formtree

Package formtree derives a tree of form values from an http.Request.PostForm.
The tree has the kind of structure that json.Unmarshal builds when
unmarhsalling to a map[string]interface{}. The leaf nodes of the tree all
have type []string and are the values taken from the PostForm.


## Key interpretation

PostForm keys are interpreted in the same way that Gorilla Schema interprets
them when populating a struct, so the form values that had the the key

    "fields.0.content.3.postcode"

would be located at (pseudocode)

    tree["fields"][0]["content"][3]["postcode"]

Each of those indexing operations would require a type assertion to cast from
interface{} to map[string]interface{} or []interface{} or []string.


## Example

Given these elements of a form

    <div class="document-id">
        <input type="hidden" name="document-id" value="476128394763523">
    </div>

    <div class="field-0">
        <input type="text" name="fields.0.type" value="location">
        <input type="text" name="fields.0.label" value="HQ">
        <input type="text" name="fields.0.content.address" value="BlAh">
        <input type="text" name="fields.0.content.postcode" value="814h">
    </div>

    <div class="field-1">
        <input type="text" name="fields.1.type" value="location">
        <input type="text" name="fields.1.label" value="outlets">
        <input type="text" name="fields.1.content.0.address" value="addr 1">
        <input type="text" name="fields.1.content.0.postcode" value="pc 1">
        <input type="text" name="fields.1.content.1.address" value="addr 2">
        <input type="text" name="fields.1.content.1.postcode" value="pc 2">
    </div>

which produce this http.Request.PostForm (output from [go-spew](https://github.com/davecgh/go-spew))

    (url.Values) (len=11) {
    	(string) (len=24) "fields.0.content.address": ([]string) (len=1 cap=1) {
    		(string) (len=4) "BlAh"
    	},
    	(string) (len=13) "fields.1.type": ([]string) (len=1 cap=1) {
    		(string) (len=8) "location"
    	},
    	(string) (len=14) "fields.1.label": ([]string) (len=1 cap=1) {
    		(string) (len=7) "outlets"
    	},
    	(string) (len=26) "fields.1.content.0.address": ([]string) (len=1 cap=1) {
    		(string) (len=6) "addr 1"
    	},
    	(string) (len=27) "fields.1.content.0.postcode": ([]string) (len=1 cap=1) {
    		(string) (len=4) "pc 1"
    	},
    	(string) (len=14) "fields.0.label": ([]string) (len=1 cap=1) {
    		(string) (len=2) "HQ"
    	},
    	(string) (len=13) "fields.0.type": ([]string) (len=1 cap=1) {
    		(string) (len=8) "location"
    	},
    	(string) (len=25) "fields.0.content.postcode": ([]string) (len=1 cap=1) {
    		(string) (len=4) "814h"
    	},
    	(string) (len=26) "fields.1.content.1.address": ([]string) (len=1 cap=1) {
    		(string) (len=6) "addr 2"
    	},
    	(string) (len=27) "fields.1.content.1.postcode": ([]string) (len=1 cap=1) {
    		(string) (len=4) "pc 2"
    	},
    	(string) (len=11) "document-id": ([]string) (len=1 cap=1) {
    		(string) (len=15) "476128394763523"
    	}
    }

we obtain this tree

    (map[string]interface {}) (len=2) {
    	(string) (len=6) "fields": ([]interface {}) (len=2 cap=2) {
    		(map[string]interface {}) (len=3) {
    			(string) (len=5) "label": ([]string) (len=1 cap=1) {
    				(string) (len=2) "HQ"
    			},
    			(string) (len=4) "type": ([]string) (len=1 cap=1) {
    				(string) (len=8) "location"
    			},
    			(string) (len=7) "content": (map[string]interface {}) (len=2) {
    				(string) (len=8) "postcode": ([]string) (len=1 cap=1) {
    					(string) (len=4) "814h"
    				},
    				(string) (len=7) "address": ([]string) (len=1 cap=1) {
    					(string) (len=4) "BlAh"
    				}
    			}
    		},
    		(map[string]interface {}) (len=3) {
    			(string) (len=7) "content": ([]interface {}) (len=2 cap=2) {
    				(map[string]interface {}) (len=2) {
    					(string) (len=7) "address": ([]string) (len=1 cap=1) {
    						(string) (len=6) "addr 1"
    					},
    					(string) (len=8) "postcode": ([]string) (len=1 cap=1) {
    						(string) (len=4) "pc 1"
    					}
    				},
    				(map[string]interface {}) (len=2) {
    					(string) (len=7) "address": ([]string) (len=1 cap=1) {
    						(string) (len=6) "addr 2"
    					},
    					(string) (len=8) "postcode": ([]string) (len=1 cap=1) {
    						(string) (len=4) "pc 2"
    					}
    				}
    			},
    			(string) (len=4) "type": ([]string) (len=1 cap=1) {
    				(string) (len=8) "location"
    			},
    			(string) (len=5) "label": ([]string) (len=1 cap=1) {
    				(string) (len=7) "outlets"
    			}
    		}
    	},
    	(string) (len=11) "document-id": ([]string) (len=1 cap=1) {
    		(string) (len=15) "476128394763523"
    	}
    }
