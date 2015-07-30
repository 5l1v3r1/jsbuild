package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const dependencyKeyword = "//deps "

// A ScriptFile holds both the clean absolute path to a script and the clean absolute paths its
// dependencies.
type ScriptFile struct {
	Path         string
	Dependencies []string
}

// ReadScriptFile generates a script file and reads the file's dependencies from its heading
// comment.
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

// NewDepGraph creates a dependency graph from an array of *ScriptFiles.
// This will fail if a dependency is listed which is not included in the script files list.
func NewDepGraph(scriptFiles []*ScriptFile) (*DepGraph, error) {
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

// TopologicalSort sorts the files in this graph and returns an ordered list of file paths. A file
// path will always appear in the list before all of its dependents.
// If the graph cannot be topologically sorted, an error will be returned.
// The receiving DepGraph will be destroyed in the process of topological sorting.
func (d *DepGraph) TopologicalSort() ([]string, error) {
	// This is an implementation of Kahn's topological sorting algorithm.

	bottomNodes := make([]*depGraphNode, 0, len(d.nodes))
	for _, node := range d.nodes {
		for _, edge := range node.edges {
			if edge.dependent == node {
				continue OuterLoop
			}
		}
		bottomNodes = append(bottomNodes, node)
	}

	sources := make([]string, 0, len(d.nodes))
	for len(bottomNodes) > 0 {
		node := bottomNodes[len(bottomNodes)-1]
		sources = append(sources, node.path)
		bottomNodes = bottomNodes[0 : len(bottomNodes)-1]
		edges := node.edges
		d.removeNode(node)
		for _, edge := range edges {
			if len(edge.dependent.edges) == 0 {
				bottomNodes = append(bottomNodes, edge.dependent)
			}
		}
	}

	if len(d.nodes) > 0 {
		return nil, errors.New("the graph is not acyclic")
	} else {
		return sources, nil
	}
}

func (d *DepGraph) removeNode(node *depGraphNode) {
	for i, aNode := range d.nodes {
		if aNode == node {
			d.nodes[i] = d.nodes[len(d.nodes)-1]
			d.nodes = d.nodes[:len(d.nodes)-1]
			break
		}
	}
	for _, edge := range node.edges {
		edge.remove()
	}
}

type depGraphNode struct {
	edges []*depGraphEdge
	path  string
}

type depGraphEdge struct {
	dependent  *depGraphNode
	dependency *depGraphNode
}

func (e *depGraphEdge) remove() {
	for _, node := range []*depGraphNode{e.dependent, e.dependency} {
		for i, edge := range node.edges {
			if edge == e {
				node.edges[i] = node.edges[len(node.edges)-1]
				node.edges = node.edges[:len(node.edges)-1]
				break
			}
		}
	}
}
