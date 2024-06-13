package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pathPtr := flag.String("path", dir, "Path to begin deleting assertions.")
	backupPtr := flag.Bool("b", false, "Backup original source files by appending '.bak' to filenames.")
	recursePtr := flag.Bool("r", false, "Recursively delete assertions in subfolders.")
	restorePtr := flag.Bool("restore", false, "Restore all backups.")

	flag.Parse()

	finfo, err := os.Stat(*pathPtr)
	if err != nil {
		panic(err)
	} else if !finfo.IsDir() {
		fmt.Printf("Path: %s is not a valid directory.\n", *pathPtr)
		return
	}

	recurseStr := ""
	if *recursePtr {
		recurseStr = " recursively"
	}
	backupStr := "..."
	if *backupPtr {
		backupStr = " with backups..."
	}

	if *restorePtr {
		fmt.Printf("Restoring backups%s...", recurseStr)
		if *recursePtr {
			restoreGoFilesRecursive(*pathPtr)
		} else {
			restoreGoFilesLocal(*pathPtr)
		}
	} else {
		fmt.Printf("Deleting assertions%s from %s%s\n", recurseStr, *pathPtr, backupStr)
		if *recursePtr {
			backupGoFilesRecursive(*pathPtr)
			deleteAssertionsRecursive(*pathPtr)
			if !(*backupPtr) {
				removeBackupGoFilesRecursive(*pathPtr)
			}
		} else {
			backupGoFilesLocal(*pathPtr)
			deleteAssertionsLocal(*pathPtr)
			if !(*backupPtr) {
				removeBackupGoFilesLocal(*pathPtr)
			}
		}
	}

	fmt.Println("Done.")
}

func restoreGoFilesLocal(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go") {
			filepath := path.Join(dir, f.Name())
			_, err = os.Stat(filepath + ".bak")
			if err != nil {
				fmt.Println("No .bak files found.")
				return
			}
			err = os.Remove(filepath)
			if err != nil {
				panic(err)
			}
			err = os.Rename(filepath+".bak", filepath)
			if err != nil {
				panic(err)
			}
		}
	}
}

func restoreGoFilesRecursive(root string) {
	filepath.WalkDir(root, fs.WalkDirFunc(func(path string, d fs.DirEntry, err error) error {
		file, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			restoreGoFilesLocal(path)
		}
		return err
	}))
}

// backupGoFilesLocal renames go files by appending .bak to the filename.
func backupGoFilesLocal(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go") {
			filepath := path.Join(dir, f.Name())
			_, err = os.Stat(filepath + ".bak")
			if err == nil {
				panic("Backup files found. Please remove them first.")
			}
			err = os.Rename(filepath, filepath+".bak")
			if err != nil {
				panic(err)
			}
		}
	}
}

func backupGoFilesRecursive(root string) {
	filepath.WalkDir(root, fs.WalkDirFunc(func(path string, d fs.DirEntry, err error) error {
		file, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			backupGoFilesLocal(path)
		}
		return err
	}))
}

// removeBackupGoFilesLocal deletes all .go.bak files in a directory.
func removeBackupGoFilesLocal(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go.bak") {
			filepath := path.Join(dir, f.Name())
			err = os.Remove(filepath)
			if err != nil {
				panic(err)
			}
		}
	}
}

func removeBackupGoFilesRecursive(root string) {
	filepath.WalkDir(root, fs.WalkDirFunc(func(path string, d fs.DirEntry, err error) error {
		file, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			removeBackupGoFilesLocal(path)
		}
		return err
	}))
}

// deleteAssertions truncates the original .go file, scans the backup file
// line by line, and only writes to the .go file if it does not contain assert
// content.
//
// It assumes the file is "well formed", i.e., a multiline import
// statement ends with a single ')' on its own line and each assertion
// only takes up a single line.
func deleteAssertions(filepath string) {
	infile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	outfile, err := os.Create(filepath[:len(filepath)-4])
	if err != nil {
		panic(err)
	}
	defer infile.Close()
	defer outfile.Close()

	lineScanner := bufio.NewScanner(infile)
	line := ""
	for lineScanner.Scan() {
		line = strings.TrimSpace(lineScanner.Text())
		if len(line) >= 8 && strings.Compare(line[:7], "import ") == 0 {
			line = line[7:]
		}
		if (len(line) >= 6 && strings.Compare(line[:6], "assert") == 0) ||
			(len(line) >= 34 && strings.Compare(line[24:33], "go-assert") == 0) {
			continue
		}
		outfile.WriteString(lineScanner.Text() + "\n")
	}
}

func deleteAssertionsLocal(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go.bak") {
			filepath := path.Join(dir, f.Name())
			deleteAssertions(filepath)
		}
	}
}

func deleteAssertionsRecursive(root string) {
	filepath.WalkDir(root, fs.WalkDirFunc(func(path string, d fs.DirEntry, err error) error {
		file, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			deleteAssertionsLocal(path)
		}
		return err
	}))
}
