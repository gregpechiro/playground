package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"time"
)

func runCmd(timeout int, command string, args ...string) (string, error) {

	// instantiate new command
	cmd := exec.Command(command, args...)

	// get pipe to standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "cmd.StdoutPipe()", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "cmd.StderrPipe()", err
	}

	// start process via command
	if err := cmd.Start(); err != nil {
		return "cmd.Start()", err
	}

	// setup a buffer to capture standard output
	var out bytes.Buffer

	var er bytes.Buffer

	// create a channel to capture any errors from wait
	done := make(chan error)
	go func() {
		if _, err := out.ReadFrom(stdout); err != nil {
			panic("buf.Read(stdout) error: " + err.Error())
		}
		if _, err := er.ReadFrom(stderr); err != nil {
			panic("buf.Read(stderr) error: " + err.Error())
		}
		done <- cmd.Wait()
	}()

	// block on select, and switch based on actions received
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return "failed to kill ", TimeOutErr
		}
		return "timeout reached, process killed", TimeOutErr
	case err := <-done:
		if err != nil {
			close(done)
			return er.String(), ExecErr
		}
		return out.String(), nil
	}
}

var ExecErr = errors.New("Execution Error")

var TimeOutErr = errors.New("Time Out")

func pretty(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.ToUpper(string(s[0])) + s[1:]
	for i := 0; i < len(s); i++ {
		if s[i] == byte(' ') {
			s = s[:i+1] + strings.ToUpper(string(s[i+1])) + s[i+2:]
		}
	}
	return s
}

var themes = []string{
	"ambiance",
	"chaos",
	"chrome",
	"clouds",
	"clouds_midnight",
	"cobalt",
	"crimson_editor",
	"dawn",
	"dreamweaver",
	"eclipse",
	"github",
	"gruvbox",
	"idle_fingers",
	"iplastic",
	"katzenmilch",
	"kr_theme",
	"kuroir",
	"merbivore",
	"merbivore_soft",
	"mono_industrial",
	"monokai",
	"pastel_on_dark",
	"solarized_dark",
	"solarized_light",
	"sqlserver",
	"terminal",
	"textmate",
	"tomorrow",
	"tomorrow_night",
	"tomorrow_night_blue",
	"tomorrow_night_bright",
	"tomorrow_night_eighties",
	"twilight",
	"vibrant_ink",
	"xcode",
}
