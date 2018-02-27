package main

import (
	"flag"
	"os"
	"strings"
	"text/template"
)

const envDelimiter string = "="
const paramsListDelimiter string = ","
const paramsListSuffix string = "LIST"

type envContainer struct {
	Value string
	List  []string
}

func getArgs() (string, string, string) {
	inFileName := flag.String("i", "", "Input file name")
	outFileName := flag.String("o", "", "Output file name")
	prefix := flag.String("p", "", "Prefix for environment variables")
	flag.Parse()

	if len(*inFileName) == 0 {
		os.Stderr.WriteString("Input file name not present\n")
		os.Exit(1)
	}

	if len(*outFileName) == 0 {
		os.Stderr.WriteString("Output file name not present\n")
		os.Exit(1)
	}

	if len(*prefix) == 0 {
		os.Stderr.WriteString("Prefix not present\n")
		os.Exit(1)
	}

	return *inFileName, *outFileName, *prefix
}

func getParams(prefix string, suffix string, delimiter string) map[string]envContainer {
	prefixShift := len(prefix) + 1
	suffixLen := len(suffix)
	tplParams := make(map[string]envContainer)

	for _, e := range os.Environ() {
		pair := strings.Split(e, envDelimiter)
		if strings.HasPrefix(pair[0], prefix) {
			var paramName string
			var paramValue string
			var paramList []string
			if strings.HasSuffix(pair[0], suffix) {
				suffixShift := len(pair[0]) - suffixLen - 1
				paramName = pair[0][prefixShift:suffixShift]
				paramValue = pair[1]
				trimmer := strings.Join([]string{delimiter, " "}, "")
				paramList = strings.Split(strings.Trim(pair[1], trimmer), delimiter)
			} else {
				paramName = pair[0][prefixShift:]
				paramValue = pair[1]
				paramList = []string{pair[1]}
			}
			tplParams[paramName] = envContainer{paramValue, paramList}
		}
	}

	return tplParams
}

func main() {
	var inFileName string
	var outFileName string
	var prefix string

	inFileName, outFileName, prefix = getArgs()

	tmpl, err := template.ParseFiles(inFileName)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(outFileName)
	if err != nil {
		panic(err)
	}

	tplParams := getParams(prefix, paramsListSuffix, paramsListDelimiter)

	tmpl.Execute(out, tplParams)
}
