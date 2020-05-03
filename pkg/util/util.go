package util

import (
	"errors"
	"fmt"
	"os"
)

func MapToArray(m map[string]interface{}) []interface{} {
	args := make([]interface{}, 0)
	for k, v := range m {
		args = append(args, k)
		args = append(args, v)
	}
	return args
}

func GetStr(name string, data map[string]interface{}, def string) string {
	if val, ok := data[name]; ok {
		return val.(string)
	}
	return def
}

func GetBool(name string, data map[string]interface{}, def bool) bool {
	if val, ok := data[name]; ok {
		return val.(bool)
	}
	return def
}

func CheckFile(fileName string) error {
	if len(fileName) == 0 {
		return errors.New("the path to file is not defined")
	}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("file not found `%s`", fileName))
	}
	if _, err := os.Stat(fileName); os.IsPermission(err) {
		return errors.New(fmt.Sprintf("there is no access to file `%s`", fileName))
	}
	return nil
}
