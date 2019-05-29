package templates

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"net/http"
)

var (
	tplExt         = ".html"
	templatesCache = make(map[string]*pongo2.Template)
)

func init() {
	if err := pongo2.RegisterFilter("markdown", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsSafeValue(string(blackfriday.Run([]byte(in.String())))), nil
	}); err != nil {
		fmt.Println("err", err.Error())
	}
}

func Render(data map[string]interface{}, body io.Writer, tplName string) error {

	tpl, ok := templatesCache[tplName]
	if !ok {
		tpl = pongo2.Must(pongo2.FromFile(tplName + tplExt))
		//templatesCache[tplName] = tpl
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

func RenderHtml(ctx context.Context, w http.ResponseWriter, response map[string]interface{}) error {
	name := ctx.Value("method").(string)

	buf := new(bytes.Buffer)
	if err := Render(response, buf, "views/"+name); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(buf.Bytes())); err != nil {
		return err
	}

	return nil
}

func PageNotFound() {

}
