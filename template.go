package embtemplate

import (
	"fmt"
	"os"
	"path/filepath"
	texttmpl "text/template"

	rice "github.com/GeertJohan/go.rice"
)

type embTemplate struct {
	Box   *rice.Box
	Files []string
}

var embeddedTemplate = &embTemplate{}

// LoadTemplates finds and parses all templates embedded in the application
func LoadTemplates() (*texttmpl.Template, error) {
	return embeddedTemplate.LoadTemplates()
}

func (embT *embTemplate) LoadTemplates() (*texttmpl.Template, error) {
	var err error
	config := &rice.Config{LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateWorkingDirectory, rice.LocateFS}}
	embT.Box, err = config.FindBox("templates")
	if err != nil {
		return nil, err
	}
	embT.Files = embT.listTemplates()
	return embT.parseFiles(nil, embT.Files...)
}

func (embT *embTemplate) listTemplates() (files []string) {
	embT.Box.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func ListTemplates() []string {
	return embeddedTemplate.Files
}

// parseFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func (embT *embTemplate) parseFiles(t *texttmpl.Template, filenames ...string) (*texttmpl.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := embT.Box.Bytes(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filepath.Base(filename)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *texttmpl.Template
		if t == nil {
			t = texttmpl.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
