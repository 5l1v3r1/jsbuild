package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	/*name := flag.String("name", "package", "the package name")
	version := flag.String("version", "", "the version name")
	licenseFile := flag.String("license", "", "the filename for the license")
	output := flag.String("output", "built.js", "the destination file")
	*/
	flag.Parse()

	inputFiles := flag.Args()
	scriptFiles := make([]*ScriptFile, len(inputFiles))
	for i, file := range inputFiles {
		var err error
		scriptFiles[i], err = ReadScriptFile(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	depGraph, err := NewDepGraph(scriptFiles)
	if err != nil {
		log.Fatal(err)
	}
	sortedPaths, err := depGraph.TopologicalSort()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sorted files", sortedPaths)
}
