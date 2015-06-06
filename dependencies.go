package main

import (
	"ioutil"
	"os"
	"path/filepath"
)

const dependencyKeyword = "dependency "

type ScriptFile struct {
	Path         string
	Dependencies []string
}

func ReadScriptFile(path string) (*ScriptFile, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	
	deps, err := parseDependencyComments(contents)
	if err != nil {
		return nil, err
	}
	
	res := &ScriptFile{absPath, make([]string, len(deps))}
	for i, dep := range deps {
		if res.Dependencies[i], err = filepath.Abs(dep); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func parseDependencyComments(script []byte) ([]string, error) {
	scriptStr := string(script)
	lines := strings.Split(string(script), "\n")
	result := []string{}
	for i, untrimmedLine := range lines {
		line := strings.TrimSpace(untrimmedLine)
		if strings.HasPrefix(line, dependencyKeyword) {
			result = append(result, line[len(dependencyKeyword):])
		}
	}
	return result, nil
}
