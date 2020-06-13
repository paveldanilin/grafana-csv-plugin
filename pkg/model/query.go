package model

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"strings"
)

type Query struct {
	RefID	string	`json:"refId,omitempty"`
	Format	string	`json:"format"`
	Query	string	`json:"query"`
}

func CreateQueryFrom(query datasource.Query) (*Query, error) {
	model := &Query{}
	err := json.Unmarshal([]byte(query.ModelJson), &model)
	if err != nil {
		return nil, err
	}
	model.RefID = query.RefId


	model.Query = strings.TrimSpace(model.Query)
	if len(model.Query) == 0 {
		model.Query = fmt.Sprintf("SELECT * FROM %s LIMIT 1, 15", csv.TableName)
	}

	return model, nil
}

func (m *Query) String() string {
	jsonBytes, _ := json.Marshal(m)
	return string(jsonBytes)
}
