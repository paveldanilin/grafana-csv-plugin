package main

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/time_filter"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/time_group"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/unix_epoch_from"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro/unix_epoch_to"
)

const (
	WelcomeMessage = "CSV plugin has been started"
	ExitMessage    = "CSV plugin has been stopped"
	Name           = "grafana_csv_plugin"
	Version        = "2.0.0"
)

func main() {
  	// GF logger
	var logger = hclog.New(&hclog.LoggerOptions{
		Name:  Name,
		Level: hclog.Trace,
	})
	logger.Info(WelcomeMessage, "version", Version)

	csvDb, err := csv.NewDB(100, 0, logger)
	if err != nil {
		logger.Error("Could not create CSV database", "error", err.Error())
		return
	}
	if err := csvDb.Init(); err != nil {
		logger.Error("Could not init CSV database", "error", err.Error())
		return
	}

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
		      Name: &datasource.DatasourcePluginImpl{Plugin: &CSVFileDatasource{
		      		MainLogger: logger,
		      		Db: csvDb,
		      }},
		},

		// A non-nil value here enables gRPC serving for this pluginLogger...
		GRPCServer: plugin.DefaultGRPCServer,
	})

	logger.Info(ExitMessage, "version", Version)
}
