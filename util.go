package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetLanguages(path string) map[string][]string {
	languages := make(map[string][]string)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			languages[file.Name()] = GetVersions(path + "/" + file.Name())
		}
	}
	return languages
}

func GetVersions(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	var versions []string
	for _, file := range files {
		if file.IsDir() {
			versions = append(versions, file.Name())
		}
	}
	// sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	sort.SliceStable(versions, func(i, j int) bool {
		ver1, ver2 := versions[i], versions[j]
		if ver1 == ver2 {
			return true
		}

		ver1S, ver2S := strings.Split(ver1, "."), strings.Split(ver2, ".")

		ver1Slen, ver2Slen := len(ver1S), len(ver2S)

		if ver1Slen > ver2Slen {
			for k := 0; k < (ver1Slen - ver2Slen); k++ {
				ver2S = append(ver2S, "0")
			}
		}
		if ver2Slen > ver1Slen {
			for k := 0; k < (ver2Slen - ver1Slen); k++ {
				ver1S = append(ver1S, "0")
			}
		}
		ver1Slen, ver2Slen = len(ver1S), len(ver2S)
		for k := 0; k < ver1Slen; k++ {
			ver1Int, _ := strconv.Atoi(ver1S[k])
			ver2Int, _ := strconv.Atoi(ver2S[k])
			if ver1Int == ver2Int {
				continue
			}
			if ver1Int < ver2Int {
				return false
			}
			if ver1Int > ver2Int {
				return true
			}
		}
		return false

	})
	return versions
}

func IsCurrent(ver string) bool {
	return strings.Contains(Current, ver)
}

func SetCurrent() error {
	cmd := exec.Command(DIR+"env/versions/current/bin/go", "version")
	cmd.Env = []string{"GOROOT=" + DIR + "env/versions/current", "GOPATH=" + DIR + "env/libs"}
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("\nutil.go >> SetCurrent() >> cmd.CombinedOutput() >> %v\n", err)
		log.Printf("\toutput >> %s\n", b)
		return err
	}
	re := regexp.MustCompile(`[0-9]\.[0-9]*\.*[0-9]*`)
	Current = re.FindString(string(b))
	return nil
}

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

func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
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
	//"gruvbox",
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
