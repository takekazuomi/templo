package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/urfave/cli/v2"
)

func importCsv(name string) []map[string]string {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
		panic(err)
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

func templateEvaulate(wr io.Writer, name string, data []map[string]string) {

	t := template.Must(template.ParseFiles(name))
	if err := t.Execute(wr, map[string][]map[string]string{"env": data}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	(&cli.App{
		Name:  "templo",
		Usage: "Expand the template",
		Action: func(c *cli.Context) error {
			r := importCsv("data/sample.csv")
			templateEvaulate(os.Stdout, "data/loop-sample.tmpl", r)
			return nil
		},
	}).Run(os.Args)
}
