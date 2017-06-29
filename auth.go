package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/cagnosolutions/web"
)

func init() {
	mux.AddRoutes(login, logout)
	mux.AddSecureRoutes(ADMIN, goGet, setCurrent)
}

var ADMIN = web.Auth{
	Roles:    []string{"ADMIN"},
	Msg:      "You are not authorized",
	Redirect: "/",
}

var login = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") != "playground" || r.FormValue("password") != "administrator" {
		http.Redirect(w, r, "/", 303)
		return
	}
	web.Login(w, r, "ADMIN")
	web.SetSuccessRedirect(w, r, "/", "Logged in")
	return
}}

var goGet = web.Route{"POST", "/goget", func(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("url")
	if !strings.HasPrefix(path, "github.com/cagnosolutions/") && !strings.HasPrefix(path, "github.com/gregpechiro/") && !strings.HasPrefix(path, "github.com/scottcagno/") {
		web.SetErrorRedirect(w, r, "/", "Error getting "+r.FormValue("url"))
		return
	}

	cmd := exec.Command(DIR+"env/versions/current/bin/go", "get", "-u", r.FormValue("url"))
	cmd.Env = []string{"PATH=/usr/bin", "GOROOT=" + DIR + "env/versions/current", "GOPATH=" + DIR + "env/libs"}
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("main.go >> run >> exec.Command.CombinedOutput(\"go\" \"get\", \"-u\", %s) >> %v\n\n", err, r.FormValue("url"))
		log.Printf("GO GET -- %s\n", b)
		web.SetErrorRedirect(w, r, "/", "Error getting "+r.FormValue("url"))
		return
	}
	web.SetSuccessRedirect(w, r, "/", "Got "+r.FormValue("url"))
	return
}}

var setCurrent = web.Route{"POST", "/setCurrent", func(w http.ResponseWriter, r *http.Request) {
	if err := SetCurrent(); err != nil {
		AjaxResponse(w, map[string]interface{}{
			"msg": fmt.Sprintf("Error setting current:\n%s\n", err.Error()),
		})
		return
	}
	AjaxResponse(w, map[string]interface{}{
		"msg": fmt.Sprintf("Successfully set current:\n%s\n", Current),
	})
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	http.Redirect(w, r, "/", 303)
	return
}}
