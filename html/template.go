package html

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Template is used to generate a HTML page
type Template struct {
	LayoutDir string
	PartDir   string
	PageDir   string
}

// NewTemplate returns a *Template
func NewTemplate(root string) *Template {
	return &Template{
		LayoutDir: filepath.Join(root, "layout"),
		PartDir:   filepath.Join(root, "part"),
		PageDir:   filepath.Join(root, "page"),
	}
}

// Execute parses a file with the given layout template
//
// Usage:
//
//        data := map[string]string{ "Title": "A Title" }
//        t := html.NewTemplate("resource/template")
//        b, err := t.Execute("main.html", "post.html", data)
//
func (t *Template) Execute(layout, filename string, data interface{}) ([]byte, error) {
	p, err := newpage(t, layout)
	if err != nil {
		return nil, err
	}

	return p.execute(filename, data)
}

type page struct {
	template *template.Template
	pagedir  string
}

// newpage returns a *page
func newpage(o *Template, layout string) (*page, error) {
	filename := filepath.Join(o.LayoutDir, layout)
	t, err := template.ParseFiles(filename)
	if cerr := checkErr(err, filename); cerr != nil {
		return nil, errors.WithStack(cerr)
	}

	parts := []string{}
	walkparts := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "parts file walk failed")
		}
		if f.IsDir() {
			return nil
		}

		parts = append(parts, path)
		return nil
	}
	if err = filepath.Walk(o.PartDir, walkparts); err != nil {
		return nil, errors.WithStack(err)
	}

	t, err = t.ParseFiles(parts...)
	if cerr := checkErr(err, o.PartDir+"..."); cerr != nil {
		return nil, errors.WithStack(cerr)
	}

	// Execution stops immediately with an error.
	t.Option("missingkey=error")

	return &page{template: t, pagedir: o.PageDir}, nil
}
func (p *page) execute(filename string, data interface{}) ([]byte, error) {
	file := filepath.Join(p.pagedir, filename)
	content, err := ioutil.ReadFile(file)
	if cerr := checkErr(err, file); cerr != nil {
		return nil, errors.WithStack(cerr)
	}

	t, err := p.template.Parse(string(content))
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse content")
	}

	b := bytes.Buffer{}

	if err := t.Execute(&b, data); err != nil {
		return nil, errors.WithStack(err)
	}

	return b.Bytes(), nil
}

func checkErr(err error, filename string) error {
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return errors.Wrapf(err, "%s file not found", filename)
	}
	if os.IsPermission(err) {
		return errors.Wrapf(err, "%s forbidden access", filename)
	}
	return errors.WithStack(err)
}
