package service

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/model"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/xsql"
	"time"
)

type ICsvService interface {
	TestConnection(ds *model.Datasource) error
	Query(ds *model.Datasource, q *model.Query, scope *macro.Scope) (*xsql.QueryResult, error)
}

type csvService struct {
	db csv.DB
	logger ILoggerService
}

func NewCsvService(db csv.DB, logger ILoggerService) ICsvService {
	return &csvService{db: db, logger: logger}
}

func (s *csvService) TestConnection(ds *model.Datasource) error {
	return util.CheckFile(ds.Filename)
}

func (s *csvService) Query(ds *model.Datasource, q *model.Query, scope *macro.Scope) (*xsql.QueryResult, error) {
	tableColumns := make([]csv.Column, 0)
	for _, dsColumn := range ds.Columns {
		tableColumns = append(tableColumns, csv.Column{
			Type: csv.ColumnTypeFromString(dsColumn.Type),
			Name: dsColumn.Name,
		})
	}

	err := s.db.LoadCSV(ds.Name, &csv.FileDescriptor{
		Filename:         ds.Filename,
		Delimiter:        rune(ds.CsvDelimiter[0]),
		Comment:          rune(ds.CsvComment[0]),
		TrimLeadingSpace: ds.CsvTrimLeadingSpace,
		FieldsPerRecord:  0, // Implies that each row contains the same count of fields as the header row
		Columns: tableColumns,
	})
	if err != nil {
		return nil, err
	}

	interpolated, err := macro.Interpolate(q.Query, scope)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Query(interpolated.Text(), s.prepareGeneratedColumns(interpolated.MetaList(), scope))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *csvService) prepareGeneratedColumns(metaList []macro.Meta, scope *macro.Scope) []xsql.GeneratedColumn {
	generatedColumns := make([]xsql.GeneratedColumn, 0)
	for _, meta := range metaList {
		generator := s.meta2column(meta, scope)
		if generator != nil {
			generatedColumns = append(generatedColumns, *generator)
		}
	}
	return generatedColumns
}

func (s *csvService) meta2column(meta macro.Meta, scope *macro.Scope) *xsql.GeneratedColumn {
	// TODO: refactor
	switch meta.Name {
	case "$autoTime":
		// 1m
		// 1ms
		// 1h
		step := "1m"
		if len(meta.Options) >= 1 {
			step = meta.Options[0]
		}
		stepTime := util.StrToDur(step)
		return &xsql.GeneratedColumn{
			Name:    "autotime",
			Func: func(values []interface{}) interface{} {
				if scope.HasVar("autotime") {
					t := scope.GetVar("autotime").(time.Time)
					scope.SetVar("autotime", t.Add(stepTime))
				} else {
					scope.SetVar("autotime", scope.GetVar("sysdate").(time.Time))
				}
				return scope.GetVar("autotime").(time.Time)
			},
			Type: nil,
		}
	}
	return nil
}
