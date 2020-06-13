## CSV datasource for Grafana 6.x.x

#### Install
- Copy files from the `dist` folder to your [Grafana plugin directory](https://grafana.com/docs/grafana/latest/plugins/installation/#grafana-plugin-directory)
- Restart Grafana
- Check datasource list as http://you-grafana/datasources/new

> The plugin is unsigned, hence you may face the error:
>
> `lvl=eror msg=“Failed to load plugin” logger=plugins error=“plugin “grafana-csv-plugin” is unsigned”`
>
> To make it work you should make configuration as described in [documentation](https://grafana.com/docs/grafana/latest/installation/configuration/#allow-loading-unsigned-plugins)

#### Features
- Read local CSV file
- Read CSV file via SFTP
- SQL queries (under the hood CSV will be converted to in-memory SQLite3 DB)


#### Query
- [SQLite3](https://www.sqlite.org/index.html)
- Data from CSV are loaded into `DataTable` table per each Datasource

#### Build graphs
- data/SacramentocrimeJanuary2006.csv
```sql
SELECT strftime('%Y-%m-%d', cdatetime) as "time", district as "metric", count(*) as "value"
FROM DataTable
GROUP BY strftime('%Y-%m-%d', cdatetime), district
ORDER BY cdatetime
```

###### Simple table

![](./doc/image/grid.png)

#### Config
- Read local file

![](./doc/image/config_local.png)

- Read remote file

![](./doc/image/config_sftp.png)

#### Build
- npm run build

#### Test
- npm test


###### Example data set: /data
###### Icon: https://freeicons.io/vector-file-types-icons/csv-icon-2272
