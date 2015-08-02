package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
)

func main() {
	name := flag.String("name", "app", "the package name")
	version := flag.String("version", "", "the version name")
	licenseFile := flag.String("license", "", "the filename for the license")
	output := flag.String("output", "built.js", "the destination file")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("no input files")
	}

	var res bytes.Buffer

	if *version != "" {
		res.WriteString(CommentOut(*name+" version "+*version) + "\n")
	}

	if *version != "" && *licenseFile != "" {
		res.WriteString("//\n")
	}

	if *licenseFile != "" {
		if licenseData, err := ioutil.ReadFile(*licenseFile); err != nil {
			log.Fatal(err)
		} else {
			res.WriteString(CommentOut(string(licenseData)))
			res.WriteString("\n")
		}
	}

	res.WriteString("(function() {\n\n")
	res.WriteString(IndentCode("  ", GenerateExportsCode(*name)))
	res.WriteString("\n\n")

	scriptFiles := make([]*ScriptFile, len(flag.Args()))
	for i, file := range flag.Args() {
		var err error
		if scriptFiles[i], err = ReadScriptFile(file); err != nil {
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

	if fileData, err := JoinSourceFiles(sortedPaths); err != nil {
		log.Fatal(err)
	} else {
		res.WriteString(IndentCode("  ", string(fileData)))
	}

	res.WriteString("\n}\n")

	finishedCode := []byte(CleanEmptyLines(res.String()))
	if err := ioutil.WriteFile(*output, finishedCode, 0755); err != nil {
		log.Fatal(err)
	}

	log.Print("done!")
}

// GenerateExportsCode creates the code which makes an "exports" variable.
func GenerateExportsCode(packageName string) string {
	var res bytes.Buffer
	res.WriteString("var exports;\n")

	ifStatement := IfStatement{}
	for _, object := range []string{"self", "window"} {
		pn := PackageName(object + "." + packageName)
		condition := "'undefined' !== typeof " + object
		ifStatement.Conditions = append(ifStatement.Conditions, condition)

		code := pn.CreationCode() + "\nexports = " + object + "." + packageName + ";"
		ifStatement.Blocks = append(ifStatement.Blocks, code)
	}

	ifStatement.Conditions = append(ifStatement.Conditions, "'undefined' !== typeof module")
	ifStatement.Blocks = append(ifStatement.Blocks, "exports = module.exports;")

	res.WriteString(ifStatement.String())

	return res.String()
}

// JoinSourceFiles reads source files from paths and joins them together.
func JoinSourceFiles(sourceFiles []string) (string, error) {
	var buf bytes.Buffer

	for _, filePath := range sourceFiles {
		if fileData, err := ioutil.ReadFile(filePath); err != nil {
			return "", err
		} else {
			buf.Write(fileData)
			buf.WriteString("\n\n")
		}
	}

	return buf.String(), nil
}
