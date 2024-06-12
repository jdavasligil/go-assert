package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func main() {
    dir, err := os.Getwd(); if err != nil {
        panic(err)
    }
    pathPtr := flag.String("path", dir, "Path to begin deleting assertions.")
    backupPtr := flag.Bool("b", false, "Backup original source files by appending '.bak' to filenames.")
    recursePtr := flag.Bool("r", false, "Recursively delete assertions in subfolders.")

    flag.Parse()

    finfo, err := os.Stat(*pathPtr); if err != nil {
        panic(err)
    } else if !finfo.IsDir() {
        panic("FATAL: path is not a directory.")
    }

    recurseStr := ""
    if *recursePtr { recurseStr = "recursively " }
    backupStr := "..."
    if *backupPtr { backupStr = " with backups..." }

    fmt.Printf("Deleting assertions %sfrom %s%s\n", recurseStr, *pathPtr, backupStr)

    backupGoFilesLocal(*pathPtr)

    // either delete asserts recursively or just the target directory
    deleteAssertionsLocal(*pathPtr)

    // if backup is false, delete backups

    fmt.Println("Done.")
}

func backupGoFile(filepath string) {
    infile, err := os.Open(filepath); if err != nil {
        panic(err)
    }
    outfile, err := os.Create(filepath + ".bak"); if err != nil {
        panic(err)
    }
    defer infile.Close()
    defer outfile.Close()

    infile.WriteTo(outfile)
}

func backupGoFilesLocal(dir string) {
    files, err := os.ReadDir(dir); if err != nil {
        panic(err)
    }
    for _, f := range files {
        if isGoFile(f.Name()) {
            absPath := path.Join(dir, f.Name())
            backupGoFile(absPath)
        }
    }
}

// find and ignore "assert.Assert(...)\n"
func deleteAssertions(filepath string) {
    infile, err := os.Open(filepath + ".bak"); if err != nil {
        panic(err)
    }
    outfile, err := os.Create(filepath + ".test"); if err != nil {
        panic(err)
    }
    defer infile.Close()
    defer outfile.Close()

    reader := bufio.NewReader(infile)
    pattern := [5]byte{byte('s'),byte('s'),byte('e'),byte('r'),byte('t')}
    _ = pattern
    idx := 0
    _ = idx
    var b byte
    for {
        b, err = reader.ReadByte(); if !(err == nil || errors.Is(err, io.EOF)) {
            panic(err)
        } else if errors.Is(err, io.EOF) {
            break
        }
        if b == byte('a') {
            next, err := reader.Peek(5); if err != nil {
                if errors.Is(err, io.EOF) {
                    break
                } else {
                    panic(err)
                }
            }
            // If pattern is matched then an assert has been found.
            // Go into 'skip over' mode and keep scanning bytes until
            // closing parenthesis is found (using stack counter).
            // Skip over newline, then break loop.
            fmt.Printf("a%s\n", string(next))
        }
        //fmt.Printf("%c", b)
    }
}

func deleteAssertionsLocal(dir string) {
    files, err := os.ReadDir(dir); if err != nil {
        panic(err)
    }
    for _, f := range files {
        if isGoFile(f.Name()) {
            absPath := path.Join(dir, f.Name())
            deleteAssertions(absPath)
        }
    }
}

func isGoFile(f string) bool {
    return strings.HasSuffix(f, ".go")
}
