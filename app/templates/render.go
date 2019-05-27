package templates

import (
	"encoding/json"
	"github.com/flosch/pongo2"
	"io"
)

var (
	tplExt         = ".html"
	templatesCache = make(map[string]*pongo2.Template)
)

func Render(data interface{}, body io.Writer, tplName string) error {

	tpl, ok := templatesCache[tplName]
	if !ok {
		tpl = pongo2.Must(pongo2.FromFile(tplName + tplExt))
		templatesCache[tplName] = tpl
	}

	b, _ := json.Marshal(data)

	var ctxData pongo2.Context
	if err := json.Unmarshal(b, &ctxData); err != nil {
		// todo
	}

	if err := tpl.ExecuteWriter(ctxData, body); err != nil {
		return err
	}

	return nil
}

func PageNotFound() {

}
