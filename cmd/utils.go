package cmd

import (
	"encoding/json"
	"time"
)

func jsonUnmarshalImpl(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func currentDate() string {
	return time.Now().Format("2006-01-02")
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
