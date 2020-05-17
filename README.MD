## Grafana CSV datasource plugin

Grafana datasource which allows:
- Read local CSV file
- Read remote CSV file by SFTP
- Filter data (thanks to a great [expr library](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md))

![](./doc/image/grid.png)

![](./doc/image/expr.png)

![](./doc/image/expr1.png)

## Config
- Local

![](./doc/image/config_local.png)

- Remote

![](./doc/image/config_sftp.png)

## Build
- npm run build

## Test
- npm test

## Install
- Copy `dist` to `/etc/grafana/data/plugins/grafana-csv-plugin/`

###### Icon: https://freeicons.io/vector-file-types-icons/csv-icon-2272
