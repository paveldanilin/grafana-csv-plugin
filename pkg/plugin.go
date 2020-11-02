package main

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/time_filter"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/time_group"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/unix_epoch_from"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/unix_epoch_to"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/service"
)

const (
	WelcomeMessage = "CSV plugin has been started"
	ExitMessage    = "CSV plugin has been stopped"
	Name           = "grafana_csv_plugin"
	Version        = "3.0.0"
)

func main() {
	loggerService := service.NewLoggerService(hclog.New(&hclog.LoggerOptions{
		Name:  Name,
		Level: hclog.Trace,
	}), map[string]interface{}{
		"plugin": Name,
		"version": Version,
	})
	loggerService.Info(WelcomeMessage)

	csvDb, err := service.NewSqliteDb(100, 0, loggerService)
	if err != nil {
		loggerService.Error("Could not create CSV database", "error", err.Error())
		return
	}
	if err := csvDb.Init(); err != nil {
		loggerService.Error("Could not init CSV database", "error", err.Error())
		return
	}

	csvService := service.NewCsvService(csvDb, loggerService)
	grafanaService := service.NewGrafanaService(loggerService, csvService)
	pluginService := service.NewPluginService(service.PluginConfig{
		Logger:         loggerService,
		CsvService:     csvService,
		GrafanaService: grafanaService,
	})

	macro.Register(time_filter.MacroName, time_filter.Processor)
	macro.Register(unix_epoch_from.MacroName, unix_epoch_from.Processor)
	macro.Register(unix_epoch_to.MacroName, unix_epoch_to.Processor)
	macro.Register(time_group.MacroName, time_group.Processor)

	// Start plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "grafana_plugin_type",
			MagicCookieValue: "datasource",
		},

		Plugins: map[string]plugin.Plugin{
		      Name: &datasource.DatasourcePluginImpl{Plugin: NewCsvDatasourcePlugin(pluginService)},
		},

		// A non-nil value here enables gRPC serving for this pluginLogger...
		GRPCServer: plugin.DefaultGRPCServer,
	})

	loggerService.Info(ExitMessage)
}
