package templates

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/shurcooL/github_flavored_markdown"
	"io"
	"net/http"
	"strconv"
)

var (
	tplExt         = ".html"
	templatesCache = make(map[string]*pongo2.Template)
)

func init() {
	if err := pongo2.RegisterFilter("markdown", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsSafeValue(string(github_flavored_markdown.Markdown([]byte(in.String())))), nil
	}); err != nil {
		fmt.Println("err", err.Error())
	}

	if err := pongo2.RegisterFilter("toString", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsValue(strconv.Itoa(in.Integer())), nil
	}); err != nil {
		fmt.Println("err", err.Error())
	}

	//if err := pongo2.RegisterFilter("paginator", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	//
	//	var res []string
	//
	//	for i := 1; i < (in.Integer() / param.Integer()); i++ {
	//		offset := (i - 1) * 10
	//		res = append(res, fmt.Sprintf(`<li class="active"><a href="/post/?pageSize=10&offset=%d">%d</a></li>`, offset, i))
	//	}
	//
	//	return pongo2.AsValue([]byte(strings.Join(res, ""))), nil
	//}); err != nil {
	//	fmt.Println("err", err.Error())
	//}
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
