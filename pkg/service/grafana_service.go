package service

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/xsql"
	"io"
	"strings"
	"time"
)

const (
	GrafanaFormatTimeSeries = "time_series"
	GrafanaFormatTable = "table"
	TsTimeColumnMap = "autotime,time"
	TsMetricColumnMap = "metric"
	TsValueColumnMap = "value"
)

type IGrafanaService interface {
	ErrorResponse(refID string, result *datasource.DatasourceResponse, errorMessage string)
	ToQueryResult(refID string, format string, result *xsql.QueryResult) *datasource.QueryResult
	ToTimeSeries(result *xsql.QueryResult, colMap TsColumnsMap) (map[string][]*datasource.Point, error)
	ToTable(result *xsql.QueryResult) (*datasource.Table, error)
}

// Comma separated lists of column names
type TsColumnsMap struct {
	Metric string
	Value string
	Time string
}

type grafanaService struct {
	logger     ILoggerService
	csvService ICsvService
}

func NewGrafanaService(logger ILoggerService, csvService ICsvService) IGrafanaService {
	return &grafanaService{logger: logger, csvService: csvService}
}

func (s *grafanaService) ErrorResponse(refID string, result *datasource.DatasourceResponse, errorMessage string) {
	result.Results = make([]*datasource.QueryResult, 0)
	result.Results = append(result.Results, s.errQueryResult(refID, errorMessage))
}

func (s *grafanaService) ToQueryResult(refID string, format string, result *xsql.QueryResult) *datasource.QueryResult {
	switch strings.ToLower(format) {
	case GrafanaFormatTimeSeries:
		timeSeries, err := s.ToTimeSeries(result, TsColumnsMap{
			Metric: TsMetricColumnMap,
			Value:  TsValueColumnMap,
			Time:   TsTimeColumnMap,
		})
		if err != nil {
			return s.errQueryResult(refID, fmt.Sprintf("Query failed: %s", err.Error()))
		}
		queryResult := &datasource.QueryResult{
			RefId: refID,
			Series: []*datasource.TimeSeries{},
		}
		for seriesName, seriesPoints := range timeSeries {
			queryResult.Series = append(queryResult.Series, &datasource.TimeSeries{
				Name:   seriesName,
				Points: seriesPoints,
			})
		}
		return queryResult
	case GrafanaFormatTable:
		table, err := s.ToTable(result)
		if err != nil {
			return s.errQueryResult(refID, fmt.Sprintf("Query failed: %s", err.Error()))
		}
		return &datasource.QueryResult{
			RefId: refID,
			Tables: []*datasource.Table{table},
		}
	}

	return s.errQueryResult(refID, fmt.Sprintf("Query failed: unknown format [%s]", format))
}

func (s *grafanaService) ToTimeSeries(result *xsql.QueryResult, colMap TsColumnsMap) (map[string][]*datasource.Point, error) {
	series := make(map[string][]*datasource.Point)

	// -- Columns
	columnNames := result.Columns()
	columnTypes := result.ColumnTypes()

	timeColIndex := -1
	metricColIndex := -1
	valueColIndex := -1

	// Detect column index
	for i, colName := range columnNames {
		if timeColIndex == -1 && util.InArray(strings.Split(colMap.Time, ","), colName) {
			timeColIndex = i
		} else if metricColIndex == -1 && util.InArray(strings.Split(colMap.Metric, ","), colName) {
			metricColIndex = i
		} else if valueColIndex == -1 && util.InArray(strings.Split(colMap.Value, ","), colName) {
			valueColIndex = i
		}
	}

	// Look for the first column with DB type DATE
	if timeColIndex == -1 {
		// Searches for first column with DATE oracle type and NOT nullable
		for i := range columnNames {
			nullable, _ := columnTypes[i].Nullable()
			if columnTypes[i].DatabaseTypeName() == "DATE" && nullable == false  {
				timeColIndex = i
				break
			}
		}
	}

	s.logger.Debug("Convert to time series",
		"metricIndex", metricColIndex,
		"timeIndex", timeColIndex,
		"valueIndex", valueColIndex,
		"columns", columnNames,
	)

	for {
		row, err := result.Next()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		for seriesName, point := range toGraphPoints(row, columnNames, timeColIndex, metricColIndex, valueColIndex) {
			_, exists := series[seriesName]
			if exists == false {
				series[seriesName] = make([]*datasource.Point, 0)
			}

			series[seriesName] = append(series[seriesName], point)
		}
	}

	return series, nil
}

func (s *grafanaService) ToTable(result *xsql.QueryResult) (*datasource.Table, error) {
	table := &datasource.Table{
		Columns: []*datasource.TableColumn{},
		Rows: make([]*datasource.TableRow, 0),
	}

	// -- Columns
	columns := result.Columns()

	for _, columnName := range columns {
		table.Columns = append(table.Columns, &datasource.TableColumn{Name: columnName})
	}

	s.logger.Debug("Convert to a table",
		"columns", columns,
	)

	// -- Rows
	for {
		row, err := result.Next()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		table.Rows = append(table.Rows, toTableRow(row))
	}

	return table, nil
}

func (s *grafanaService) errQueryResult(refID string, message string) *datasource.QueryResult {
	return &datasource.QueryResult{
		Error: message,
		RefId: refID,
	}
}

// Converts row values to graph points => map[SeriesName] = Point{}
func toGraphPoints(values []interface{}, columnNames []string, timeColIndex int, metricColIndex int, valueColIndex int) map[string]*datasource.Point {
	points := make(map[string]*datasource.Point)
	metricTime := time.Now().UnixNano() / 1000000

	if metricColIndex == -1 {
		// Each column except `time` represents a time series
		for i := range values {
			if i == timeColIndex {
				// Skip time column
				continue
			}

			if timeColIndex != -1 {
				switch tv := values[timeColIndex].(type) {
				case time.Time:
					metricTime = util.TimeToEpochMs(tv)
					break
				case string:
					ttv, _ := dateparse.ParseAny(tv)
					metricTime = util.TimeToEpochMs(ttv)
					break
				}

			}

			seriesName := columnNames[i]
			metricValue, err := util.ToFloat64(values[i])
			if err != nil {
				metricValue = 0
			}

			points[seriesName] = &datasource.Point{
				Timestamp:	metricTime,
				Value:		metricValue,
			}
		}
	} else {
		if timeColIndex != -1 {
			switch tv := values[timeColIndex].(type) {
			case time.Time:
				metricTime = util.TimeToEpochMs(tv)
				break
			case string:
				ttv, _ := dateparse.ParseAny(tv)
				metricTime = util.TimeToEpochMs(ttv)
				break
			}
		}

		seriesName := fmt.Sprint(values[metricColIndex])
		metricValue, err := util.ToFloat64(values[valueColIndex])
		if err != nil {
			metricValue = 0
		}

		points[seriesName] = &datasource.Point{
			Timestamp:	metricTime,
			Value:		metricValue,
		}
	}

	return points
}


func toTableRow(row []interface{}) *datasource.TableRow {
	values := make([]*datasource.RowValue, 0)
	for _, value := range row {
		rowValue := toRowValue(value)
		values = append(values, rowValue)
	}
	return &datasource.TableRow{
		Values: values,
	}
}

func toRowValue(value interface{}) *datasource.RowValue {
	switch typedValue := value.(type) {
	case time.Time:
		return &datasource.RowValue{
			Kind:		datasource.RowValue_TYPE_INT64,
			Int64Value:	typedValue.UnixNano() / int64(time.Millisecond),
		}
	case int64:
		return &datasource.RowValue{
			Kind:		datasource.RowValue_TYPE_INT64,
			Int64Value:	typedValue,
		}
	case float64:
		return &datasource.RowValue{
			Kind:		datasource.RowValue_TYPE_DOUBLE,
			DoubleValue:	typedValue,
		}
	case string:
		return &datasource.RowValue{
			Kind:        datasource.RowValue_TYPE_STRING,
			StringValue: typedValue,
		}
	case nil:
		return &datasource.RowValue{
			Kind: datasource.RowValue_TYPE_NULL,
		}
	}
	return nil
}
