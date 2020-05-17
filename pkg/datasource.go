package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/model"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/sftp"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"golang.org/x/net/context"
	"time"
)

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
	ds.debugf("DatasourceModel=%s", dsModel.String())

	queryModel, err := model.CreateQueryFrom(*req.Queries[0])
	if err != nil {
		errMsg := fmt.Sprintf("Could not create query: %s", err.Error())
		ds.logError(errMsg)
		ds.resultWithError(result, errMsg)
		return result, nil
	}
	ds.debugf("QueryModel=%s", queryModel.String())

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

func (ds *CSVFileDatasource) testConnection(dsModel *model.DatasourceModel) error {
	if dsModel.AccessMode == model.AccessMode_LOCAL {
		return ds.testConnectionLocal(dsModel)
	}
	if dsModel.AccessMode == model.AccessMode_SFTP {
		return ds.testConnectionSftp(dsModel)
	}
	return errors.New(fmt.Sprintf("unknown access mode `%s`", dsModel.AccessMode))
}

func (ds *CSVFileDatasource) testConnectionLocal(dsModel *model.DatasourceModel) error {
	return util.CheckFile(dsModel.Filename)
}

func (ds *CSVFileDatasource) testConnectionSftp(dsModel *model.DatasourceModel) error {
	return sftp.Test(sftp.ConnectionConfig{
		Host:          dsModel.SftpHost,
		Port:          dsModel.SftpPort,
		User:          dsModel.SftpUser,
		Password:      dsModel.SftpPassword,
		Timeout:       time.Second * 30, // TODO: Move to UI
		IgnoreHostKey: dsModel.SftpIgnoreHostKey,
	})
}

func (ds *CSVFileDatasource) performQuery(dsModel *model.DatasourceModel, queryModel *model.QueryModel) *datasource.QueryResult {
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

	return ds.queryLocalCsv(queryModel, &csv.FileDescriptor{
		Filename:         csvFilename,
		Delimiter:        rune(dsModel.CsvDelimiter[0]),
		Comment:          rune(dsModel.CsvComment[0]),
		TrimLeadingSpace: dsModel.CsvTrimLeadingSpace,
		FieldsPerRecord:  0, // Implies that each row contains the same count of fields as the header row
	})
}

func (ds *CSVFileDatasource) getRemoteFile(file string, dsModel *model.DatasourceModel) (string, error) {
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

func (ds *CSVFileDatasource) queryLocalCsv(queryModel *model.QueryModel, csvDesc *csv.FileDescriptor) *datasource.QueryResult {
	ds.debugf("Going to query data from `%s`", csvDesc.Filename)
	if err := util.CheckFile(csvDesc.Filename); err != nil {
		ds.logError(fmt.Sprintf("Could not read data file `%s`: %s", csvDesc.Filename, err.Error()))
		return ds.errQueryResult(queryModel.RefID, fmt.Sprintf("Could not read data file `%s`", csvDesc.Filename))
	}

	csvFile, err := csv.Read(csvDesc)
	if err != nil {
		ds.logError(fmt.Sprintf("Failed to read data file `%s`: %s", csvDesc.Filename, err.Error()))
		return ds.errQueryResult(queryModel.RefID, fmt.Sprintf("Failed to read data file `%s`", csvDesc.Filename))
	}
	ds.debugf("Table=%s", csvFile.String())


	if len(queryModel.Query) > 0 {
		ds.debugf("Apply query `%s`", queryModel.Query)
		csvFile, err = csvFile.Filter(queryModel.Query)
		if err != nil {
			ds.logError(fmt.Sprintf("Filtering is failed: %s. expr=`%s`", err.Error(), queryModel.Query))
			return ds.errQueryResult(queryModel.RefID, fmt.Sprintf("%s", err.Error()))
		}
	}

	ds.debugf("Rows count=%d", csvFile.RowsCount())

	table := ds.toTable(csvFile)

	return &datasource.QueryResult{
		RefId: queryModel.RefID,
		Tables: []*datasource.Table{table},
	}
}

func (ds *CSVFileDatasource) toTable(csvTable *csv.Table) *datasource.Table {
	table := &datasource.Table{
		Columns: []*datasource.TableColumn{},
		Rows:    make([]*datasource.TableRow, 0),
	}

	for _, columnName := range csvTable.Columns {
		table.Columns = append(table.Columns, &datasource.TableColumn{Name: columnName})
	}

	for _, row := range csvTable.Rows {
		tableRow := &datasource.TableRow{
			Values: make([]*datasource.RowValue, 0),
		}
		for _, value := range row {
			switch typedValue := value.(type) {
			case int64:
				tableRow.Values = append(tableRow.Values, &datasource.RowValue{
					Kind:		datasource.RowValue_TYPE_INT64,
					Int64Value:	typedValue,
				})
				break
			case float64:
				tableRow.Values = append(tableRow.Values, &datasource.RowValue{
					Kind:		datasource.RowValue_TYPE_DOUBLE,
					DoubleValue:	typedValue,
				})
				break
			case string:
				tableRow.Values = append(tableRow.Values, &datasource.RowValue{
					Kind:        datasource.RowValue_TYPE_STRING,
					StringValue: typedValue,
				})
				break
			case nil:
				tableRow.Values = append(tableRow.Values, &datasource.RowValue{
					Kind: datasource.RowValue_TYPE_NULL,
				})
				break
			}
		}
		table.Rows = append(table.Rows, tableRow)
	}

	return table
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
