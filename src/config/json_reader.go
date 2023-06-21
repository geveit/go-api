package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonReader struct {
	data map[string]any
}

func (j *JsonReader) readConfigJson() (map[string]any, error) {
	if j.data != nil {
		return j.data, nil
	}

	filePath := "./config.json"
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config.json: %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("Failed to unmarshall json data: %w", err)
	}

	j.data = data

	return data, nil
}

var jsonReader *JsonReader

func getJsonReader() *JsonReader {
	if jsonReader != nil {
		return jsonReader
	}

	jsonReader = &JsonReader{}
	return jsonReader
}
