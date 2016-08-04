package main

import (
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

	"github.com/cagnosolutions/web"
)

const salt = "[replace this with something unique]"

var mux = web.NewMux()
var tmpl *web.TmplCache
var projects = "projects"
var Mutex sync.RWMutex

func init() {

	web.Funcs["pretty"] = pretty

	tmpl = web.NewTmplCache()

}

func main() {

	mux.AddRoutes(index, run, format, share, view)

	fmt.Println("REMEMBER TO REGISTER ANY NEW ROUTES")
	log.Fatal(http.ListenAndServe(":8888", mux))
}

var index = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "index.tmpl", web.Model{
		"themes": themes,
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
		"code":   string(b),
		"loaded": true,
		"themes": themes,
	})
}}

var run = web.Route{"POST", "/run", func(w http.ResponseWriter, r *http.Request) {
	dir := "temp"
	path := projects + "/" + dir
	resp := make(map[string]interface{})

	Mutex.Lock()
	defer Mutex.Unlock()
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("main.go >> run >> os.MkdirAll() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		ajaxResponse(w, resp)
		return
	}

	dat := []byte(r.FormValue("dat"))

	if err := ioutil.WriteFile(path+"/main.go", dat, 0644); err != nil {
		log.Printf("main.go >> run >> ioutil.WriteFile() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		ajaxResponse(w, resp)
		return
	}

	/*currentDir, err := os.Getwd()
	if err != nil {
		log.Printf(`main.go >> run >> os.Getwd >> %v\n\n`, err)
		resp["output"] = "Server error. Please try again."
		resp["error"] = true
		ajaxResponse(w, resp)
		return
	}

	if err := os.Chdir(path); err != nil {
		log.Printf("main.go >> run >> os.Chdir() >> %v\n\n", err)
		resp["output"] = "Server error. Please try again."
		resp["error"] = true
		ajaxResponse(w, resp)
		return
	}*/

	cmd := exec.Command("go", "clean")
	cmd.Dir = path
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("main.go >> run >> exec.Command.CombinedOutput(\"go\", \"clean\") >> %v\n\n", err)
		resp["output"] = "Server error. Please try again."
		resp["error"] = true
		ajaxResponse(w, resp)
		return
	}

	cmd = exec.Command("go", "build")
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
		ajaxResponse(w, resp)
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
		ajaxResponse(w, resp)
		return
	}

	/*out, err := runCmd(10, "go", "run", path+"/main.go")
	if err != nil {
		if err == BuildErr {
			log.Printf("main.go >> format >> cmd.CombinedOutput() >> %v\n\n", err)
			//out := fmt.Sprintf("%s", b)
			out = strings.Replace(out, "# command-line-arguments\n", "", -1)
			out = strings.Replace(out, "temp/main.go:", "Line ", -1)
			resp["output"] = out
		} else if err == TimeOutErr {
			resp["output"] = "Process took too long."
		} else {
			resp["output"] = "Server error. Please try again."
		}
		resp["error"] = true
		ajaxResponse(w, resp)
		return
	}*/

	/*if err := os.Chdir(currentDir); err != nil {
		log.Printf("main.go >> run >> os.Chdir() >> %v\n\n", err)
		resp["output"] = "Server error. Please try again."
		resp["error"] = true
		ajaxResponse(w, resp)
		return
	}*/

	resp["error"] = false
	resp["output"] = out
	ajaxResponse(w, resp)
	return
}}

var format = web.Route{"POST", "/format", func(w http.ResponseWriter, r *http.Request) {
	dir := "temp"
	path := projects + "/" + dir
	resp := make(map[string]interface{})

	Mutex.Lock()
	defer Mutex.Unlock()
	if err := os.MkdirAll(path, 0755); err != nil {
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		ajaxResponse(w, resp)
		return
	}

	dat := []byte(r.FormValue("dat"))

	if err := ioutil.WriteFile(path+"/main.go", dat, 0644); err != nil {
		log.Printf("main.go >> run >> ioutil.WriteFile() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again"
		ajaxResponse(w, resp)
		return
	}

	var com string
	//var cmd *exec.Cmd
	if r.FormValue("imp") == "true" {
		com = "goimports"
		//cmd = exec.Command("goimports", path+"/main.go")
	} else {
		com = "gofmt"
		//cmd = exec.Command("gofmt", path+"/main.go")
	}

	cmd := exec.Command(com, "main.go")
	cmd.Dir = path
	b, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("main.go >> format >> cmd.CombinedOutput(%q, %q) >> %v\n\n", com, "main.go", err)
		//resp["output"] = fmt.Sprintf("Server Error. Please try again.")

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
		ajaxResponse(w, resp)
		return
	}

	resp["error"] = false
	resp["output"] = fmt.Sprintf("%s", b)
	ajaxResponse(w, resp)
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
		ajaxResponse(w, resp)
		return
	}

	if err := ioutil.WriteFile(path+"/main.go", dat, 0644); err != nil {
		log.Printf("main.go >> run >> ioutil.WriteFile() >> %v\n\n", err)
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		ajaxResponse(w, resp)
		return
	}

	resp["error"] = false
	resp["output"] = dir
	ajaxResponse(w, resp)
	return

}}

func ajaxResponse(w http.ResponseWriter, resp map[string]interface{}) {
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
