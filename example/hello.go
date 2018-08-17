package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/dradtke/vdom2"
	"honnef.co/go/js/dom"
)

const tmpl = `<p>{{ .Message }}</p>`

type Hello struct {
	Message string

	tmpl *template.Template
	root dom.Element
}

func (hello *Hello) Render() {
	// Execute the template with the given todo and write to a buffer
	buf := bytes.NewBuffer(nil)
	if err := hello.tmpl.Execute(buf, hello); err != nil {
		panic(err)
	}
	vdom.Apply(hello.root, buf.Bytes())
}

func main() {
	root := dom.GetWindow().Document().GetElementByID("app")
	tmpl := template.Must(template.New("hello").Parse(tmpl))
	hello := &Hello{Message: "Hello", tmpl: tmpl, root: root}
	hello.Render()

	time.AfterFunc(time.Second, func() {
		hello.Message = "World!"
		hello.Render()
	})
}
