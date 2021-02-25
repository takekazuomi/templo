package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/urfave/cli/v2"
)

func importCsv(name string) []map[string]string {
	fmt.Println(name)

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

func templateEvaulate(wr io.Writer, name string, data []map[string]string) {

	t := template.Must(template.ParseFiles(name))
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
			templateEvaulate(os.Stdout, c.String("template"), r)
			return nil
		},
	})

	app.Run(os.Args)
}
