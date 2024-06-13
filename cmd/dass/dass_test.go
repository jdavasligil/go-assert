package main

import (
	"os"
	"strings"
	"testing"
)

const max_tests = 16

type test struct {
	testfile []byte
	expected []byte
}

func TestDeleteAssertions(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testDir := wd[:len(wd)-8] + "test"
	testDirInfo, err := os.Stat(testDir)
	if err != nil {
		t.Fatal(err)
	} else if !testDirInfo.IsDir() {
		t.Fatalf("%s is not a directory.", testDir)
	}
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	tests := make([]test, 16)
	_ = tests
	testIndices := make(map[int]struct{})

	// Read test data into memory
	for _, f := range files {
		filepath := testDir + "/" + f.Name()
		fileType := f.Name()[:len(f.Name())-1]
		fileIdx := int(byte(f.Name()[len(f.Name())-1]) - 48)
		if strings.Compare(fileType, "expected") == 0 {
			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Fatal(err)
			}
			tests[fileIdx].expected = content
		}
		if strings.Compare(fileType, "testfile") == 0 {
			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Fatal(err)
			}
			tests[fileIdx].testfile = content
		}
		if tests[fileIdx].testfile != nil && tests[fileIdx].expected != nil {
			testIndices[fileIdx] = struct{}{}
		}
	}

	tempPath := testDir + "/" + "temp.go"
	for idx := range testIndices {
		os.WriteFile(tempPath, tests[idx].testfile, os.FileMode(0644))
		backupGoFilesLocal(testDir)
		deleteAssertionsLocal(testDir)
		got, err := os.ReadFile(tempPath)
		if err != nil {
			t.Fatal(err)
		}
		linecount := 0
		bytesDiffer := false
		expectedLen := len(tests[idx].expected)
		for i := range got {
			if got[i] == '\n' {
				linecount++
			}
			if i >= expectedLen || got[i] != tests[idx].expected[i] {
				bytesDiffer = true
				break
			}
		}
		if bytesDiffer {
			t.Logf("Test %d fails to match expected file on line %d.", idx, linecount)
			os.WriteFile(testDir+"/"+"testdump"+string(byte(idx+48)), got, os.FileMode(0644))
			t.Fail()
		}
		restoreGoFilesLocal(testDir)
	}
	os.Remove(tempPath)
}
