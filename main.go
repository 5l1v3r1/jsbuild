package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	name := flag.String("name", "package", "the package name")
	version := flag.String("version", "", "the version name")
	licenseFile := flag.String("license", "", "the filename for the license")
	output := flag.String("output", "built.js", "the destination file")
	flag.Parse()

	var data []byte
	if *licenseFile != "" {
		if licenseData, err := ioutil.ReadFile(*licenseFile); err != nil {
			log.Fatal(err)
		} else {
			licenseStr := string(licenseData)
			licenseLines := strings.Split(licenseStr, "\n")
			for i, line := range licenseLines {
				licenseLines[i] = "// " + line
			}
			data = append([]byte(strings.Join(licenseLines, "\n")+"\n"), data...)
		}
	}
	if *version != "" {
		data = append([]byte("// "+*name+" version "+*version+"\n\n"), data...)
	}

	if body, err := generateCodeBody(*name, *version, *licenseFile, flag.Args()); err != nil {
		log.Fatal(err)
	} else {
		data = append(data, body...)
	}

	if err := ioutil.WriteFile(*output, data, 0755); err != nil {
		log.Fatal(err)
	}

	log.Print("done!")
}

func generateCodeBody(name, version, licenseFile string, files []string) ([]byte, error) {
	scriptFiles := make([]*ScriptFile, len(files))
	for i, file := range files {
		var err error
		if scriptFiles[i], err = ReadScriptFile(file); err != nil {
			return nil, err
		}
	}
	depGraph, err := NewDepGraph(scriptFiles)
	if err != nil {
		return nil, err
	}
	sortedPaths, err := depGraph.TopologicalSort()
	if err != nil {
		return nil, err
	}

	data := []byte("(function() {\n")

	// TODO: write exports code.

	for _, filePath := range sortedPaths {
		if fileData, err := ioutil.ReadFile(filePath); err != nil {
			return nil, err
		} else {
			lines := strings.Split(string(fileData), "\n")
			for i, line := range lines {
				line = "  " + line
				if strings.TrimSpace(line) == "" {
					line = ""
				}
				lines[i] = line
			}
			data = append(data, []byte(strings.Join(lines, "\n"))...)
		}
	}

	return append(data, []byte("\n})();\n")...), nil
}
