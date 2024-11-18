package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/dogenzaka/tsv"
	"github.com/iancoleman/strcase"
)

type ExitCode struct {
	Code        int
	GoStyleName string
	PyStyleName string
	Description string
	Group       string
}

type ExitCodeTSVRow struct {
	Code        int
	Name        string
	PSOverride  string
	GSOverride  string
	Description string
	Group       string
}

//go:embed exitcodes.tsv
var ecSrc string

//go:embed golang.tmpl
var goTmpl string

//go:embed python.tmpl
var pyTmpl string

func main() {
	goOutfile := flag.String("go", "", "Golang output file")
	pyOutfile := flag.String("py", "", "Python output file")
	flag.Parse()

	if *goOutfile == "" && *pyOutfile == "" {
		log.Fatalln("At least one of -go or -py must be given.")
		return
	}

	var ecs []ExitCode
	var ecRow ExitCodeTSVRow
	rowN := 0
	parser := tsv.NewParserWithoutHeader(strings.NewReader(ecSrc), &ecRow)
	for {
		eof, err := parser.Next()
		rowN++
		if eof {
			break
		}
		if err != nil {
			log.Fatalf("row %d: failed to parse: %s", rowN, err)
		}

		ec := ExitCode{
			Code:        ecRow.Code,
			Description: ecRow.Description,
			Group:       ecRow.Group,
		}

		name := ecRow.Name
		if strings.HasPrefix(name, "EX_") {
			name = strings.TrimPrefix(name, "EX_")
		} else if strings.HasPrefix(name, "EXIT_") {
			name = strings.TrimPrefix(name, "EXIT_")
		}

		if ecRow.PSOverride != "" {
			ec.PyStyleName = ecRow.PSOverride
		} else {
			ec.PyStyleName = name
		}

		if ecRow.GSOverride != "" {
			ec.GoStyleName = ecRow.GSOverride
		} else {
			ec.GoStyleName = strcase.ToCamel(name)
		}

		ecs = append(ecs, ec)
	}

	if *goOutfile != "" {
		if err := execTmpl("go", *goOutfile, goTmpl, ecs); err != nil {
			log.Fatalf("Failed to execute Go template: %v", err)
		}
	}
	if *pyOutfile != "" {
		if err := execTmpl("py", *pyOutfile, pyTmpl, ecs); err != nil {
			log.Fatalf("Failed to execute Python template: %v", err)
		}
	}
}

func execTmpl(name, outfile, templ string, ecs []ExitCode) error {
	t, err := template.New(name).Parse(templ)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	if err := t.Execute(f, ecs); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
