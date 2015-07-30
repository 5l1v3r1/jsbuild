package main

import (
	"errors"
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

// A DepGraph represents a graph of source files. The edges in the graph are their dependencies.
// Using a DepGraph, it is possible to perform a topological sort.
type DepGraph struct {
	nodes []*depGraphNode
}

func NewDepGraph(scriptFiles []ScriptFile) (*DepGraph, error) {
	nodes := make([]*depGraphNode, len(scriptFiles))
	nodesPerPath := map[string]*depGraphNode{}
	for i, f := range scriptFiles {
		nodes[i] = &depGraphNode{[]*depGraphEdge{}, f.Path}
		nodesPerPath[f.Path] = nodes[i]
	}
	for _, f := range scriptFiles {
		for _, dep := range f.Dependencies {
			if depNode, ok := nodesPerPath[dep]; !ok {
				return nil, errors.New("dependency not included: " + dep)
			} else {
				thisNode := nodesPerPath[f.Path]
				edge := &depGraphEdge{thisNode, nodesPerPath[dep]}
				depNode.edges = append(depNode.edges, edge)
				thisNode.edges = append(thisNode.edges, edge)
			}
		}
	}
	return &DepGraph{nodes}, nil
}

type depGraphNode struct {
	edges []*depGraphEdge
	path  string
}

type depGraphEdge struct {
	dependent  *depGraphNode
	dependency *depGraphNode
}
