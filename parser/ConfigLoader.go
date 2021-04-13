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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func LoadConfig(configYamlFilePath string) Config {

	yamlFileContents, err := ioutil.ReadFile(configYamlFilePath)
	if err != nil {
		log.Fatalf("Unable to load Yaml file: '%s'. Given error: %v", configYamlFilePath, err)
	}

	config := Config{}

	err = yaml.Unmarshal(yamlFileContents, &config)
	if err != nil {
		log.Fatalf("Unable to unmarshall Yaml file: '%s'. Given error: %v", configYamlFilePath, err)
	}

	validateConfig(configYamlFilePath, &config)

	return config
}

func validateConfig(configYamlFilePath string, config *Config) {

	if config.PatternFilesToParse == "" {
		config.PatternFilesToParse = DefaultConfig.PatternFilesToParse
	}

	if config.SourceDir == "" {
		config.SourceDir = DefaultConfig.SourceDir

	} else {
		config.SourceDir = removeSlashIfInPathLastCharacter(config.SourceDir)
	}

	if config.PatternFilesToExclude == "" {
		config.PatternFilesToExclude = DefaultConfig.PatternFilesToExclude
	}

	configFileInfo, _ := os.Stat(configYamlFilePath)
	config.PatternFilesToExclude += " " + configFileInfo.Name()

	if config.PatternDirsToExclude == "" {
		config.PatternDirsToExclude = DefaultConfig.PatternDirsToExclude
	}

	if config.OutputDir == "" {
		config.OutputDir = DefaultConfig.OutputDir

	} else {
		config.OutputDir = removeSlashIfInPathLastCharacter(config.OutputDir)
	}

	if config.OutputDir == config.SourceDir {
		log.Fatalf("In the given config Yaml file: '%s' has the value of OutputDir equal to SourceDir. This is not possible since this script risks overriding your source files in SourceDir.", configYamlFilePath)
	}
}

func removeSlashIfInPathLastCharacter(path string) string {
	length := len(path)
	if path[length-1] == '/' {
		return path[:length-1]
	}

	return path
}
