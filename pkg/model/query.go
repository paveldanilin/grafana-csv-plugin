package model

import (
	"encoding/json"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
)

type QueryModel struct {
	RefID string
	Format string
	ColumnSort string
	CustomColumnOrder []string
}

func CreateQueryFrom(query datasource.Query) (*QueryModel, error) {
	inputModel := make(map[string]interface{})
	err := json.Unmarshal([]byte(query.ModelJson), &inputModel)
	if err != nil {
		return nil, err
	}

	tmpCustomColumnOrder := inputModel["customColumnOrder"].([]interface{})
	customColumnOrder := make([]string, 0)
	for _, v := range tmpCustomColumnOrder {
		customColumnOrder = append(customColumnOrder, v.(string))
	}

	model := &QueryModel{
		RefID:             query.RefId,
		Format:            util.GetStr("format", inputModel, ""),
		ColumnSort:        util.GetStr("columnSort", inputModel, ""),
		CustomColumnOrder: customColumnOrder,
	}

	return model, nil
}
