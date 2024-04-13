package utils

import "text/template"

func MustTemplate(name string, funcs map[string]interface{}, tmpl string) *template.Template {
	parse, err := template.New(name).Funcs(funcs).Parse(tmpl)
	if err != nil {
		panic(err)
	}
	return parse
}
