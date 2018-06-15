package main

import (
	"flag"
	"os"
	"strings"
	"text/template"
	"bufio"
	"io"
)

const defaultParamPrefix = "GEV"
const envDelimiter = "="
const paramsListDelimiter = ","
const paramsListSuffix = "LIST"
const templateSuffix = ".tpl"

func getArgs() (string, string, string) {
	inFileName := flag.String("i", "", "Input file name")
	outFileName := flag.String("o", "", "Output file name")
	prefix := flag.String("p", defaultParamPrefix, "Prefix for environment variables")
	flag.Parse()

	if *outFileName == "" && strings.HasSuffix(*inFileName, templateSuffix) {
		*outFileName = strings.TrimSuffix(*inFileName, templateSuffix)
	}

	return *inFileName, *outFileName, *prefix
}

func getParams(prefix string, suffix string, delimiter string) map[string]interface{} {
	prefixShift := len(prefix) + 1
	suffixLen := len(suffix)
	tplParams := make(map[string]interface{})

	for _, e := range os.Environ() {
		pair := strings.Split(e, envDelimiter)
		if strings.HasPrefix(pair[0], prefix) {
			var (
				paramName string
				paramValue string
				paramList []string
			)
			if strings.HasSuffix(pair[0], suffix) {
				suffixShift := len(pair[0]) - suffixLen - 1
				paramName = pair[0][prefixShift:suffixShift]
				trimmer := strings.Join([]string{delimiter, " "}, "")
				paramList = strings.Split(strings.Trim(pair[1], trimmer), delimiter)
				tplParams[paramName] = paramList
			} else {
				paramName = pair[0][prefixShift:]
				paramValue = pair[1]
				tplParams[paramName] = paramValue
			}

		}
	}

	return tplParams
}

func main() {
	var (
		inFileName string
		outFileName string
		prefix string
		out *os.File
		tmpl *template.Template
		err error
	)

	inFileName, outFileName, prefix = getArgs()

	if inFileName != "" {
		tmpl, err = template.ParseFiles(inFileName)
		if err != nil {
			panic(err)
		}
	} else {
		var (
			content []rune
			b rune
		)
		reader := bufio.NewReader(os.Stdin)
		for {
			b, _, err = reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			content = append(content, b)
		}
		tmpl, _ = template.New("piped").Parse(string(content))
	}

	if outFileName != "" {
		out, err = os.Create(outFileName)
		if err != nil {
			panic(err)
		}
	} else {
		out = os.Stdout
	}

	tplParams := getParams(prefix, paramsListSuffix, paramsListDelimiter)

	tmpl.Execute(out, tplParams)
}
