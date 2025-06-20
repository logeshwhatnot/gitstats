package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"slices"
	"strings"
)

func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "./gitstats"

	return dotFile
}

func openFile(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return f
}

func parseFileLinesToSlice(filepath string) []string {
	f := openFile(filepath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return lines
}

func dumpStringSliceToFiles(repos []string, filepath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filepath, []byte(content), 0755)
}

func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !slices.Contains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func addNewSliceElementsToFile(filepath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filepath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringSliceToFiles(repos, filepath)
}

func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}

func recursiveFolderScan(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

func scan(folder string) {
	fmt.Printf("Found Folders:\n\n")
	repositories := recursiveFolderScan(folder)
	filepath := getDotFilePath()
	addNewSliceElementsToFile(filepath, repositories)
	fmt.Printf("\n\nSuccessfully added\n")
}
