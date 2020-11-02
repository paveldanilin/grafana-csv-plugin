package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/model"
	"sync"
	"time"
)

type IPluginService interface {
	Query(req *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error)
}

type PluginConfig struct {
	Logger         ILoggerService
	CsvService     ICsvService
	GrafanaService IGrafanaService
}

type pluginService struct {
	logger         ILoggerService
	csvService     ICsvService
	grafanaService IGrafanaService
}

func NewPluginService(config PluginConfig) IPluginService {
	return &pluginService{
		logger:         config.Logger,
		csvService:     config.CsvService,
		grafanaService: config.GrafanaService,
	}
}

func (s *pluginService) Query(req *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	s.debugRequest(req)

	result := &datasource.DatasourceResponse{
		Results: make([]*datasource.QueryResult, 0),
	}

	if len(req.Queries) == 0 {
		s.logger.Warning("No queries, nothing to execute")
		return nil, errors.New("csv-plugin: nothing to execute")
	}

	dsModel, err := model.CreateDatasourceFrom(*req)
	if err != nil {
		s.logger.Error("Could not create datasource", "error", err.Error())
		return nil, errors.New("csv-plugin: could not create datasource")
	}
	s.logger.Debug("Datasource", dsModel.String())

	queryTime := time.Now()
	queryTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())

	var wg sync.WaitGroup
	for _, reqQuery := range req.Queries {
		wg.Add(1)
		go func (dsReq *datasource.DatasourceRequest, dsQuery *datasource.Query, dsResp *datasource.DatasourceResponse, now time.Time) {
			defer wg.Done()

			queryModel, err := model.CreateQueryFrom(*dsQuery)
			if err != nil {
				s.logger.Error("Could not create query", "error", err.Error())
				s.grafanaService.ErrorResponse(dsQuery.RefId, dsResp, err.Error())
				return
			}
			s.logger.Debug("Query", queryModel.String())

			// RefId is hardcoded in datasource.js
			if queryModel.RefID == "[tests-connection]" {
				err := s.csvService.TestConnection(dsModel)
				if err != nil {
					s.grafanaService.ErrorResponse(dsQuery.RefId, dsResp, err.Error())
				}
				return
			}

			if queryModel.IsEmpty() {
				s.logger.Warning("Skip empty query")
				s.grafanaService.ErrorResponse(dsQuery.RefId, dsResp, "Empty query")
				return
			}

			datetimeFrom, _ := dateparse.ParseAny(fmt.Sprintf("%d", dsReq.TimeRange.FromEpochMs))
			datetimeTo, _ := dateparse.ParseAny(fmt.Sprintf("%d", dsReq.TimeRange.ToEpochMs))

			// TODO: use context instead of scope obj
			queryScope := macro.NewScope()
			queryScope.SetVar("timeToMs",	dsReq.TimeRange.ToEpochMs)
			queryScope.SetVar("timeFromMs",	dsReq.TimeRange.FromEpochMs)
			queryScope.SetVar("timeFrom",	datetimeFrom)
			queryScope.SetVar("timeTo", 	datetimeTo)
			queryScope.SetVar("sysdate",      now)

			dsResp.Results = append(dsResp.Results, s.performQuery(dsModel, queryModel, queryScope))
		}(req, reqQuery, result, queryTime)
	}
	wg.Wait()

	return result, nil
}

func (s *pluginService) performQuery(dsModel *model.Datasource, queryModel *model.Query, scope *macro.Scope) *datasource.QueryResult {
	result, err := s.csvService.Query(dsModel, queryModel, scope)
	if err != nil {
		return &datasource.QueryResult{
			Error: fmt.Sprintf("Query failed: %s", err.Error()),
			RefId: queryModel.RefID,
		}
	}
	defer result.Release()

	return s.grafanaService.ToQueryResult(queryModel.RefID, queryModel.Format, result)
}

func (s *pluginService) debugRequest(req *datasource.DatasourceRequest) {
	if s.logger.IsDebug() == false {
		return
	}
	reqJson, _ := json.Marshal(req)
	s.logger.Debug("Datasource request", "jsonRequest", reqJson)
}
