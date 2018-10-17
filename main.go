package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/cagnosolutions/web"

	"github.com/gregpechiro/playground/sizeof"
)

//const DIR = "/home/greg/code/go/src/github.com/gregpechiro/playground/"

const DIR = "/opt/playground/"

const salt = "[replace this with something unique]"

var mux = web.NewMux()
var tmpl *web.TmplCache
var projects = "projects"
var Mutex sync.RWMutex
var SizeOfTemp *template.Template
var Current string

func init() {

	web.Funcs["pretty"] = pretty
	web.Funcs["isCurrent"] = IsCurrent
	SizeOfTemp = prepTemplate()

	tmpl = web.NewTmplCache()
	SetCurrent()
}

func main() {
	os.Setenv("PATH", "")

	mux.AddRoutes(index, run, format, share, view, sizeOf)

	fmt.Println("REMEMBER TO REGISTER ANY NEW ROUTES")
	log.Fatal(http.ListenAndServe(":8888", mux))
}

var index = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {

	tmpl.Render(w, r, "index.tmpl", web.Model{
		"themes":   themes,
		"versions": GetVersions(),
		"current":  Current,
	})
}}

var view = web.Route{"GET", "/:id", func(w http.ResponseWriter, r *http.Request) {
	Mutex.RLock()
	b, err := ioutil.ReadFile(projects + "/" + r.FormValue(":id") + "/" + "main.go")
	Mutex.RUnlock()
	if err != nil {
		fmt.Fprintf(w, "snippet not found")
		return
	}

	tmpl.Render(w, r, "index.tmpl", web.Model{
		"code":     string(b),
		"loaded":   true,
		"themes":   themes,
		"versions": GetVersions(),
		"current":  Current,
	})
}}

var run = web.Route{"POST", "/run", func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	version := r.FormValue("version")
	if version == "" {
		version = "current"
	}
	if _, err := os.Stat("env/versions/" + version); err != nil {
		log.Printf("main.go >> run >> os.Stat(env/versions/%s) >> %v\n\n", version, err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		AjaxResponse(w, resp)
		return
	}
	dir := fmt.Sprintf("%d", time.Now().Unix())
	path := "/tmp/playground/" + dir

	Mutex.Lock()
	defer Mutex.Unlock()
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("main.go >> run >> os.MkdirAll() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		AjaxResponse(w, resp)
		return
	}

	dat := []byte(r.FormValue("dat"))

	if err := ioutil.WriteFile(path+"/main.go", dat, 0644); err != nil {
		log.Printf("main.go >> run >> ioutil.WriteFile() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		AjaxResponse(w, resp)
		return
	}

	cmd := exec.Command(DIR+"env/versions/"+version+"/bin/go", "clean")
	cmd.Env = []string{"GOROOT=" + DIR + "env/versions/" + version, "GOPATH=" + DIR + "env/libs"}
	cmd.Dir = path
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("main.go >> run >> exec.Command.CombinedOutput(\"go\", \"clean\") >> %v\n\n", err)
		resp["output"] = "Server error. Please try again."
		resp["error"] = true
		AjaxResponse(w, resp)
		return
	}

	cmd = exec.Command(DIR+"env/versions/"+version+"/bin/go", "build")
	cmd.Env = []string{"GOROOT=" + DIR + "env/versions/" + version, "GOPATH=" + DIR + "env/libs"}
	cmd.Dir = path
	b, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("main.go >> run >> exec.Command.CombinedOutput(\"go\", \"build\") >> %v\n\n", err)
		out := string(b)
		if i := strings.Index(out, "#"); i == 0 && strings.Contains(out, "\n") {
			out = out[strings.Index(out, "\n")+1:]
		}
		out = strings.Replace(out, "./main.go:", "Line ", -1)
		resp["output"] = out
		resp["error"] = true
		AjaxResponse(w, resp)
		return
	}

	out, err := runCmd(5, path+"/"+dir)
	if err != nil {
		log.Printf("main.go >> run >> runCmd() >> %v\n\n", err)
		if err == ExecErr {
			resp["output"] = out
		} else if err == TimeOutErr {
			resp["output"] = "Process took too long."
		} else {
			resp["output"] = "Server error. Please try again."
		}
		resp["error"] = true
		AjaxResponse(w, resp)
		return
	}

	resp["error"] = false
	resp["output"] = out
	AjaxResponse(w, resp)
	return
}}

var format = web.Route{"POST", "/format", func(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]interface{})
	version := r.FormValue("version")
	if version == "" {
		version = "current"
	}
	if _, err := os.Stat("env/versions/" + version); err != nil {
		log.Printf("main.go >> format >> os.Stat(env/versions/%s) >> %v\n\n", version, err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		AjaxResponse(w, resp)
		return
	}
	dir := fmt.Sprintf("%d", time.Now().Unix())
	path := "/tmp/playground/" + dir

	Mutex.Lock()
	defer Mutex.Unlock()
	if err := os.MkdirAll(path, 0755); err != nil {
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		AjaxResponse(w, resp)
		return
	}

	dat := []byte(r.FormValue("dat"))

	if err := ioutil.WriteFile(path+"/main.go", dat, 0644); err != nil {
		log.Printf("main.go >> run >> ioutil.WriteFile() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again"
		AjaxResponse(w, resp)
		return
	}

	var com string
	if r.FormValue("imp") == "true" {
		com = DIR + "env/tools/bin/goimports"
	} else {
		com = DIR + "env/versions/" + version + "/bin/gofmt"
	}

	cmd := exec.Command(com, "main.go")
	cmd.Env = []string{"GOROOT=" + DIR + "env/versions/" + version, "GOPATH=" + DIR + "env/libs"}
	cmd.Dir = path
	b, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("main.go >> format >> cmd.CombinedOutput(%q, %q) >> %v\n\n", com, "main.go", err)

		out := string(b)
		if i := strings.Index(out, "#"); i == 0 && strings.Contains(out, "\n") {
			out = out[strings.Index(out, "\n")+1:]
		}
		outL := strings.Split(out, ":")
		if len(outL) > 2 {
			out = "Line " + outL[1] + ":" + outL[3]
		} else {
			out = strings.Replace(out, "main.go:", "Line ", -1)
		}

		resp["output"] = out

		resp["error"] = true
		AjaxResponse(w, resp)
		return
	}

	resp["error"] = false
	resp["output"] = fmt.Sprintf("%s", b)
	AjaxResponse(w, resp)
	return
}}

var share = web.Route{"POST", "/share", func(w http.ResponseWriter, r *http.Request) {
	dat := []byte(r.FormValue("dat"))
	dir := GetId(dat)
	path := projects + "/" + dir
	resp := make(map[string]interface{})

	if err := os.MkdirAll(path, 0755); err != nil {
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		AjaxResponse(w, resp)
		return
	}

	if err := ioutil.WriteFile(path+"/main.go", dat, 0644); err != nil {
		log.Printf("main.go >> run >> ioutil.WriteFile() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		AjaxResponse(w, resp)
		return
	}

	resp["error"] = false
	resp["output"] = dir
	AjaxResponse(w, resp)
	return

}}

var sizeOf = web.Route{"POST", "/sizeof", func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	code := r.FormValue("code")
	if code == "" {
		resp["error"] = true
		resp["output"] = "No code was submitted!"
		AjaxResponse(w, resp)
		return
	}

	toRender := &struct {
		Code   string
		Result *sizeof.ViewData
		Error  string
	}{Code: code}

	result, err := sizeof.ParseCode(code)
	if err != nil {
		toRender.Error = err.Error()
	} else {
		toRender.Result = sizeof.CreateViewData(result)
	}

	buf := new(bytes.Buffer)
	err = SizeOfTemp.Execute(buf, toRender)
	if err != nil {
		log.Printf("main.go >> sizeOf >> t.Execute() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Error processing. Please try again"
		AjaxResponse(w, resp)
		return
	}
	resp["error"] = false
	resp["output"] = buf.String()
	AjaxResponse(w, resp)
	return
}}

func AjaxResponse(w http.ResponseWriter, resp map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(resp)
	msg := string(b)
	if err != nil {
		msg = `{"error":true,"msg":"Server error. Please try again."}`
	}
	fmt.Fprint(w, msg)
}

func GetId(doc []byte) string {
	h := sha1.New()
	io.WriteString(h, salt)
	h.Write(doc)
	sum := h.Sum(nil)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(sum)))
	base64.URLEncoding.Encode(b, sum)
	return string(b)[:10]
}

func prepTemplate() *template.Template {
	var fns = template.FuncMap{
		"unvischunk": func(x int, len int) bool {
			return x > 2 && x < (len-1)
		},
	}
	t, err := template.New("Sizeof").Funcs(fns).Parse(SizeOfDisplay)
	if err != nil {
		panic(err)
	}
	return t
}
