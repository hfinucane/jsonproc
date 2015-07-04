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

type blob struct {
	Path     string   `json:"path"`
	Files    []string `json:"files,omitempty"`
	Dirs     []string `json:"dirs,omitempty"`
	Contents *string  `json:"contents,omitempty"`
	Err      string   `json:"err,omitempty"`
	Mode     string   `json:"mode,omitempty"`
}

var BUFMAX = 1024 * 1024 * 4
var DIRMAX = 1024

func readFile(path string) (rval *blob) {
	rval = new(blob)

	fd, err := os.Open(path)

	if err != nil {
		rval.Err = err.Error()
		return
	}

	buf := make([]byte, BUFMAX)
	length, err := fd.Read(buf)
	if err != nil {
		rval.Err = err.Error()
		return
	}
	contents := string(buf[:length])
	rval.Contents = &contents
	return
}

func readDir(path string) (rval *blob) {
	rval = new(blob)

	fd, err := os.Open(path)

	if err != nil {
		fmt.Println("not a valid path", err)
		rval.Err = err.Error()
		return
	}
	files, err := fd.Readdir(DIRMAX)
	if err != nil {
		rval.Err = err.Error()
		return
	}

	// Not ideal
	rval.Files = make([]string, 0, len(files))
	rval.Dirs = make([]string, 0, len(files))
	for _, literal := range files {
		if literal.IsDir() {
			rval.Dirs = append(rval.Dirs, literal.Name())
		} else {
			rval.Files = append(rval.Files, literal.Name())
		}
	}
	return
}

func readPath(path string) (rval *blob) {
	rval = &blob{Path: path}
	fileinfo, err := os.Stat(path)

	if err != nil {
		rval.Err = err.Error()
		return
	}
	rval.Mode = fileinfo.Mode().String()

	if fileinfo.Mode().IsRegular() {
		latest := readFile(path)
		rval.Contents = latest.Contents
		rval.Err = latest.Err
	} else if fileinfo.Mode().IsDir() {
		latest := readDir(path)
		rval.Files = latest.Files
		rval.Dirs = latest.Dirs
		rval.Err = latest.Err
	}
	return
}

func main() {
	listen := flag.String("listen", ":9234", "What to listen on- you should prefer to bind to a local interface, like 10.0.1.3:9234")

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := readPath(path.Join("/proc", r.URL.Path))
		b_str, err := json.Marshal(*b)
		if err != nil {
			log.Println("marshalling error", err, b)
		}
		if err != nil || b.Err != "" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, string(b_str))
	})

	log.Fatal(http.ListenAndServe(*listen, nil))
}
