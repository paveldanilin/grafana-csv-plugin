package model

import (
	"encoding/json"
	"errors"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

type Datasource struct {
	ID			int64	`json:"id,omitempty"`
	OrgID			int64	`json:"orgId,omitempty"`
	Name			string	`json:"name,omitempty"`
	Type			string	`json:"type:omitempty"`

	Filename		string	`json:"filename"`

	// CSV options
	CsvDelimiter		string	`json:"csvDelimiter"`
	CsvComment		string	`json:"csvComment"`
	CsvTrimLeadingSpace	bool	`json:"csvTrimLeadingSpace"`

	Columns			[]struct {
		Name	string	`json:"name"`
		Type	string	`json:"type"`
	} `json:"columns"`
}

func CreateDatasourceFrom(req datasource.DatasourceRequest) (*Datasource, error) {
	model := &Datasource{}
	err := json.Unmarshal([]byte(req.Datasource.JsonData), &model)
	if err != nil {
		return nil, err
	}

	model.ID = req.Datasource.Id
	model.OrgID = req.Datasource.OrgId
	model.Name = req.Datasource.Name
	model.Type = req.Datasource.Type

	if len(model.CsvDelimiter) == 0 {
		model.CsvDelimiter = ","
	}

	if len(model.CsvComment) == 0 {
		model.CsvComment = "#"
	}

  	return model, validateLocalDatasourceModel(model)
}

func (m *Datasource) String() string {
	jsonBytes, _ := json.Marshal(m)
	return string(jsonBytes)
}

func validateLocalDatasourceModel(model *Datasource) error {
	if len(model.Filename) == 0 {
		return errors.New("the path to the CSV file is not defined")
	}
	return nil
}
