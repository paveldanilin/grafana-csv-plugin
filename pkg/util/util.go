package util

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
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

func FileChanged(filename string, oldSize, oldModTime int64) bool {
	curSize, curModTime := FileStat(filename)
	return curSize != oldSize || curModTime != oldModTime
}

// Returns fileSize, fileModTime
func FileStat(filename string) (int64, int64) {
	fileStat, err := os.Stat(filename)
	if err != nil {
		return 0, 0
	}
	return fileStat.Size(), fileStat.ModTime().Unix()
}

func TimeToEpochMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond) // 1000000
}

func ToFloat64(v interface{}) (float64, error) {
	if v == nil {
		return float64(0), nil
	}

	switch i := v.(type) {
	case string:
		return strconv.ParseFloat(i, 64)
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int:
		return float64(i), nil
	case int32:
		return float64(i), nil
	}
	return float64(0), errors.New(fmt.Sprintf("Could not convert value `%v`.(%s) => float64", v, reflect.TypeOf(v).String()))
}
