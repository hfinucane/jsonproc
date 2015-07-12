package main

import (
	"path"
	"testing"
)

func TestSunnyFiles(t *testing.T) {
	for _, sunny_day := range [...]string{"loadavg", "/self/cmdline", "//version"} {
		proc_result := readProcPath(path.Join("/proc", sunny_day))

		if proc_result.Err != "" || proc_result.Contents == nil {
			t.Errorf("Expected a result, with no errors, got %v", proc_result)
		}
	}
}

func TestSunnyDirectories(t *testing.T) {
	for _, sunny_day := range [...]string{"/", "", "/self", "/self/"} {
		proc_result := readProcPath(path.Join("/proc", sunny_day))

		if proc_result.Err != "" || (len(proc_result.Files) == 0 && len(proc_result.Dirs) == 0) {
			t.Errorf("Expected a result, with no errors, got %v", proc_result)
		}
	}
}

func TestMissingPath(t *testing.T) {
	proc_result := readProcPath(path.Join("/proc", "x"))

	if proc_result.Err == "" {
		t.Errorf("Expected an error, got %v", proc_result)
	}
	if proc_result.Contents != nil {
		t.Errorf("Expected no contents, got %v", proc_result)
	}
}

func TestNoPermissions(t *testing.T) {
	proc_result := readProcPath(path.Join("/proc", "kpagecount"))

	if proc_result.Err == "" {
		t.Errorf("Expected an error, got %v", proc_result)
	}
	if proc_result.Contents != nil {
		t.Errorf("Expected no contents, got %v", proc_result)
	}
}
