package tpl

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// Templates is a specialized set of templates from "text/template"
type Templates struct {
	common *template.Template
	set    map[string]*template.Template
}

// New allocates a new set of templates that can be inited via ParseDir function
func New() *Templates {
	return &Templates{
		common: template.New(""),
		set:    make(map[string]*template.Template),
	}
}

// Must is a helper that wraps a call to a function returning (*Templates, error)
// and panics if the error is non-nil. It is intended for use in variable initializations
// such as
//	var t = tpl.Must(tpl.New().ParseDir("html", []string{".html"}))
func Must(t *Templates, err error) *Templates {
	if err != nil {
		panic(err)
	}

	return t
}

// Delims sets the action delimiters to the specified strings, to be used in
// subsequent calls to Parse, ParseFiles, or ParseGlob. Nested template
// definitions will inherit the settings. An empty delimiter stands for the
// corresponding default: {{ or }}.
// The return value is the template, so calls can be chained.
func (t *Templates) Delims(left, right string) *Templates {
	t.common.Delims(left, right)

	return t
}

// Funcs adds the elements of the argument map to the template's function map.
// It must be called before the template is parsed.
// It panics if a value in the map is not a function with appropriate return
// type. However, it is legal to overwrite elements of the map. The return
// value is the template, so calls can be chained.
func (t *Templates) Funcs(funcMap template.FuncMap) *Templates {
	t.common.Funcs(funcMap)

	return t
}

// Execute applies a parsed template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
// A template may be executed safely in parallel, although if parallel
// executions share a Writer the output may be interleaved.
func (t *Templates) Execute(wr io.Writer, name string, data interface{}) error {
	tpl := t.Lookup(name)
	if tpl == nil {
		return fmt.Errorf("template %s is not found", name)
	}

	return tpl.Execute(wr, data)
}

// Lookup returns the template with the given name or nil if there is no such template.
func (t *Templates) Lookup(name string) *template.Template {
	return t.set[name]
}

// ParseDir parses all files from specified directory that names end with specified extensions.
// Templates may have been extended if they contain extends tag on the first line.
func (t *Templates) ParseDir(dir string, exts ...string) (*Templates, error) {
	root := strings.TrimSuffix(filepath.ToSlash(filepath.Clean(dir)), "/") + "/"

	files, err := parseDir(root, exts...)
	if err != nil {
		return t, err
	}

	for k := range files {
		if files[k].parent == nil {
			err := t.addCommon(files[k])
			if err != nil {
				return t, err
			}
		}
	}

	for k := range files {
		chain := []string{k}
		parent := files[k].parent

		for parent != nil {
			p, ok := files[*parent]
			if !ok {
				break
			}

			chain = append(chain, *parent)
			parent = p.parent
		}

		for i := len(chain) - 1; i >= 0; i-- {
			err = t.add(files[chain[i]])
			if err != nil {
				return t, err
			}
		}
	}

	return t, nil
}

func (t *Templates) addCommon(f *tplFile) error {
	if f == nil {
		return nil
	}

	_, err := t.common.New(f.path).Parse(string(f.content))

	return err
}

func (t *Templates) add(f *tplFile) error {
	if f == nil {
		return nil
	}

	if f.parent == nil {
		tpl := template.Must(t.common.Clone()).New(f.path)

		_, err := tpl.Parse(string(f.content))
		if err != nil {
			return err
		}

		t.set[f.path] = tpl

		return nil
	}

	parent, ok := t.set[*f.parent]
	if !ok {
		return fmt.Errorf("extend template mentioned in %s is not found", f.abspath)
	}

	tpl := template.Must(parent.Clone())

	_, err := tpl.Parse(string(f.content))
	if err != nil {
		return err
	}

	t.set[f.path] = tpl

	return nil
}
