//go:build !solution

package ciletters

import (
	"bufio"
	_ "embed"
	"strings"
	"text/template"
)

func get8first(str string) string {
	return str[:8]
}

func get10last(str string) string {
	scanner := bufio.NewScanner(strings.NewReader(str))
	buf := strings.Builder{}
	arr := make([]string, 0)
	for scanner.Scan() {
		tmp := scanner.Text()
		arr = append(arr, tmp)
	}
	for _, s := range arr[max(0, len(arr)-10):] {
		buf.WriteString("            ")
		buf.WriteString(s)
		buf.WriteRune('\n')
	}
	return buf.String()
}

//go:embed temp_pipeline.txt
var format string

func MakeLetter(n *Notification) (ans string, err error) {

	funcMap := template.FuncMap{
		"first8bytes": get8first,
		"last10lines": get10last,
	}

	tmpl, _ := template.New("test").Funcs(funcMap).Parse(format)
	buf := strings.Builder{}
	err = tmpl.Execute(&buf, n)
	ans = buf.String()
	return ans, err
}
