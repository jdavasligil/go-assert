package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
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
		panic("FATAL: path is not a directory.")
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
		restoreGoFilesLocal(*pathPtr)
	} else {
		fmt.Printf("Deleting assertions%s from %s%s\n", recurseStr, *pathPtr, backupStr)
		backupGoFilesLocal(*pathPtr)
		// either delete asserts recursively or just the target directory
		deleteAssertionsLocal(*pathPtr)
		// if backup is false, delete backups
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
				panic("FATAL: No backup files found.")
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

// backupGoFilesLocal renames go files by appending .bak
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
				panic("Backup files found. Please clean them first.")
			}
			err = os.Rename(filepath, filepath+".bak")
			if err != nil {
				panic(err)
			}
		}
	}
}

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

	// Look for "import".
	searchCount := 0
	line := ""
	singleImport := true
	for lineScanner.Scan() {
		line = lineScanner.Text()
		if strings.Compare(line[:6], "import") == 0 {
			break
		}
		searchCount++
		// If not found then exit early
		if searchCount >= 64 {
			return
		}
	}
	if strings.Compare(line[7:13], "assert") != 0 {
		outfile.WriteString(line + "\n")
		singleImport = false
	}
	if !singleImport {
		for lineScanner.Scan() {
			line = lineScanner.Text()
			fmt.Println(line[31:37])
			if strings.Compare(line[31:37], "assert") == 0 {
				break
			}
			outfile.WriteString(line + "\n")
			if line[0] == ')' {
				break
			}
		}
	}
	// Scan for assert.Asserts
	for lineScanner.Scan() {
		line = strings.TrimSpace(lineScanner.Text())
		if strings.Compare(line[:6], "assert") != 0 {
			outfile.WriteString(lineScanner.Text() + "\n")
		}
	}
}

// find and ignore "assert.Assert(...)\n"
//func deleteAssertions(filepath string) {
//	infile, err := os.Open(filepath)
//	if err != nil {
//		panic(err)
//	}
//	outfile, err := os.Create(filepath[:len(filepath)-4])
//	if err != nil {
//		panic(err)
//	}
//	defer infile.Close()
//	defer outfile.Close()
//
//	reader := bufio.NewReader(infile)
//	var b byte
//	idx := 0
//	for {
//		b, err = reader.ReadByte()
//		if err != nil && !errors.Is(err, io.EOF) {
//			panic(err)
//		}
//		if errors.Is(err, io.EOF) {
//			break
//		}
//		if b == byte('a') {
//			next, err := reader.Peek(len(pattern))
//			if err != nil {
//				if errors.Is(err, io.EOF) {
//					break
//				} else {
//					panic(err)
//				}
//			}
//			// If pattern is matched then an assert has been found.
//			// Go into 'skip over' mode and keep scanning bytes until
//			// closing parenthesis is found (using counter).
//			// Skip over newline, then break loop.
//			if bytes.Compare(next, pattern) == 0 {
//				pcount := 1
//				_, err = reader.ReadBytes('(')
//				if err != nil {
//					panic(err)
//				}
//				for pcount != 0 {
//					b, err := reader.ReadByte()
//					if err != nil {
//						panic(err)
//					}
//					if b == byte('(') {
//						pcount++
//					}
//					if b == byte(')') {
//						pcount--
//					}
//				}
//				_, err := reader.ReadBytes('\n')
//				if err != nil && !errors.Is(err, io.EOF) {
//					panic(err)
//				}
//				if errors.Is(err, io.EOF) {
//					break
//				}
//				continue
//			}
//		}
//		buf[idx] = b
//		idx++
//		if idx == BUF_SIZE {
//			idx = 0
//			_, err := outfile.Write(buf)
//			if err != nil && !errors.Is(err, io.EOF) {
//				panic(err)
//			}
//			if errors.Is(err, io.EOF) {
//				break
//			}
//		}
//		//fmt.Printf("%c", b)
//	}
//	outfile.Write(buf[:idx])
//}

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
