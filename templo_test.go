package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func dummyMap(data map[string]string) {
	fmt.Println(data)
}

func dummy(header []string, records []string) {
	fmt.Println(header)
	fmt.Println(records)
}

func TestCSVSuccess(t *testing.T) {
	r := importCsv("data/sample.csv")
	for i, v := range r {
		if i == 0 && v["f1"] != "1" {
			t.Errorf("line: %v, result f1: %v", i, v["f1"])
		}
		if i == 1 && v["f1"] != "2" {
			t.Errorf("line: %v, result f1: %v", i, v["f1"])
		}
	}
}

func TestTemplateEvaulateSuccess(t *testing.T) {
	r := importCsv("data/sample.csv")
	var out bytes.Buffer
	templateEvaulate(&out, "data/loop-sample.tmpl", r)
	s := out.String()
	if !strings.Contains(s, "0: f1 1, f2: apple, f3: orange, f4: 14") {
		t.Errorf("erroe unexpeded result %v", s)
	}
}

func TestTemplateEvaulateFwRuleSuccess(t *testing.T) {
	r := importCsv("data/fwrules.csv")
	var out bytes.Buffer
	templateEvaulate(&out, "data/fwrules.bicep.tmpl", r)
	s := out.String()

	file, err := os.Create("data/fwfutes.bicep")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(s)
	file.Close()

	if !strings.Contains(s, "0: f1 1, f2: apple, f3: orange, f4: 14") {
		t.Errorf("erro unexpeded result %v", s)
	}

}
