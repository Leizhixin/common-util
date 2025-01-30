package fileutil

import (
	"encoding/json"
	"github.com/Leizhixin/common-util/middleware/beanutil"
	"os"
)

func WriteResultToFile(raw interface{}, resFileName string) error {
	resultFile, err := os.Create(resFileName)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(resultFile)
	switch raw.(type) {
	case map[interface{}]interface{}:
		raw = beanutil.ParseResultToMap(raw.(map[interface{}]interface{}))
	case []interface{}:
		rawList := raw.([]interface{})
		beanutil.HandleList(rawList)
	}
	content, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	_, err = resultFile.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func WriteStringResultToFile(content string, resFileName string) error {
	resultFile, err := os.Create(resFileName)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(resultFile)
	_, err = resultFile.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}

func ReadJsonFileContent(fileName string, ptr interface{}) error {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, ptr)
	if err != nil {
		return err
	}
	return nil
}
