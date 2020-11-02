package service

import (
	"github.com/hashicorp/go-hclog"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
)

type ILoggerService interface {
	IsDebug() bool
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

type loggerService struct {
	args map[string]interface{}
	logger hclog.Logger
}

func NewLoggerService(logger hclog.Logger, args map[string]interface{}) ILoggerService {
	return &loggerService{
		args:   args,
		logger: logger,
	}
}

func (s *loggerService) IsDebug() bool {
	return s.logger.IsDebug()
}

func (s *loggerService) Info(message string, args ...interface{}) {
	s.logger.Info(message, s.prepareArgs(args))
}

func (s *loggerService) Error(message string, args ...interface{}) {
	s.logger.Error(message, s.prepareArgs(args))
}

func (s *loggerService) Warning(message string, args ...interface{}) {
	s.logger.Warn(message, s.prepareArgs(args))
}

func (s *loggerService) Debug(message string, args ...interface{}) {
	s.logger.Debug(message, s.prepareArgs(args))
}

func (s *loggerService) prepareArgs(args ...interface{}) []interface{} {
	prepared := util.MapToArray(s.args)
	return append(prepared, args...)
}
