## Grafana CSV datasource

Grafana backend plugin which allows you:
- Read local CSV file
- Download remote CSV file by SFTP

![](./doc/image/grid.png)

## Config
- Local

![](./doc/image/config_local.png)

- Remote

![](./doc/image/config_sftp.png)

## Build
- npm run build

## Test
- go test ./tests -v

## Install
- Copy `dist` to `/etc/grafana/data/plugins/grafana-csv-plugin/`

###### Icon: https://freeicons.io/vector-file-types-icons/csv-icon-2272
