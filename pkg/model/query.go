package model

import (
	"encoding/json"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

type QueryModel struct {
	RefID	string	`json:"refId,omitempty"`
	Format	string	`json:"format"`
	Query	string	`json:"query"`
}

func CreateQueryFrom(query datasource.Query) (*QueryModel, error) {
	model := &QueryModel{}
	err := json.Unmarshal([]byte(query.ModelJson), &model)
	if err != nil {
		return nil, err
	}
	model.RefID = query.RefId
	return model, nil
}

func (m *QueryModel) String() string {
	jsonBytes, _ := json.Marshal(m)
	return string(jsonBytes)
}
