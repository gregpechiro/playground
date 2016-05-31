package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/cagnosolutions/web"
)

var mux = web.NewMux()

var tmpl *web.TmplCache

func init() {

	tmpl = web.NewTmplCache()

}

func main() {

	mux.AddRoutes(index, run, format)

	fmt.Println("REMEMBER TO REGISTER ANY NEW ROUTES")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

var index = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "index.tmpl", web.Model{})
}}

var run = web.Route{"POST", "/run", func(w http.ResponseWriter, r *http.Request) {
	path := "temp"
	resp := make(map[string]interface{})
	if err := os.MkdirAll(path, 0755); err != nil {
		resp["error"] = true
		resp["msg"] = "Server error. Please try again."
		ajaxResponse(w, resp)
		return
	}
	dat := []byte(r.FormValue("dat"))
	ioutil.WriteFile(path+"/main.go", dat, 0644)*/

	out, err := runCmd(10, "go", "run", path+"/main.go")

	if err != nil {
		if err == BuildErr {
			log.Printf("main.go >> format >> cmd.CombinedOutput() >> %v\n", err)
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
	}

	resp["error"] = false
	resp["output"] = out
	ajaxResponse(w, resp)
	return
}}

var format = web.Route{"POST", "/format", func(w http.ResponseWriter, r *http.Request) {
	path := "temp"
	resp := make(map[string]interface{})
	if err := os.MkdirAll(path, 0755); err != nil {
		resp["error"] = true
		resp["output"] = "Server error. Please try again."
		ajaxResponse(w, resp)
		return
	}
	dat := []byte(r.FormValue("dat"))
	ioutil.WriteFile(path+"/main.go", dat, 0644)
	var cmd *exec.Cmd
	if r.FormValue("imp") == "true" {
		cmd = exec.Command("goimports", path+"/main.go")
	} else {
		cmd = exec.Command("gofmt", path+"/main.go")
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("main.go >> format >> cmd.CombinedOutput() >> %v\n", err)
		resp["error"] = true
		resp["msg"] = fmt.Sprintf("Server Error. Please try again.\n%s", b)
		ajaxResponse(w, resp)
		return
	}
	resp["error"] = false
	resp["output"] = fmt.Sprintf("%s", b)
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
	fmt.Fprintf(w, msg)
}
