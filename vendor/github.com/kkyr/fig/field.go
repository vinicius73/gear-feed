package fig

import (
	"fmt"
	"reflect"
	"strings"
)

// flattenCfg recursively flattens a cfg struct into
// a slice of its constituent fields.
func flattenCfg(cfg interface{}, tagKey string) []*field {
	root := &field{
		v:        reflect.ValueOf(cfg).Elem(),
		t:        reflect.ValueOf(cfg).Elem().Type(),
		sliceIdx: -1,
	}
	fs := make([]*field, 0)
	flattenField(root, &fs, tagKey)
	return fs
}

// flattenField recursively flattens a field into its
// constituent fields, filling fs as it goes.
func flattenField(f *field, fs *[]*field, tagKey string) {
	for (f.v.Kind() == reflect.Ptr || f.v.Kind() == reflect.Interface) && !f.v.IsNil() {
		f.v = f.v.Elem()
		f.t = f.v.Type()
	}

	switch f.v.Kind() {
	case reflect.Struct:
		for i := 0; i < f.t.NumField(); i++ {
			unexported := f.t.Field(i).PkgPath != ""
			embedded := f.t.Field(i).Anonymous
			if unexported && !embedded {
				continue
			}
			child := newStructField(f, i, tagKey)
			*fs = append(*fs, child)
			flattenField(child, fs, tagKey)
		}

	case reflect.Slice, reflect.Array:
		switch f.t.Elem().Kind() {
		case reflect.Struct, reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
			for i := 0; i < f.v.Len(); i++ {
				child := newSliceField(f, i, tagKey)
				flattenField(child, fs, tagKey)
			}
		}
	}
}

// newStructField is a constructor for a field that is a struct
// member. idx is the field's index in the struct. tagKey is the
// key of the tag that contains the field alt name (if any).
func newStructField(parent *field, idx int, tagKey string) *field {
	f := &field{
		parent:   parent,
		v:        parent.v.Field(idx),
		t:        parent.v.Field(idx).Type(),
		st:       parent.t.Field(idx),
		sliceIdx: -1,
	}
	f.structTag = parseTag(f.st.Tag, tagKey)
	return f
}

// newStructField is a constructor for a field that is a slice
// member. idx is the field's index in the slice. tagKey is the
// key of the tag that contains the field alt name (if any).
func newSliceField(parent *field, idx int, tagKey string) *field {
	f := &field{
		parent:   parent,
		v:        parent.v.Index(idx),
		t:        parent.v.Index(idx).Type(),
		st:       parent.st,
		sliceIdx: idx,
	}
	f.structTag = parseTag(f.st.Tag, tagKey)
	return f
}

// field is a settable field of a config object.
type field struct {
	parent *field

	v        reflect.Value
	t        reflect.Type
	st       reflect.StructField
	sliceIdx int // >=0 if this field is a member of a slice.

	structTag
}

// name is the name of the field. if the field contains an alt name
// in the struct that name is used, else  it falls back to
// the field's name as defined in the struct.
// if this field is a slice field, then its name is simply its
// index in the slice.
func (f *field) name() string {
	if f.sliceIdx >= 0 {
		return fmt.Sprintf("[%d]", f.sliceIdx)
	}
	if f.altName != "" {
		return f.altName
	}
	return f.st.Name
}

// path is a dot separated path consisting of all the names of
// the field's ancestors starting from the topmost parent all the
// way down to the field itself.
func (f *field) path() (path string) {
	var visit func(f *field)
	visit = func(f *field) {
		if f.parent != nil {
			visit(f.parent)
		}
		path += f.name()
		// if it's a slice/array we don't want a dot before the slice indexer
		// e.g. we want A[0].B instead of A.[0].B
		if f.t.Kind() != reflect.Slice && f.t.Kind() != reflect.Array {
			path += "."
		}
	}
	visit(f)
	return strings.Trim(path, ".")
}

// parseTag parses a fields struct tags into a more easy to use structTag.
// key is the key of the struct tag which contains the field's alt name.
func parseTag(tag reflect.StructTag, key string) (st structTag) {
	if val, ok := tag.Lookup(key); ok {
		i := strings.Index(val, ",")
		if i == -1 {
			i = len(val)
		}
		st.altName = val[:i]
	}

	if val := tag.Get("validate"); val == "required" {
		st.required = true
	}

	if val, ok := tag.Lookup("default"); ok {
		st.setDefault = true
		st.defaultVal = val
	}

	return
}

// structTag contains information gathered from parsing a field's tags.
type structTag struct {
	altName    string // the alt name of the field as defined in the tag.
	required   bool   // true if the tag contained a required validation key.
	setDefault bool   // true if tag contained a default key.
	defaultVal string // the value of the default key.
}
