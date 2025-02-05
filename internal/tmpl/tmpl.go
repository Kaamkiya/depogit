package tmpl

import "html/template"

var Templates = template.Must(template.ParseGlob("templates/*.tmpl"))
