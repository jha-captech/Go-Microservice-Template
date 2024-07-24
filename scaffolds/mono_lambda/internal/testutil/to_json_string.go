package testutil

import (
	"encoding/json"
	"fmt"
)

func ToJSONString(data any) string {
	JSONString, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("Failed to Marshal data to JSON. \ndata: %v\nerr: %v", data, err))
	}
	return string(JSONString)
}
