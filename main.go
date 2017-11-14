package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"path/filepath"

	"path"

	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/unders/tmplgen/html"
	"github.com/unders/tmplgen/toml"
)

var from = ""
var to = ""
var dataFile = ""
var indexFile = "index.html"

const usage = `Usage of tmplgen:
    tmplgen -from=dir -to=dir files

Example:
    tmplgen -from=templates -to=websites/public page/blog/a-very-long-name.html`

func init() {
	flag.StringVar(&from, "from", "", "from directory")
	flag.StringVar(&to, "to", "", "to directory")
	flag.Parse()

	if from == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	if to == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	dataFile = filepath.Join(from, "data.toml")
	indexFile = filepath.Join(to, indexFile)
}

func main() {
	pages := flag.Args()
	l := len(pages)
	if l == 0 {
		fmt.Println(usage)
		return
	}
	p := pages[0]
	if l > 1 {
		p = dataFile
	}

	p = path.Clean(p)

	data, err := toml.ReadFile(dataFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	template := html.NewTemplate(from)
	pager := Pager{template: template, data: data, destDir: path.Clean(to)}
	if p == "all" {
		pager.renderAll()
		return
	}

	if p == dataFile {
		pager.renderAll()
		return
	}

	if strings.HasPrefix(p, template.LayoutDir) {
		pager.renderAll()
		return
	}

	if strings.HasPrefix(p, template.PartDir) {
		pager.renderAll()
		return
	}

	pager.render(p)
}

// Pager creates pages
type Pager struct {
	template *html.Template
	data     *toml.Data
	destDir  string
}

func (p Pager) renderAll() {
	fmt.Println("render all pages")

	pages := []string{}
	walkparts := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "page file walk failed")
		}
		if f.IsDir() {
			return nil
		}

		pages = append(pages, path)
		return nil
	}
	if err := filepath.Walk(p.template.PageDir, walkparts); err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range pages {
		p.render(file)
	}
}

func (p Pager) render(file string) {
	filename := strings.TrimPrefix(file, p.template.PageDir+"/")
	layout := p.data.Layout(filename)
	b, err := p.template.Execute(layout, filename, p.data)
	if err != nil {
		fmt.Printf("template error: %s [ layout: %s, file: %s ]\n", err, layout, filename)
		return
	}

	destFile := path.Join(p.destDir, filename)
	if err := os.MkdirAll(path.Dir(destFile), os.ModePerm); err != nil {
		fmt.Printf("could not create dir %s, error: %s\n", path.Dir(destFile), err)
	}
	if err := ioutil.WriteFile(destFile, b, os.ModePerm); err != nil {
		fmt.Printf("write %s error: %s [ layout %s, file: %s ]\n", destFile, err, layout, filename)
		return
	}
	fmt.Printf("write file: %s [ layout: %s, file: %s ]\n", destFile, layout, filename)
}
