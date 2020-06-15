## CSV datasource for Grafana 6.x.x

#### Install
- Copy files from the `dist` folder to your [Grafana plugin directory](https://grafana.com/docs/grafana/latest/plugins/installation/#grafana-plugin-directory)
- Restart Grafana
- Check datasource list as http://you-grafana/datasources/new

> The plugin is unsigned, hence you may face the error:
>
> `lvl=eror msg=“Failed to load plugin” logger=plugins error=“plugin “grafana-csv-plugin” is unsigned”`
>
> To get it to work you should make configuration as described in [documentation](https://grafana.com/docs/grafana/latest/installation/configuration/#allow-loading-unsigned-plugins)

#### Features
- Read local CSV file
- Read CSV file via SFTP
- SQL queries (under the hood CSV will be converted into in-memory SQLite3 DB)


#### Query
- [SQLite3](https://www.sqlite.org/index.html)
- Each DS has its own table, the name of the table coincides with the DS name (for example the DS name is `my_data`, hence in a query you should select from `my_data` table)

#### Build graphs
- CSV File: data/SacramentocrimeJanuary2006.csv
- DS name `jan_2006`

```sql
SELECT strftime('%Y-%m-%d', cdatetime) as "time", district as "metric", count(*) as "value"
FROM jan_2006
GROUP BY strftime('%Y-%m-%d', cdatetime), district
ORDER BY cdatetime
```

![](./doc/image/graph.png)

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

#### Prev version
- [1.1.0](https://github.com/paveldanilin/grafana-csv-plugin/tree/1.1.0) which doesn't support SQL, but supports the filtering expressions.


###### Example data set: /data
###### Icon: https://freeicons.io/vector-file-types-icons/csv-icon-2272
