package nethttp

import (
	"encoding/json"
	"os"
	"strings"
)

type ApiResponse struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateApiResponse(code int, message string, data interface{}) ([]byte, error) {
	result := false
	if 0 == code {
		result = true
	} else {
		appName := os.Args[0]
		szArr := strings.Split(appName, "/")
		if len(szArr) >= 2 {
			appName = szArr[1]
		}
		message = "(" + appName + "):" + message
	}

	resp := ApiResponse{result, code, message, data}
	b, err := json.Marshal(resp)
	if err != nil {
		return []byte(""), err
	}

	return b, nil
}
