# go template with csv

```shell
templo -t data/fwrules.bicep.tmpl -c data/fwrules.csv

param mysql object // external resources


// https://docs.microsoft.com/ja-jp/azure/templates/microsoft.dbformysql/2017-12-01/servers/firewallrules

resource fw0_resource 'Microsoft.DBforMySQL/servers/firewallRules@2017-12-01' = {
  name: '${mysql.name}/apple'
  properties: {
    startIpAddress: '192.168.12.1/32'
    endIpAddress: '192.168.12.1/32'
  }
}

resource fw1_resource 'Microsoft.DBforMySQL/servers/firewallRules@2017-12-01' = {
  name: '${mysql.name}/orange'
  properties: {
    startIpAddress: '192.168.12.2/32'
    endIpAddress: '192.168.12.2/32'
  }
}
```

Run in docker

```shell
$ docker run --rm -i -v ${PWD}:/app takekazuomi/templo -t data/fwrules.bicep.tmpl -c data/fwrules.csv
```

## 0.1.0 new custom function

templo implemented three custom functions `split`, `barray`, `barrayq` to decompose strings
and convert them to bicep arrays. Separate the string with `split` and make it an array of
strings. `barray/barrayq` takes a string or an array of strings and converts it to bicep
array syntax.

- split string sep -> []string
  - split string to string array
- barray []string -> string
  - string array to bicep array
- barrayq []string -> string
  - string array to bicep array with quote

[go template Functions](https://pkg.go.dev/text/template#hdr-Functions)

### barray sample

A template is the go template syntax. The example below uses the pipeline syntax.

```go-template:data/barray.tmpl
{{- range $i, $v := .env -}}
bicep array: {{ .f1 | barray }}
bicep array with quote: {{ .f1 | barrayq }}
bicep array with sep: {{ barray .f2 "|" }}
{{- end }}
```

```csv:data/barray.csv
f1,f2
"1 2", "AB|CD"
```

```sh
$ docker run --rm -i -v ${PWD}:/app takekazuomi/templo -t data/barray.tmpl -c data/barray.csv
bicep array: [1, 2]
bicep array with quote: ['1', '2']
bicep array with sep: [AB, CD]
```

## ChangeLog

- 0.1.1
  - fix pipeline args bug
  - make more small container image


