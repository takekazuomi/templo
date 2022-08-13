package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
)

func importCsv(name string) []map[string]string {

	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comment = '#'
	r.TrimLeadingSpace = true

	var header []string
	var results []map[string]string

	for {
		records, err := r.Read()
		if err != nil {
			if err != io.EOF {
				log.Fatalln(err, records)
			}
			break
		}
		if header == nil {
			header = records
		} else {
			m := map[string]string{}
			for i, s := range records {
				m[header[i]] = s
			}
			results = append(results, m)
		}
	}
	return results
}

func reflectSplit(item reflect.Value, sep string) ([]string, error) {
	var sa []string
	switch item.Kind() {
	case reflect.Slice:
		sa = reflect.Value.Interface(item).([]string)
	case reflect.String:
		sa = strings.Split(item.String(), sep)
	default:
		return nil, fmt.Errorf("cannot split none string type %s", item.Type())
	}
	return sa, nil
}

func barray(item reflect.Value, args ...string) (string, error) {
	sep := " "
	if len(args) > 0 {
		sep = args[0]
	}
	sa, err := reflectSplit(item, sep)
	if err != nil {
		return "", err
	}
	return "[" + strings.Join(sa, ", ") + "]", nil
}

func barrayq(item reflect.Value, args ...string) (string, error) {
	sep := " "
	if len(args) > 0 {
		sep = args[0]
	}
	sa, err := reflectSplit(item, sep)
	if err != nil {
		return "", err
	}

	for i, v := range sa {
		sa[i] = "'" + v + "'"
	}

	return "[" + strings.Join(sa, ", ") + "]", nil
}

func split(s string, options ...string) ([]string, error) {
	sep := " "
	if len(options) > 0 {
		sep = options[0]
	}
	return strings.Split(s, sep), nil
}

func templateEvaluate(wr io.Writer, name string, data []map[string]string) {

	// https://pkg.go.dev/text/template#Template.Funcs
	// > t should usually have the name of one of the (base) names of the files.
	// テンプレートの名前は、ParseFilesで指定されたファイルのbase name の一つで無ければいけない。
	t, err := template.New(path.Base(name)).Funcs(template.FuncMap{
		"split":   split,
		"barray":  barray,
		"barrayq": barrayq,
	}).ParseFiles(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(wr, map[string][]map[string]string{"env": data}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := (&cli.App{
		Name:  "templo",
		Usage: "Expand the template",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "template",
				Aliases:  []string{"t"},
				Usage:    "Load template from `FILE`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "csv",
				Aliases:  []string{"c"},
				Usage:    "Load csv data from `FILE`",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			r := importCsv(c.String("csv"))
			templateEvaluate(os.Stdout, c.String("template"), r)
			return nil
		},
	})

	app.Run(os.Args)
}
