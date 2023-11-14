package http

import "html/template"

var errorAlertTmpl *template.Template = template.Must(template.ParseFiles("web/components/error_alert.html"))
