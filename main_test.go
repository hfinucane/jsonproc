package main

import (
	"fmt"
	"os"
	"testing"
)

func TestSunnyFiles(t *testing.T) {
	for _, sunnyDay := range [...]string{"loadavg", "/self/cmdline", "//version"} {
		procResult := readProcPath(sunnyDay)

		if procResult.Err != "" || procResult.Contents == nil {
			t.Errorf("Expected a result, with no errors, got %v", procResult)
		}
	}
}

func TestSunnyDirectories(t *testing.T) {
	for _, sunnyDay := range [...]string{"/", "", "/self", "/self/"} {
		procResult := readProcPath(sunnyDay)

		if procResult.Err != "" || (len(procResult.Files) == 0 && len(procResult.Dirs) == 0) {
			t.Errorf("Expected a result, with no errors, got %v", procResult)
		}
	}
}

func TestMissingPath(t *testing.T) {
	procResult := readProcPath("x")

	if procResult.Err == "" {
		t.Errorf("Expected an error, got %v", procResult)
	}
	if procResult.Contents != nil {
		t.Errorf("Expected no contents, got %v", procResult)
	}
}

func TestNoPermissions(t *testing.T) {
	procResult := readProcPath("kpagecount")

	if procResult.Err == "" {
		t.Errorf("Expected an error, got %v", procResult)
	}
	if procResult.Contents != nil {
		t.Errorf("Expected no contents, got %v", procResult)
	}
}

func TestNoTraversal(t *testing.T) {
	procResult := readProcPath("../etc/passwd")

	if procResult.Err == "" {
		t.Errorf("Expected an error, got %v", procResult)
	}
	if procResult.Contents != nil {
		t.Errorf("Expected no contents, got %v", procResult)
	}
}

func TestNoSymlinkTraversal(t *testing.T) {
	procResult := readProcPath(fmt.Sprintf("/self/%d/tasks/cwd/main_test.go", os.Getpid()))

	if procResult.Err == "" {
		t.Errorf("Expected an error, got %v", procResult)
	}
	if procResult.Contents != nil {
		t.Errorf("Expected no contents, got %v", procResult)
	}
}
