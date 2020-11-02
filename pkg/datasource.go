package main

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-plugin"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/service"
	"golang.org/x/net/context"
)

type CsvDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	pluginService service.IPluginService
}

func NewCsvDatasourcePlugin(pluginService service.IPluginService) datasource.DatasourcePlugin {
	return &CsvDatasource{
		pluginService: pluginService,
	}
}

func (ds *CsvDatasource) Query(_ context.Context, req *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	return ds.pluginService.Query(req)
}
