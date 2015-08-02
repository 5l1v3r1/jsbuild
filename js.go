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
func (f IfStatement) String(indent string) string {
	buf := &bytes.Buffer{}
	for i, condition := range f.Conditions {
		buf.WriteString(indent)
		if i != 0 {
			buf.WriteString("} else ")
		}
		buf.WriteString("if (")
		buf.WriteString(condition)
		buf.WriteString(") {\n")
		buf.WriteString(IndentCode(indent+"  ", strings.TrimSpace(f.Blocks[i])))
		buf.WriteString("\n")
	}
	buf.WriteString(indent)
	buf.WriteString("}")
	return buf.String()
}

// IndentCode adds an indent before every line in a code block.
func IndentCode(indent, code string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}
