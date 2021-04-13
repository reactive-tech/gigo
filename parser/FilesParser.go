/*
Copyright 2020 Reactive Tech Limited.
"Reactive Tech Limited" is a limited company with number 09234118 and located in England, United Kingdom.
https://www.reactive-tech.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FilesParser struct {
	config                       Config
	filesRegistry                map[string]*fileToParse
	regExpFilePathInTagInclude   *regexp.Regexp
	regExpFilePathInTagIncludeIn *regexp.Regexp
	yamlFilePath                 string
}

type fileToParse struct {
	path                   string
	info                   os.FileInfo
	contents               string
	filePathsInTagInclude  []string
	filePathInTagIncludeIn string
}

func ParseFiles(config Config, yamlFilePath string) {

	filesParser := FilesParser{
		config:                       config,
		filesRegistry:                make(map[string]*fileToParse),
		regExpFilePathInTagInclude:   regexp.MustCompile(`<gigo-include file="(.*)"`),
		regExpFilePathInTagIncludeIn: regexp.MustCompile(`<gigo-include-in file="(.*)"`),
		yamlFilePath:                 yamlFilePath,
	}

	filesParser.loadFiles()
	filesParser.replaceInFilesContentsTagInclude()
	filesParser.replaceInFilesContentsTagIncludeIn()
	filesParser.copyContentsToOutputDir()
}

func (r *FilesParser) loadFiles()  {

	_ = filepath.Walk(r.config.SourceDir, func(filePath string, fileInfo os.FileInfo, err error) error {

		if r.config.SourceDir == "." {
			filePath = "./" + filePath
		}

		if r.isDir(fileInfo) ||
			r.doesMatchPatternForDirsToExclude(fileInfo, filePath) ||
			r.doesMatchPatternForFilesToExclude(fileInfo.Name()) ||
			!r.doesMatchPatternForFilesToParse(fileInfo.Name()) ||
			r.isFileLocatedInOutputDir(filePath) {
			return nil
		}

		fileContents := r.readFile(filePath)

		r.filesRegistry[filePath] = &fileToParse{
			path:                  filePath,
			info:                  fileInfo,
			contents:              fileContents,
			filePathsInTagInclude: r.extractFilePathsInTagInclude(fileContents),
			filePathInTagIncludeIn: r.extractFilePathInTagIncludeIn(fileContents),
		}

		return nil
	})
}

func (r *FilesParser) doesMatchPatternForDirsToExclude(fileInfo os.FileInfo, filePath string) bool  {

	dirPath := filePath
	if !r.isDir(fileInfo) {
		dirPath = filepath.Dir(filePath)
	}

	patterns := strings.Split(r.config.PatternDirsToExclude, " ")

	for _, pattern := range patterns {
		if matches := strings.Contains(dirPath, pattern); matches {
			return true
		}
	}

	return false
}

func (r *FilesParser) doesMatchPatternForFilesToExclude(fileName string) bool  {
	return r.doesFileMatchesPatterns(fileName, r.config.PatternFilesToExclude)
}

func (r *FilesParser) doesMatchPatternForFilesToParse(fileName string) bool  {
	return r.doesFileMatchesPatterns(fileName, r.config.PatternFilesToParse)
}

func (r *FilesParser) doesFileMatchesPatterns(fileName string, pattern string) bool  {

	patterns := strings.Split(pattern, " ")
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return matched
		}
	}

	return false
}

func (r *FilesParser) isDir(fileInfo os.FileInfo) bool {
	return fileInfo.IsDir()
}

func (r *FilesParser) createDir(dirPath string) {
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		log.Fatalf("Unable to create the folder: '%s'. Given error: %v", dirPath, err)
	}
}

func (r *FilesParser) extractFilePathsInTagInclude(fileContents string) []string {

	var gigoFilePathsToInclude []string

	submatchall := r.regExpFilePathInTagInclude.FindAllStringSubmatch(fileContents, -1)
	for _, element := range submatchall {
		gigoFilePathsToInclude = append(gigoFilePathsToInclude, r.config.SourceDir + "/" + element[1])
	}

	return gigoFilePathsToInclude
}

func (r *FilesParser) extractFilePathInTagIncludeIn(fileContents string) string {

	submatchall := r.regExpFilePathInTagIncludeIn.FindAllStringSubmatch(fileContents, -1)
	for _, element := range submatchall {
		return r.config.SourceDir + "/" + element[1]
	}

	return ""
}

func (r *FilesParser) readFile(filePath string) string {
	read, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(read)
}

func (r *FilesParser) replaceInFilesContentsTagInclude() {
	for _, file := range r.filesRegistry {
		r.replaceInFileContentsTagInclude(file)
	}

	println("-----------------")
}

func (r *FilesParser) replaceInFileContentsTagInclude(file *fileToParse) {

	for _, gigoFilePathToInclude := range file.filePathsInTagInclude {

		filePath := strings.Replace(gigoFilePathToInclude, r.config.SourceDir + "/", "", -1)
		gigoIncludeTagName := "<gigo-include file=\""+filePath+"\" />"
		includeFile := r.filesRegistry[gigoFilePathToInclude]

		if includeFile == nil {
			log.Fatal("This include file does not exist: '" + gigoFilePathToInclude + "'. It was referenced in the file: '" + file.path + "'")
		}

		if !strings.Contains(file.contents, gigoIncludeTagName) {
			continue
		}

		println("Generating: '" + file.path + "' by appending '" + gigoFilePathToInclude + "'")

		if strings.Contains(includeFile.contents, "gigo-include") {
				r.replaceInFileContentsTagInclude(includeFile)
		}

		file.contents = strings.Replace(file.contents, gigoIncludeTagName, includeFile.contents, -1)
	}
}

func (r *FilesParser) replaceInFilesContentsTagIncludeIn() {

	for _, file := range r.filesRegistry {

		if file.filePathInTagIncludeIn == "" {
			continue
		}

		filePathInTagIncludeIn := strings.Replace(file.filePathInTagIncludeIn, r.config.SourceDir + "/", "", -1)
		gigoIncludeInTagName := "<gigo-include-in file=\""+filePathInTagIncludeIn+"\" />"
		file.contents = strings.Replace(file.contents, gigoIncludeInTagName, "", -1)

		includeInContents := r.filesRegistry[file.filePathInTagIncludeIn]
		file.contents = strings.Replace(includeInContents.contents, "<gigo-include-in-content />", file.contents, -1)

		println("Generating: '" + file.path + "' by appending its contents in '" + includeInContents.path + "'")
	}

	println("-----------------")
}

func (r *FilesParser) copyContentsToOutputDir() {
	r.createDir(r.config.OutputDir)
	r.copyFilesWithMatchingPatternsToOutputDir()
	r.copyFilesWithoutMatchingPatternsToOutputDir()
}

func (r *FilesParser) copyFilesWithMatchingPatternsToOutputDir() {

	for _, file := range r.filesRegistry {

		if r.isIncludeFile(file) {
			continue
		}

		filePath := strings.Replace(file.path, r.config.SourceDir + "/", r.config.OutputDir + "/", -1)

		println("Outputting generated file: '" + filePath + "'")

		r.createDir(filepath.Dir(filePath))

		err := ioutil.WriteFile(filePath, []byte(file.contents), 0644)
		if err != nil {
			log.Fatalf("Unable to write in file: '%s'. Given error: %v", filePath, err)
		}
	}

	println("-----------------")
}

func (r *FilesParser) copyFilesWithoutMatchingPatternsToOutputDir()  {

	_ = filepath.Walk(r.config.SourceDir, func(filePath string, fileInfo os.FileInfo, err error) error {

		if r.config.SourceDir == "." {
			filePath = "./" + filePath
		}

		outputFilePath := strings.Replace(filePath, r.config.SourceDir + "/", r.config.OutputDir + "/", -1)

		if r.isSourceDir(filePath) ||
			r.doesMatchPatternForDirsToExclude(fileInfo, filePath) ||
			r.doesMatchPatternForFilesToExclude(fileInfo.Name()) ||
			r.isFileLocatedInOutputDir(filePath) ||
			r.doesExistInRegistry(filePath) {
			return nil
		}

		if !r.isDir(fileInfo) {

			r.createDir(filepath.Dir(outputFilePath))

			println("Outputting static file: '" + outputFilePath + "'")
			r.copy(filePath, outputFilePath)
		}

		return nil
	})

}

func (r *FilesParser) isSourceDir(filePath string) bool {
	return filePath == r.config.SourceDir
}

func (r *FilesParser) isIncludeFile(fileToSearch *fileToParse) bool {
	for _, file := range r.filesRegistry {
		if doesContain(file.filePathsInTagInclude, fileToSearch.path) ||
			file.filePathInTagIncludeIn == fileToSearch.path {
			return true
		}
	}
	return false
}

func (r *FilesParser) isFileLocatedInOutputDir(filePath string) bool {
	return strings.Contains(filePath, r.config.OutputDir)
}

func doesContain(gigoIncludePaths []string, pathToSearch string) bool {
	for _, gigoIncludePath := range gigoIncludePaths {
		if gigoIncludePath == pathToSearch {
			return true
		}
	}
	return false
}

func (r *FilesParser) doesExistInRegistry(filePathToSearch string) bool {
	for _, file := range r.filesRegistry {
		if file.path == filePathToSearch {
			return true
		}
	}
	return false
}

func (r *FilesParser) copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
