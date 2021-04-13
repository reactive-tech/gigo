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

package main

import (
	"my.domain/gigo/parser"
	"os"
)

func main() {

	println("-----------------")
	println("Generating a static-website.")

	args := os.Args[1:]

	config := parser.DefaultConfig
	configYamlFilePath := ""

	if len(args) == 1 {
		configYamlFilePath = args[0]
		println("Using the config file: '" + configYamlFilePath + "'.")
		config = parser.LoadConfig(configYamlFilePath)
	}

	println("Parsing all files with patterns: '" + config.PatternFilesToParse + "' in the source folder: '" + config.SourceDir + "'.")
	println("We will copy the generated files in the output folder: '" + config.OutputDir + "'.")
	println("Files to parse : '" + config.PatternFilesToParse + "'.")
	println("Files to exclude : '" + config.PatternFilesToExclude + "'.")
	println("Dirs to exclude : '" + config.PatternDirsToExclude + "'.")

	println("-----------------")

	parser.ParseFiles(config, configYamlFilePath)
}
