package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

type blob struct {
	Path        string
	files, dirs *[]string
	Contents    *string
	Err         string
	Mode        string
}

var BUFMAX = 1024 * 1024 * 4

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

func readDir(path string) *blob {
	return &blob{}
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
		return readDir(path)
	}
	return
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := readPath(path.Join("/proc", r.URL.Path))
		b_str, err := json.Marshal(*b)
		if err != nil {
			log.Println("marshalling error", err, b)
		}
		fmt.Fprintf(w, string(b_str))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
