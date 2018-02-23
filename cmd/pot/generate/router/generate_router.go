package router

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/wzshiming/go-swagger/swagger"
	"github.com/wzshiming/go-swagger/swaggergen"
)

func GenerateRouter(pkg, routers, controllers, out string) error {

	if pkg == "" {
		pkg = "routers"
	}

	swagger := &swagger.Swagger{}
	swaggergen.GB(swagger, routers, controllers)

	t := "temp"
	temp := template.New(t)
	_, err := temp.Parse(mkRouter)
	if err != nil {
		return err
	}

	// 解析模板
	buf := bytes.NewBuffer(nil)
	err = temp.ExecuteTemplate(buf, "temp", map[string]interface{}{
		"Swagger":  swagger,
		"Package":  pkg,
		"Function": strings.Split(path.Base(out), ".")[0],
	})
	if err != nil {
		return err
	}

	// 格式化源码
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	out = path.Join(out)
	// 写文件
	d, _ := ioutil.ReadFile(out)
	if string(d) == string(src) {
		fmt.Println("[pot] Unchanged " + out)
	} else {
		err = ioutil.WriteFile(out, src, 0666)
		if err != nil {
			return err
		}
		fmt.Println("[pot] Generate " + out)
	}
	return nil
}

var mkRouter = `// Code generated by "pot gen rou"; DO NOT EDIT.
package {{.Package}}
{{if .Swagger.Paths}}
import (
	controllers "{{.Swagger.Extensions.Package}}"
	"net/http"

	"gopkg.in/pot.v1"
	"gopkg.in/pot.v1/router"
)

func {{.Function}}(p *pot.Pot, rs ...*router.Router) (r *router.Router) {
	
	if len(rs) != 0 {
		r = rs[0]
	}
	
	if r == nil {
		r = router.NewRouter()
	}

	paths := r.PathPrefix("{{.Swagger.BasePath}}").Subrouter()

	// init controllers
	{{range $k, $v := .Swagger.Paths}}
	{{with .Get}}
	paths.
		Methods(http.MethodGet).
		Path("{{$k}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := controllers.{{.Extensions.Controllers}}{}
			t.Init(p, w, r)
			t.{{.Extensions.Methods}}()
		})
	{{end}}
	{{with .Post}}
	paths.
		Methods(http.MethodPost).
		Path("{{$k}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := controllers.{{.Extensions.Controllers}}{}
			t.Init(p, w, r)
			t.{{.Extensions.Methods}}()
		})
	{{end}}
	{{with .Put}}
	paths.
		Methods(http.MethodPut).
		Path("{{$k}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := controllers.{{.Extensions.Controllers}}{}
			t.Init(p, w, r)
			t.{{.Extensions.Methods}}()
		})
	{{end}}
	{{with .Delete}}
	paths.
		Methods(http.MethodDelete).
		Path("{{$k}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := controllers.{{.Extensions.Controllers}}{}
			t.Init(p, w, r)
			t.{{.Extensions.Methods}}()
		})
	{{end}}
	{{with .Options}}
	paths.
		Methods(http.MethodOptions).
		Path("{{$k}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := controllers.{{.Extensions.Controllers}}{}
			t.Init(p, w, r)
			t.{{.Extensions.Methods}}()
		})
	{{end}}
	{{with .Head}}
	paths.
		Methods(http.MethodHead).
		Path("{{$k}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := controllers.{{.Extensions.Controllers}}{}
			t.Init(p, w, r)
			t.{{.Extensions.Methods}}()
		})
	{{end}}
	{{with .Patch}}
	paths.
		Methods(http.MethodPatch).
		Path("{{$k}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := controllers.{{.Extensions.Controllers}}{}
			t.Init(p, w, r)
			t.{{.Extensions.Methods}}()
		})
	{{end}}
	{{end}}
	return r
}
{{end}}
`
