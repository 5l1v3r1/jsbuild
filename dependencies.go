package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const dependencyKeyword = "//deps "

type ScriptFile struct {
	Path         string
	Dependencies []string
}

func ReadScriptFile(path string) (*ScriptFile, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	absPath = filepath.Clean(absPath)

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
		if filepath.IsAbs(dep) {
			res.Dependencies[i] = filepath.Clean(dep)
		} else {
			res.Dependencies[i] = filepath.Clean(filepath.Join(filepath.Dir(absPath), dep))
		}
	}
	return res, nil
}

func parseDependencyComments(script []byte) ([]string, error) {
	lines := strings.Split(string(script), "\n")
	result := []string{}
	for _, untrimmedLine := range lines {
		line := strings.TrimSpace(untrimmedLine)
		if strings.HasPrefix(line, dependencyKeyword) {
			result = append(result, strings.Split(line[len(dependencyKeyword):], " ")...)
		}
	}
	return result, nil
}
