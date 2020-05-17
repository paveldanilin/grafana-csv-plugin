package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func MapToArray(m map[string]interface{}) []interface{} {
	args := make([]interface{}, 0)
	for k, v := range m {
		args = append(args, k)
		args = append(args, v)
	}
	return args
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

func IsNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

func IsInt(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err == nil
}
