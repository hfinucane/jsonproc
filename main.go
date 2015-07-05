/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

type ProcResult struct {
	Path     string   `json:"path"`
	Files    []string `json:"files,omitempty"`
	Dirs     []string `json:"dirs,omitempty"`
	Contents *string  `json:"contents,omitempty"`
	Err      string   `json:"err,omitempty"`
	Mode     string   `json:"mode,omitempty"`
}

var BUFMAX = 1024 * 1024 * 4
var DIRMAX = 1024

func readFile(path string) (contents *string, err error) {
	fd, err := os.Open(path)

	if err != nil {
		return
	}

	buf := make([]byte, BUFMAX)
	length, err := fd.Read(buf)
	if err != nil {
		return
	}

	contents = new(string)
	*contents = string(buf[:length])
	return
}

func readDir(path string) (files, dirs []string, err error) {
	fd, err := os.Open(path)

	if err != nil {
		return
	}
	direntries, err := fd.Readdir(DIRMAX)
	if err != nil {
		return
	}

	// Not ideal
	files = make([]string, 0, len(direntries))
	dirs = make([]string, 0, len(direntries))
	for _, literal := range direntries {
		if literal.IsDir() {
			dirs = append(dirs, literal.Name())
		} else {
			files = append(files, literal.Name())
		}
	}
	return
}

func readPath(path string) (rval *ProcResult) {
	rval = &ProcResult{Path: path}
	fileinfo, err := os.Stat(path)

	if err != nil {
		rval.Err = err.Error()
		return
	}
	rval.Mode = fileinfo.Mode().String()

	if fileinfo.Mode().IsRegular() {
		rval.Contents, err = readFile(path)
	} else if fileinfo.Mode().IsDir() {
		rval.Files, rval.Dirs, err = readDir(path)
	}
	if err != nil {
		rval.Err = err.Error()
	}
	return
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	b := readPath(path.Join("/proc", r.URL.Path))
	b_str, err := json.Marshal(*b)
	if err != nil {
		log.Println("marshalling error", err, b)
	}
	if err != nil || b.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, string(b_str))
}

func main() {
	listen := flag.String("listen", ":9234", "What to listen on- you should prefer to bind to a local interface, like 10.0.1.3:9234")
	flag.IntVar(&BUFMAX, "file-limit", BUFMAX, "Maximum amount of files to read")
	flag.IntVar(&DIRMAX, "dir-limit", DIRMAX, "Maximum number of directory entries to read")

	flag.Parse()

	http.HandleFunc("/", jsonHandler)

	log.Fatal(http.ListenAndServe(*listen, nil))
}
