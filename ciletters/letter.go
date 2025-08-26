//go:build !solution

package ciletters

import (
	"strings"
	"text/template"

	_ "embed"
)

//go:embed template.txt
var fileStr string

func parseLog(log string) []string {
	lines := strings.Split(log, "\n")

	var lineIndex int

	l := len(lines)

	if l > 10 {
		lineIndex = l - 10
	} else {
		lineIndex = 0
	}

	return lines[lineIndex:]
}

func MakeLetter(n *Notification) (string, error) {
	var err error
	var s string

	funcMap := template.FuncMap{
		"parselog": parseLog,
	}

	tmpl, err := template.New("letter").Funcs(funcMap).Parse(fileStr)

	if err == nil {
		var b strings.Builder

		err = tmpl.Execute(&b, n)

		if err == nil {
			s = b.String()
		}
	}

	return s, err
}

//Для тестирования
// notification := Notification{
// 	Project: GitlabProject{
// 		GroupID: "go-spring-2021",
// 		ID:      "gopher",
// 	},
// 	Branch: "master",
// 	Commit: Commit{
// 		Hash:    "8967153efffffff",
// 		Message: "Solve urlfetch",
// 		Author:  "gopher",
// 	},
// 	Pipeline: Pipeline{
// 		ID: 194613,
// 	},
// }
