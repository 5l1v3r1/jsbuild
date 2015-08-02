package main

import (
	"bytes"
	"strings"
)

// An IfStatement makes it possible to format an if-else statement.
type IfStatement struct {
	Conditions []string
	Blocks     []string
}

// String returns the properly indented and formatted if-else statement.
// There will be no trailing newline.
func (f IfStatement) String() string {
	buf := &bytes.Buffer{}
	for i, condition := range f.Conditions {
		if i != 0 {
			buf.WriteString("} else ")
		}
		buf.WriteString("if (")
		buf.WriteString(condition)
		buf.WriteString(") {\n")
		buf.WriteString(IndentCode("  ", strings.TrimSpace(f.Blocks[i])))
		buf.WriteString("\n")
	}
	buf.WriteString("}")
	return buf.String()
}

// A PackageName is a "." separated namespace for a JavaScript object. For example,
// "window.app.MyClass".
type PackageName string

// CreationCode generates code which creates an object with the package name.
func (p PackageName) CreationCode() string {
	components := strings.Split(string(p), ".")
	var res bytes.Buffer
	for i := 1; i < len(components); i++ {
		var ifStatement IfStatement
		objectName := strings.Join(components[:i+1], ".")

		condition := "!" + objectName
		ifStatement.Conditions = append(ifStatement.Conditions, condition)

		code := objectName + " = " + "{};"
		ifStatement.Blocks = append(ifStatement.Blocks, code)

		res.WriteString(ifStatement.String())
		if i < len(components)-1 {
			res.WriteString("\n")
		}
	}
	return res.String()
}

// CleanEmptyLines removes whitespace from empty lines.
func CleanEmptyLines(code string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines[i] = ""
		}
	}
	return strings.Join(lines, "\n")
}

// CommentOut generates a code block which contains a commented string.
func CommentOut(code string) string {
	return IndentCode("// ", strings.TrimSpace(code))
}

// IndentCode adds an indent before every line in a code block.
func IndentCode(indent, code string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}
