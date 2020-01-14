package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
	TemplateDir        = "views/"
)

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render the view with the predefined layout.
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

//NewView parses the template files
func NewView(layout string, files ...string) *View {
	addTemplatePaths(files)
	addTemplateExt(files)
	files = append(files,
		layoutPaths()...,
	)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// addTemplatePaths takes in a slice of strings
// representing files paths for Templates and it prepends
// the TemplateDir directory to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePaths(files []string) []string {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
	return files
}

// addTemplatePaths takes in a slice of strings
// representing files paths for Templates and it appends
// the TemplateExt to each string in the slice
//
// Eg the input {"views/home"} would result in the output
// {"views/home.gohtml"} if TemplateDir == ".gohtml"
func addTemplateExt(files []string) []string {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
	return files
}

func layoutPaths() []string {
	filePaths, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return filePaths
}
