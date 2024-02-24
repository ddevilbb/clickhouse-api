package model

import (
	"clickhouse-api/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// swagger:model TestData
type TestData struct {
	// Id of TestData
	// in: string
	Id string `json:"id" db:"id"`

	// Sign of TestData
	// in: int8
	Sign int8 `json:"sign" db:"sign"`

	// Version of TestData
	// in: int32
	Version uint32 `json:"version" db:"version"`

	// Data of TestData
	// in: string
	Data string `json:"data" db:"data"`

	// CreatedAt of TestData
	// in: time.Time
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (testData *TestData) IsValid() (map[string]interface{}, bool) {
	if testData.Data == "" {
		return helper.Message(false, helper.FieldIsRequired("data")), false
	}

	return helper.Message(true, "success"), true
}

func CreateTestDataFromRequestBody(req *http.Request) (*TestData, error) {
	testData := &TestData{}
	err := json.NewDecoder(req.Body).Decode(testData)

	return testData, err
}

func CreateTestDataFromJsonString(jsonString string) (*TestData, error) {
	testData := &TestData{}
	if err := json.Unmarshal([]byte(jsonString), testData); err != nil {
		return nil, fmt.Errorf("json decode error: %+v", err)
	}

	return testData, nil
}
