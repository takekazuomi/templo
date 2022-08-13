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

// swapSplit properly swaps and splits arg1 and arg2
// When called from pipeline, sep comes in arg1. When called by function, sep is entered in arg2.
func swapSplit(arg1 reflect.Value, arg2 reflect.Value) (item []string, err error) {
	var s, sep = "", " "

	// check arg1
	switch arg1.Kind() {
	case reflect.Slice:
		item = reflect.Value.Interface(arg1).([]string)
	case reflect.String:
		s = arg1.String()
	case reflect.Int:
		sep = string(rune(int64(reflect.Value.Interface(arg1).(int))))
	default:
		return nil, fmt.Errorf("cannot split none string type %s", arg1.Type())
	}

	// check arg2
	switch arg2.Kind() {
	case reflect.Slice:
		item = reflect.Value.Interface(arg2).([]string)
	case reflect.String:
		if len(s) > 0 || len(item) > 0 {
			sep = arg2.String()
		} else {
			s = arg2.String()
		}
	case reflect.Int:
		sep = string(rune(int64(reflect.Value.Interface(arg2).(int))))
	case reflect.Invalid:
	default:
		return nil, fmt.Errorf("cannot split none string type %s", arg2.Type())
	}

	if len(s) > 0 {
		item = strings.Split(s, sep)
	}
	return
}

func barray(item reflect.Value, args ...reflect.Value) (string, error) {
	var arg2 reflect.Value
	if len(args) > 0 {
		arg2 = args[0]
	}
	sa, err := swapSplit(item, arg2)
	if err != nil {
		return "", err
	}

	return "[" + strings.Join(sa, ", ") + "]", nil
}

func barrayq(item reflect.Value, args ...reflect.Value) (string, error) {
	var arg2 reflect.Value
	if len(args) > 0 {
		arg2 = args[0]
	}
	sa, err := swapSplit(item, arg2)
	if err != nil {
		return "", err
	}

	for i, v := range sa {
		sa[i] = "'" + v + "'"
	}

	return "[" + strings.Join(sa, ", ") + "]", nil
}

func split(item reflect.Value, args ...reflect.Value) ([]string, error) {
	var arg2 reflect.Value
	if len(args) > 0 {
		arg2 = args[0]
	}
	return swapSplit(item, arg2)
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
