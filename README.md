# go template with csv

```shell
templo -t data/fwrules.bicep.tmpl -c data/fwrules.csv

param mysql object // external resouces


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
docker run --rm -i takekazuomi/templo -t data/fwrules.bicep.tmpl -c data/fwrules.csv
```

## TODO

`build/*` is root permission, it's hard to live.



