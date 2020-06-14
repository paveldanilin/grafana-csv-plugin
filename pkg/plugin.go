package main

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
)

const (
	WelcomeMessage = "CSV plugin has been started"
	ExitMessage    = "CSV plugin has been stopped"
	Name           = "grafana_csv_plugin"
	Version        = "2.0.0"
)

func main() {
  	// Grafana logger
	var logger = hclog.New(&hclog.LoggerOptions{
		Name:  Name,
		Level: hclog.Trace,
	})

	// Welcome -> grafana log
	logger.Info(WelcomeMessage, "version", Version)

	// Start plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "grafana_plugin_type",
			MagicCookieValue: "datasource",
		},

		Plugins: map[string]plugin.Plugin{
		      Name: &datasource.DatasourcePluginImpl{Plugin: &CSVFileDatasource{
		      		MainLogger: logger,
		      		CsvDbManager: csv.NewDbManager(logger),
		      }},
		},

		// A non-nil value here enables gRPC serving for this pluginLogger...
		GRPCServer: plugin.DefaultGRPCServer,
	})

	logger.Info(ExitMessage, "version", Version)
}
