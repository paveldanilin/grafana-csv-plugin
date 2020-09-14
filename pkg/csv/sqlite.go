package csv

import (
	"database/sql"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/hashicorp/go-hclog"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/model"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"io"
	"strconv"
	"strings"
	"time"
	"sync"
)

type DbSqlite struct {
	db *sql.DB
	logger hclog.Logger
	mux sync.Mutex
}

const metaCsvTable = "_meta_csv_"

// If maxIdleCons <= 0, no idle connections are retained
// If connMaxLifetime <= 0, connections are reused forever.
func NewDB(maxIdleCons int, connMaxLifetime time.Duration, logger hclog.Logger) (DB, error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleCons)
	db.SetConnMaxLifetime(connMaxLifetime)

	return &DbSqlite{db: db, logger: logger}, nil
}

func (sqlite *DbSqlite) Init() error {
	sqlite.logger.Debug("Init CSV DB")
	return sqlite.createMetaCsvTable()
}

func (sqlite *DbSqlite) Query(sql string) (*QueryResult, error) {
	sqlite.logger.Debug("Query", "sql", sql)
	rows, err := sqlite.db.Query(sql)
	if err != nil {
		sqlite.logger.Error("Query failed", "error", err.Error())
		return nil, err
	}
	return newQueryResult(rows)
}

func (sqlite *DbSqlite) LoadCSV(tableName string, descriptor *FileDescriptor) error {
	sqlite.mux.Lock()
	defer sqlite.mux.Unlock()

	sqlite.logger.Info("Loading CSV", "table", tableName, "filename", descriptor.Filename)

	var metaCsv *model.Meta
	reload := false
	tableExists, err := sqlite.ifTableExists(tableName)
	if err != nil {
		return err
	}

	if tableExists {
		metaCsv = sqlite.getMetaCsv(tableName)
		if metaCsv == nil {
			sqlite.logger.Debug("CSV already loaded", "table", tableName, "filename", descriptor.Filename, "meta", "nil", "reload", false)
			return nil
		}

		fSize, fModTime := util.FileStat(descriptor.Filename)
		if fSize == metaCsv.FileSize && fModTime == metaCsv.FileModTime {
			// the file is not changed
			sqlite.logger.Debug("CSV already loaded", "table", tableName, "filename", descriptor.Filename, "changed", false, "reload", false)
			return nil
		}

		// The file is changed, we should reload it
		sqlite.logger.Debug("CSV already loaded", "table", tableName, "filename", descriptor.Filename, "changed", true, "reload", true)
		reload = true
		metaCsv.FileSize = fSize
		metaCsv.FileModTime = fModTime
	}

	reader, err := newCsvReader(descriptor)
	if err != nil {
		sqlite.logger.Debug("Failed to create CSV reader", "error", err.Error(), "filename", descriptor.Filename)
		return err
	}
	defer reader.close()

	if reload {
		_ = sqlite.updateMetaCsv(metaCsv)
	} else {
		_ = sqlite.createMetaCsv(&model.Meta{
			TableName:   tableName,
			FileName:    descriptor.Filename,
			FileSize:    descriptor.fileSize,
			FileModTime: descriptor.fileModTime,
		})
	}

	// NewRead header
	// TODO: we should somehow handle the situation when there is no header line
	header, err := reader.csv.Read()
	if err != nil {
		sqlite.logger.Error("Failed to read the header line", "error", err.Error(), "filename", descriptor.Filename)
		return err
	}

	// Auto detect column types by the first row with data
	// Keep in mind that in case the absence of data the type will be detected incorrectly
	// In such edge situations, it would be better explicitly define column-type at the data source settings page
	firstRow, err := reader.csv.Read()
	if err != nil {
		sqlite.logger.Error("Failed to read the first data line", "error", err.Error(), "filename", descriptor.Filename)
		return err
	}
	if descriptor.Columns == nil || len(descriptor.Columns) == 0 {
		columnTypesStr := make([]string, 0)
		descriptor.Columns = make([]Column, 0)
		for i, firstRowVal := range firstRow {
			columnName := header[i]
			columnType := detectDatatype(firstRowVal)
			descriptor.Columns = append(descriptor.Columns, Column{
				Type: columnType,
				Name: columnName,
			})
			columnTypesStr = append(columnTypesStr, fmt.Sprintf("[%s](%s)", columnName, columnType))
		}

		sqlite.logger.Info("CSV column types have been auto-detected", "filename", descriptor.Filename, "columns", strings.Join(columnTypesStr, ","))
	}

	// Build map: ColumnName -> CSV column Id
	csvColumns := getColumnNames(descriptor.Columns)
	columnsMap := make(map[string]int)
	for _, columnName := range csvColumns {
		for hci, headerColumn := range header {
			if headerColumn == columnName {
				columnsMap[columnName] = hci
			}
		}
	}

	if reload {
		if err := sqlite.exec(fmt.Sprintf("DELETE FROM %s", tableName)); err != nil {
			return err
		}
	} else {
		if err := sqlite.exec(createTableFor(tableName, descriptor.Columns)); err != nil {
			return err
		}
	}

	// Prepare INSERT statement
	sqlInsert := createInsertFor(tableName, csvColumns)
	stmt, err := sqlite.db.Prepare(sqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	sqlite.logger.Debug("Begin inserting", "table", tableName, "filename", descriptor.Filename)
	insertedCount := 1

	// Insert the first row
	rowValues := valuesToRow(firstRow, descriptor.Columns, columnsMap)
	_, err = stmt.Exec(rowValues...)
	if err != nil {
		return err
	}

	// Insert rows...
	for {
		row, err := reader.csv.Read()
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			break
		}

		// CSV Row -> Insert values
		rowValues := valuesToRow(row, descriptor.Columns, columnsMap)
		_, err = stmt.Exec(rowValues...)
		if err != nil {
			return err
		}

		insertedCount++
	}

	sqlite.logger.Debug("Stop inserting", "table", tableName, "inserted", insertedCount, "filename", descriptor.Filename)

	return nil
}

func (sqlite *DbSqlite) exec(sql string) error {
	sqlite.logger.Debug("Execute", "sql", sql)
	_, err := sqlite.db.Exec(sql)
	if err != nil {
		sqlite.logger.Error("Execution failed", "sql", sql, "error", err.Error())
		return err
	}
	return nil
}

func (sqlite *DbSqlite) ifTableExists(tableName string) (bool, error) {
	rows, err := sqlite.db.Query(fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s' LIMIT 1", tableName))
	if err != nil {
		sqlite.logger.Error("Failed to check existence of table", "table", tableName, "error", err.Error())
		return false, err
	}
	defer rows.Close()

	var count int
	for {
		if !rows.Next() {
			break
		}
		count++
	}
	sqlite.logger.Debug("Table count", "count", count, "table", tableName)

	return count > 0, nil
}

func (sqlite *DbSqlite) createMetaCsvTable() error {
	metaColumns := make([]Column, 0)
	metaColumns = append(metaColumns, Column{
		Type: "TEXT",
		Name: "table_name",
	})
	metaColumns = append(metaColumns, Column{
		Type: "TEXT",
		Name: "file_name",
	})
	metaColumns = append(metaColumns, Column{
		Type: "INTEGER",
		Name: "file_size",
	})
	metaColumns = append(metaColumns, Column{
		Type: "INTEGER",
		Name: "file_mod_time",
	})
	if err := sqlite.exec(createTableFor(metaCsvTable, metaColumns)); err != nil {
		return err
	}
	return nil
}

// TODO: move to Repository
func (sqlite *DbSqlite) getMetaCsv(tableName string) *model.Meta {
	rows, err := sqlite.db.Query(
		fmt.Sprintf(
			"SELECT file_name, file_size, file_mod_time FROM %s WHERE table_name='%s' LIMIT 1",
			metaCsvTable,
			tableName,
		),
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var fileName string
		var fileSize int64
		var fileModTime int64

		err := rows.Scan(&fileName, &fileSize, &fileModTime)
		if err != nil {
			return nil
		}

		return &model.Meta{
			TableName:   tableName,
			FileName:    fileName,
			FileSize:    fileSize,
			FileModTime: fileModTime,
		}
	}

	return nil
}

// TODO: move to Repository
func (sqlite *DbSqlite) updateMetaCsv(meta *model.Meta) error {
	return sqlite.exec(
		fmt.Sprintf(
			"UPDATE %s SET file_size=%d, file_mod_time=%d WHERE table_name='%s'",
			metaCsvTable,
			meta.FileSize,
			meta.FileModTime,
			meta.TableName,
		),
	)
}

// TODO: move to Repository
func (sqlite *DbSqlite) createMetaCsv(meta *model.Meta) error {
	return sqlite.exec(
		fmt.Sprintf(
			"INSERT INTO %s VALUES('%s', '%s', %d, %d)",
			metaCsvTable,
			meta.TableName,
			meta.FileName,
			meta.FileSize,
			meta.FileModTime,
		),
	)
}

func createTableFor(tableName string, columns []Column) string {
	columnDefs := make([]string, 0)

	for _, column := range columns {
		columnDefs = append(columnDefs,
			// column data_type DEFAULT 0
			fmt.Sprintf("%s %s %s", column.Name, column.Type, getDefaultForColumn(column.Type)),
		)
	}

	return fmt.Sprintf("CREATE TABLE %s(%s)", tableName, strings.Join(columnDefs, ","))
}

func getDefaultForColumn(columnType ColumnType) string {
	switch columnType {
	case ColumnTypeReal:
		return "DEFAULT 0"
	case ColumnTypeInteger:
		return "DEFAULT 0"
	case ColumnTypeText:
		return "DEFAULT \"\""
	case ColumnTypeDate:
		return "DEFAULT CURRENT_TIMESTAMP"
	case ColumnTypeTimestamp:
		return "DEFAULT CURRENT_TIMESTAMP"
	}
	return "DEFAULT 0"
}

func getColumnNames(columns []Column) []string {
	columnNames := make([]string, 0)
	for _, column := range columns {
		columnNames = append(columnNames, column.Name)
	}
	return columnNames
}

func createInsertFor(tableName string, columnNames []string) string {
	binds := strings.TrimSuffix(strings.Repeat("?,", len(columnNames)), ",")
	return fmt.Sprintf("INSERT INTO %s (%s) values(%s)", tableName, strings.Join(columnNames, ","), binds)
}

func getColumnType(columns []Column, columnName string) *ColumnType {
	for _, column := range columns {
		if column.Name == columnName {
			return &column.Type
		}
	}
	return nil
}

// Caveat: function is not able to guess timestamp format, it will always be Integer
func detectDatatype(value string) ColumnType {
	if util.IsNumber(value) {
		if util.IsInt(value) {
			return ColumnTypeInteger
		}
		return ColumnTypeReal
	}
	_, err := dateparse.ParseAny(value)
	if err == nil {
		return ColumnTypeDate
	}
	return ColumnTypeText
}

func valuesToRow(values []string, columns []Column, columnsMap map[string]int) []interface{} {
	rowValues := make([]interface{}, 0)

	for _, column := range columns {
		if columnIndex, ok := columnsMap[column.Name]; ok {
			columnType := getColumnType(columns, column.Name)
			rowValues = append(rowValues, strToValue(values[columnIndex], columnType))
		}
	}

	return rowValues
}

func strToValue(value string, columnType *ColumnType) interface{} {
	if columnType == nil {
		return value
	}
	switch *columnType {
	case ColumnTypeDate:
		t, err := dateparse.ParseAny(value)
		if err != nil {
			return value
		}
		return t
	case ColumnTypeTimestamp:
		ival, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return value
		}
		return ival
	case ColumnTypeInteger:
		ival, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return value
		}
		return ival
	case ColumnTypeReal:
		fval, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return value
		}
		return fval
	}
	return value
}
