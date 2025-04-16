package filesystem

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/medatechnology/goutil/simplelog"
)

// This is a simple way to get the filename from a complete path
// It will return the filename only, without the path
// Example:
// filename := filesystem.Base("/home/user/test.txt")
// It will return "test.txt"
func FileName(completefilename string) string {
	return filepath.Base(completefilename)
}

// Return the directory name (everything beside the filename)
// DirPath will return the directory name of the complete path
// Example:
// dir := filesystem.DirPath("/home/user/test.txt")
// It will return "/home/user"
func DirPath(completefilename string) string {
	return filepath.Dir(completefilename)
}

// Just like ls command, in the path, list all files in array of fs.DisEntry
// It will return an array of fs.DirEntry
// Example:
// files := filesystem.Dir("/home/user", ".txt")
// It will return all files in the directory /home/user with .txt extension
// If filter is empty or "*", it will return all files in the directory
// Output:
// [file1.txt file2.txt file3.txt]
func Dir(path, filter string) []fs.DirEntry {
	fn := "Dir"
	var filtered []fs.DirEntry
	files, err := os.ReadDir(path)
	if err != nil {
		simplelog.LogInfoAny(fn, 10, "01:cannot get files in path", path, ";error:", err)
		// log.Fatal(err)
	}
	// default to all files (but not hidden files)
	if filter == "" {
		filter = "*"
	}
	if filter != "*" {
		// search if it has suffix of the filter, this is the only functionality of the filter
		// at the moment. Only matching at the end and it's case sensitive.
		for _, f := range files {
			if strings.HasSuffix(f.Name(), filter) {
				filtered = append(filtered, f)
			}
		}
		return filtered
	} else {
		return files
	}
	// DEBUGGING:
	// for _, file := range files {
	// 	fmt.Println(file.Name(), file.IsDir())
	// }
}

// Check if file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Check if file or directory exist. if directory=true then check if exist AND also if it's a dir
// Example:
// exist := filesystem.DirFileExist("/home/user/test.txt", false)
// It will return true if the file "test.txt" exist
// exist = filesystem.DirFileExist("/home/user/test", true)
// It will return true if the directory "test" exist
func DirFileExist(path string, directory bool) bool {
	finfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}
	if directory {
		if finfo.IsDir() {
			return true
		} else {
			return false
		}
	}
	return true
}

// Only for textfile, this will read the file and return as array of string per line
// This is a simple way to read a file line by line
// It will return an array of string, each line is a string
// It will return an empty array if the file is not found or error
// Example:
// content := filesystem.More("test.txt")
func More(filename string) []string {
	fn := "More"
	file, err := os.Open(filename)
	if err != nil {
		simplelog.LogInfoAny(fn, 10, "01:cannot read file ", filename, ";error:", err)
		// log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content []string
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		content = append(content, scanner.Text())
		// fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		// log.Fatal(err)
		simplelog.LogInfoAny(fn, 10, "02:error read file ", filename, ";error:", err)
	}
	return content
}
