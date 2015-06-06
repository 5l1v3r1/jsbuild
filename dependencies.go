package main

import (
	"ioutil"
	"os"
	"path/filepath"
)

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
	// TODO: this
	return nil, nil
}
