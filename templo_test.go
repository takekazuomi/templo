package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVSuccess(t *testing.T) {
	r := importCsv("data/sample.csv")
	for i, v := range r {
		assert.False(t, i == 0 && v["f1"] != "1", "line: %v, result f1: %v", i, v["f1"])
		assert.False(t, i == 1 && v["f1"] != "2", "line: %v, result f1: %v", i, v["f1"])
	}
}

func TestTemplateEvaluateSuccess(t *testing.T) {
	r := importCsv("data/sample.csv")
	var out bytes.Buffer
	templateEvaluate(&out, "data/loop-sample.tmpl", r)
	s := out.String()

	assert.Contains(t, s, "0: f1 1, f2: apple, f3: orange, f4: 14", "error unexpected result %v", s)
}

func TestTemplateEvaluateFwRuleSuccess(t *testing.T) {
	r := importCsv("data/fwrules.csv")
	var out bytes.Buffer
	templateEvaluate(&out, "data/fwrules.bicep.tmpl", r)
	s := out.String()

	file, err := os.Create("data/fwrules.bicep")
	assert.NoError(t, err)
	defer file.Close()

	file.WriteString(s)
	file.Close()

	assert.Contains(t, s, "/apple", "error unexpected result %v", s)
	assert.Contains(t, s, "/orange", "error unexpected result %v", s)
	assert.Contains(t, s, "192.168.12.2/32", "error unexpected result %v", s)
}

func TestTemplateSplit(t *testing.T) {
	tests := []struct {
		name  string
		templ string
		csv   string
		out   string
		want  []string
	}{
		{
			name:  "split-pipe",
			templ: "data/split-pipe.tmpl",
			csv:   "data/split-pipe.csv",
			out:   "data/split-pipe.out",
			want: []string{
				"[1 2]",
			},
		},
		{
			name:  "barray",
			templ: "data/barray.tmpl",
			csv:   "data/barray.csv",
			out:   "data/barray.out",
			want: []string{
				"bicep array: [1, 2]",
				"bicep array with quote: ['1', '2']",
				"bicep array with sep: [AB, CD]",
			},
		},
		{
			name:  "split",
			templ: "data/split.tmpl",
			csv:   "data/split.csv",
			out:   "data/split.out",
			want: []string{
				"[1 2] is array. len 2",
				"bicep array: [1, 2]",
				"bicep array with quote: ['1', '2']",
				"bicep array func: ['1', '2']",
			},
		},
		{
			name:  "barray pipe",
			templ: "data/barray-pipe.tmpl",
			csv:   "data/barray.csv",
			out:   "data/barray-pipe.out",
			want: []string{
				"bicep array with sep: [AB, CD]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := importCsv(tt.csv)
			var out bytes.Buffer
			templateEvaluate(&out, tt.templ, r)
			s := out.String()

			file, err := os.Create(tt.out)
			if err != nil {
				t.Error(err)
			}
			defer file.Close()

			file.WriteString(s)
			file.Close()

			for _, e := range tt.want {
				assert.Contains(t, s, e)
			}
		})
	}
}
