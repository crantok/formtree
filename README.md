# Package formtree

Package formtree creates a tree of form values from a url.Values or any other
type derived from map[string][]string . (One use of the url.Values type is
http.Request.PostForm). The created tree has the kind of structure that
json.Unmarshal builds when unmarhsalling to an empty interface. The leaf nodes
of the tree all have type string or []string and are the values taken from
the original map.

## Why?

I was decoding a form using Gorilla Schema. Part of the form was best
represented by a slice of a struct type from another package. I did not want to
add schema tags to the other package. My options were:

 0. Make the names of the elements in the form match the fields of the
 externally defined struct. (Very sensible, but I didn't want to stray from the
 naming convention in use for form elements.)
 0. Implement a copy of the external struct and add schema tags to it. (This solution lasted a while.)
 0. Read the whole form manually from the PostForm map. (Tedious.)
 0. Write a generic solution to unpack a PostForm into a convenient data structure. (Aha!)

Later on, I decided I wanted to allow my site members to edit a set of documents
where the structure could not be predicted at coding time. That's why I finally
wrote this package.


## Key interpretation

Keys (e.g. html form field names used as the keys in http.Request.PostForm) are
interpreted in the same way that Gorilla Schema interprets them when populating
a struct, so the form values that had the the key

    "fields.0.content.3.postcode"

would be located at (pseudocode)

    tree["fields"][0]["content"][3]["postcode"]

Using formtree, the syntax would be

    tree.Slice("fields").Map(0).Slice("content").Map(3).Values("postcode")

or

    tree.Slice("fields").Map(0).Slice("content").Map(3).Value("postcode")
                                                        -----
if you knew that there would be only one value.


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

which produce this http.Request.PostForm

    url.Values{
        "document-id":                 []string{"476128394763523"},
        "fields.0.type":               []string{"location"},
        "fields.0.label":              []string{"HQ"},
        "fields.0.content.address":    []string{"BlAh"},
        "fields.0.content.postcode":   []string{"814h"},
        "fields.1.type":               []string{"location"},
        "fields.1.label":              []string{"outlets"},
        "fields.1.content.0.address":  []string{"addr 1"},
        "fields.1.content.0.postcode": []string{"pc 1"},
        "fields.1.content.1.address":  []string{"addr 2"},
        "fields.1.content.1.postcode": []string{"pc 2"},
    }

we obtain this tree

    formtree.FormTree{
        "document-id": []string{"476128394763523"},
        "fields": formtree.Slice{
            formtree.FormTree{
                "type":  []string{"location"},
                "label": []string{"HQ"},
                "content": formtree.FormTree{
                    "address":  []string{"BlAh"},
                    "postcode": []string{"814h"},
                },
            },
            formtree.FormTree{
                "type":  []string{"location"},
                "label": []string{"outlets"},
                "content": formtree.Slice{
                    formtree.FormTree{
                        "address":  []string{"addr 1"},
                        "postcode": []string{"pc 1"},
                    },
                    formtree.FormTree{
                        "address":  []string{"addr 2"},
                        "postcode": []string{"pc 2"},
                    },
                },
            },
        },
    }
