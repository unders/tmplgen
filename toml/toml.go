package toml

import (
	"errors"
	"io/ioutil"

	"sort"

	"github.com/naoina/toml"
)

type layout struct {
	Path     string
	Filename string
}

// Data parsed from given file
type Data struct {
	Layouts []layout
	String  map[string]string
	Array   map[string][]map[string]string
}

// ReadFile returns *Data with values parsed from the file
func ReadFile(path string) (*Data, error) {
	data := &Data{}
	if path == "" {
		return data, errors.New("no TOML file")
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return data, errors.New("could not read TOML file")
	}

	if err := toml.Unmarshal(buf, &data); err != nil {
		return nil, err
	}

	sort.Sort(byLongestPath(data.Layouts))
	return data, nil
}

// Layout return the layout used for given filename
func (d *Data) Layout(filename string) string {
	const layout = "main.html"

	filename = "/" + filename

	size := len(d.Layouts)
	if size == 0 {
		return layout
	}
	if size == 1 {
		return d.Layouts[0].Filename
	}

	for _, l := range d.Layouts {
		size := len(l.Path)
		if size > len(filename) {
			continue
		}
		if filename[0:size] == l.Path {
			return l.Filename
		}
	}

	return layout
}

type byLongestPath []layout

func (b byLongestPath) Len() int           { return len(b) }
func (b byLongestPath) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byLongestPath) Less(i, j int) bool { return len(b[j].Path) < len(b[i].Path) }
