package grafana

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"io"
	"time"
)

func ToTimeSeries(result *csv.QueryResult) (map[string][]*datasource.Point, error) {
	series := make(map[string][]*datasource.Point)

	// -- Columns
	columnNames, err := result.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := result.ColumnTypes()
	if err != nil {
		return nil, err
	}

	timeColIndex := -1
	metricColIndex := -1
	valueColIndex := -1

	// Detect column index
	for i, colName := range columnNames {
		switch colName {
		case "metric":
			metricColIndex = i
		case "time":
			timeColIndex = i
		case "value":
			valueColIndex = i
		}
	}

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

func ToTable(result *csv.QueryResult) (*datasource.Table, error) {
	table := &datasource.Table{
		Columns: []*datasource.TableColumn{},
		Rows: make([]*datasource.TableRow, 0),
	}

	// -- Columns
	columns, err := result.Columns()
	if err != nil {
		return nil, err
	}

	for _, columnName := range columns {
		table.Columns = append(table.Columns, &datasource.TableColumn{Name: columnName})
	}

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

