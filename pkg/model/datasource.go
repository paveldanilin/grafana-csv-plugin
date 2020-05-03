package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
)

const (
	AccessMode_LOCAL = "local"
	AccessMode_SFTP  = "sftp"
)

type DatasourceModel struct {
	ID             int64
	OrgID          int64
	Name           string
	Type           string

	Filename string

	// CSV options
	CsvDelimiter string
	CsvComment string
	CsvTrimLeadingSpace bool

	// Access mode: local, sftp
	AccessMode string

	// SFTP
	SftpHost string
	SftpPort string
	SftpUser string
	SftpPassword string
	SftpWorkingDir string // Local working dir
	SftpIgnoreHostKey bool
}

func CreateDatasourceFrom(req datasource.DatasourceRequest) (*DatasourceModel, error) {
	inputModel := make(map[string]interface{})
	err := json.Unmarshal([]byte(req.Datasource.JsonData), &inputModel)
	if err != nil {
		return nil, err
	}

	model := &DatasourceModel{
		ID:            		req.Datasource.Id,
		OrgID:         		req.Datasource.OrgId,
		Name:          		req.Datasource.Name,
		Type:          		req.Datasource.Type,
		// Access Mode
		AccessMode: 		util.GetStr("accessMode", inputModel, "local"),
		// CSV options
		CsvDelimiter: 		util.GetStr("csvDelimiter", inputModel, ","),
		CsvComment: 		util.GetStr("csvComment", inputModel, "#"),
		CsvTrimLeadingSpace:	util.GetBool("csvTrimLeadingSpace", inputModel, true),

		Filename: util.GetStr("filename", inputModel, ""),

		// SFTP
		SftpHost: 		util.GetStr("sftpHost", inputModel, ""),
		SftpPort: 		util.GetStr("sftpPort", inputModel, ""),
		SftpUser: 		util.GetStr("sftpUser", inputModel, ""),
		SftpPassword:		req.Datasource.DecryptedSecureJsonData["sftpPassword"],
		SftpWorkingDir:         util.GetStr("sftpWorkingDir", inputModel, ""),
		SftpIgnoreHostKey:      util.GetBool("sftpIgnoreHostKey", inputModel, false),
	}

	if len(model.CsvDelimiter) == 0 {
		model.CsvDelimiter = ","
	}

	if len(model.CsvComment) == 0 {
		model.CsvComment = "#"
	}

	if model.AccessMode == AccessMode_SFTP {
		if len(model.SftpPort) == 0 {
			model.SftpPort = "22"
		}
	}

  	return model, validateDatasourceModel(model)
}

func validateDatasourceModel(model *DatasourceModel) error {
	if model.AccessMode == AccessMode_LOCAL {
		return validateLocalDatasourceModel(model)
	}
	if model.AccessMode == AccessMode_SFTP {
		return validateSftpDatasourceModel(model)
	}
	return errors.New(fmt.Sprintf("unknown access mode `%s`", model.AccessMode))
}

func validateLocalDatasourceModel(model *DatasourceModel) error {
	if len(model.Filename) == 0 {
		return errors.New("the path to the CSV file is not defined")
	}
	return nil
}

func validateSftpDatasourceModel(model *DatasourceModel) error {
	if len(model.Filename) == 0 {
		return errors.New("the path to the CSV file is not defined")
	}
	if len(model.SftpHost) == 0 {
		return errors.New("SFTP hot is not defined")
	}
	if len(model.SftpWorkingDir) == 0 {
		return errors.New("working dir is not defined")
	}
	return nil
}
