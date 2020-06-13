package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/grafana"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/model"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/sftp"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"golang.org/x/net/context"
	"time"
)

var csvDbManager = csv.NewDbManager()

type CSVFileDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	MainLogger hclog.Logger
}

func (ds *CSVFileDatasource) Query(ctx context.Context, req *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	ds.logRequest(req)

	result := &datasource.DatasourceResponse{
		Results: make([]*datasource.QueryResult, 0),
	}

	if len(req.Queries) == 0 {
		ds.logWarning("No queries, nothing to execute")
		return nil, errors.New("no queries, nothing to execute")
	}

	dsModel, err := model.CreateDatasourceFrom(*req)
	if err != nil {
		errMsg := fmt.Sprintf("Could not create datasource: %s", err.Error())
		ds.logError(errMsg)
		ds.resultWithError(result, errMsg)
		return result, nil
	}
	ds.debugf("Datasource=%s", dsModel.String())

	queryModel, err := model.CreateQueryFrom(*req.Queries[0])
	if err != nil {
		errMsg := fmt.Sprintf("Could not create query: %s", err.Error())
		ds.logError(errMsg)
		ds.resultWithError(result, errMsg)
		return result, nil
	}
	ds.debugf("Query=%s", queryModel.String())

	// RefId is hardcoded in datasource.js
	if queryModel.RefID == "[tests-connection]" {
		err := ds.testConnection(dsModel)
		if err != nil {
			ds.resultWithError(result, err.Error())
			return result, nil
		}
		return result, nil
	}

	result.Results = append(result.Results, ds.performQuery(dsModel, queryModel))

	return result, nil
}

func (ds *CSVFileDatasource) testConnection(dsModel *model.Datasource) error {
	if dsModel.AccessMode == model.AccessMode_LOCAL {
		return ds.testConnectionLocal(dsModel)
	}
	if dsModel.AccessMode == model.AccessMode_SFTP {
		return ds.testConnectionSftp(dsModel)
	}
	return errors.New(fmt.Sprintf("unknown access mode `%s`", dsModel.AccessMode))
}

func (ds *CSVFileDatasource) testConnectionLocal(dsModel *model.Datasource) error {
	return util.CheckFile(dsModel.Filename)
}

func (ds *CSVFileDatasource) testConnectionSftp(dsModel *model.Datasource) error {
	return sftp.Test(sftp.ConnectionConfig{
		Host:          dsModel.SftpHost,
		Port:          dsModel.SftpPort,
		User:          dsModel.SftpUser,
		Password:      dsModel.SftpPassword,
		Timeout:       time.Second * 30, // TODO: Move to UI
		IgnoreHostKey: dsModel.SftpIgnoreHostKey,
	})
}

func (ds *CSVFileDatasource) performQuery(dsModel *model.Datasource, queryModel *model.Query) *datasource.QueryResult {
	csvFilename := dsModel.Filename

	if dsModel.AccessMode == model.AccessMode_SFTP {
		ds.debugf("Going to download remote file `%s`", csvFilename)
		downloadedFile, err := ds.getRemoteFile(csvFilename, dsModel)
		if err != nil {
			return &datasource.QueryResult{
				Error: fmt.Sprintf("Could not download CSV data file: %s", err.Error()),
				RefId: queryModel.RefID,
			}
		}
		ds.debugf("File has been downloaded %s->%s", csvFilename, downloadedFile)
		csvFilename = downloadedFile
	}

	tableColumns := make([]csv.Column, 0)
	for _, dsColumn := range dsModel.Columns {
		tableColumns = append(tableColumns, csv.Column{
			Type: csv.ColumnTypeFromString(dsColumn.Type),
			Name: dsColumn.Name,
		})
	}

	csvDb, err := csvDbManager.Get(dsModel.Name, &csv.FileDescriptor{
		Filename:         csvFilename,
		Delimiter:        rune(dsModel.CsvDelimiter[0]),
		Comment:          rune(dsModel.CsvComment[0]),
		TrimLeadingSpace: dsModel.CsvTrimLeadingSpace,
		FieldsPerRecord:  0, // Implies that each row contains the same count of fields as the header row
		Columns: tableColumns,
	})
	if err != nil {
		return &datasource.QueryResult{
			Error: fmt.Sprintf("Could not parse CSV: %s", err.Error()),
			RefId: queryModel.RefID,
		}
	}

	result, err := csvDb.Query(queryModel.Query)
	if err != nil {
		return &datasource.QueryResult{
			Error: fmt.Sprintf("Query failed: %s", err.Error()),
			RefId: queryModel.RefID,
		}
	}
	defer result.Release()

	return ds.toGrafanaResult(result, queryModel)
}

func (ds *CSVFileDatasource) toGrafanaResult(result *csv.QueryResult, queryModel *model.Query) *datasource.QueryResult {
	if queryModel.Format == "time_series" {
		return ds.toGrafanaTimeseries(queryModel.RefID, result)
	}
	return ds.toGrafanaTable(queryModel.RefID, result)
}

func (ds *CSVFileDatasource) toGrafanaTimeseries(refId string, result *csv.QueryResult) *datasource.QueryResult {
	timeseries, err := grafana.ToTimeSeries(result)
	if err != nil {
		return &datasource.QueryResult{
			Error: fmt.Sprintf("Query failed: %s", err.Error()),
			RefId: refId,
		}
	}
	queryResult := &datasource.QueryResult{
		RefId: refId,
		Series: []*datasource.TimeSeries{},
	}
	for seriesName, seriesPoints := range timeseries {
		queryResult.Series = append(queryResult.Series, &datasource.TimeSeries{
			Name:   seriesName,
			Points: seriesPoints,
		})
	}
	return queryResult
}

func (ds *CSVFileDatasource) toGrafanaTable(refId string, result *csv.QueryResult) *datasource.QueryResult {
	table, err := grafana.ToTable(result)
	if err != nil {
		return &datasource.QueryResult{
			Error: fmt.Sprintf("Query failed: %s", err.Error()),
			RefId: refId,
		}
	}
	return &datasource.QueryResult{
		RefId: refId,
		Tables: []*datasource.Table{table},
	}
}

func (ds *CSVFileDatasource) getRemoteFile(file string, dsModel *model.Datasource) (string, error) {
	downloadedFile, err := sftp.GetFile(sftp.ConnectionConfig{
		Host:          dsModel.SftpHost,
		Port:          dsModel.SftpPort,
		User:          dsModel.SftpUser,
		Password:      dsModel.SftpPassword,
		Timeout:       time.Second * 30,
		IgnoreHostKey: dsModel.SftpIgnoreHostKey,
	}, file, dsModel.SftpWorkingDir)
	if err != nil {
		ds.logError(fmt.Sprintf("Could not download CSV data file: %s", err.Error()))
		return "", err
	}
	return downloadedFile, nil
}

func (ds *CSVFileDatasource) toTimeSeries(result *csv.QueryResult) *datasource.TimeSeries {
	return nil
}


func (ds *CSVFileDatasource) resultWithError(result *datasource.DatasourceResponse, errorMessage string) {
	result.Results = make([]*datasource.QueryResult, 0)
	result.Results = append(result.Results, &datasource.QueryResult{
		RefId: "A",
		Error: errorMessage,
	})
}

func (ds *CSVFileDatasource) errQueryResult(refId string, message string) *datasource.QueryResult {
	return &datasource.QueryResult{
		Error: message,
		RefId: refId,
	}
}

func (ds *CSVFileDatasource) logRequest(req *datasource.DatasourceRequest) {
	if ds.MainLogger.IsDebug() == false {
		return
	}
	logContext := make(map[string]interface{}, 0)
	logContext["method"] = "logRequest"
	reqJson, _ := json.Marshal(req)
	logContext["attributes"] = string(reqJson)
	ds.logDebug("Request", logContext)
}

func (ds *CSVFileDatasource) logDebug(msg string, context map[string]interface{}) {
	context["version"] = Version
	ds.MainLogger.Debug(msg, util.MapToArray(context))
}

func (ds *CSVFileDatasource) logInfo(msg string) {
	logContext := map[string]interface{}{}
	logContext["version"] = Version
	ds.MainLogger.Info(msg, util.MapToArray(logContext))
}

func (ds *CSVFileDatasource) logWarning(msg string) {
	logContext := map[string]interface{}{}
	logContext["version"] = Version
	ds.MainLogger.Warn(msg, util.MapToArray(logContext))
}

func (ds *CSVFileDatasource) logError(msg string) {
	logContext := map[string]interface{}{}
	logContext["version"] = Version
	ds.MainLogger.Error(msg, util.MapToArray(logContext))
}

func (ds *CSVFileDatasource) debugf(msg string, args ...interface{}) {
	ds.MainLogger.Debug(fmt.Sprintf(msg, args...))
}
