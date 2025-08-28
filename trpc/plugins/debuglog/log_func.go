package debuglog

import (
	"context"
	"encoding/json"
	"fmt"
)

// LogFunc is the struct print method function.
type LogFunc func(ctx context.Context, req, rsp any) string

// DefaultLogFunc is the default struct print method.
var DefaultLogFunc = func(ctx context.Context, req, rsp any) string {
	// 查看字段名与值的关系 {Name:Alice Age:30}
	return fmt.Sprintf(", req:%+v, rsp:%+v", req, rsp)
}

// SimpleLogFunc does not print the struct.
var SimpleLogFunc = func(ctx context.Context, req, rsp any) string {
	return ""
}

// PrettyJSONLogFunc is the method for printing formatted JSON.
var PrettyJSONLogFunc = func(ctx context.Context, req, rsp any) string {
	// 带缩进和换行的结构化 JSON，易于阅读
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	rspJSON, _ := json.MarshalIndent(rsp, "", "  ")
	return fmt.Sprintf("\nreq:%s\nrsp:%s", string(reqJSON), string(rspJSON))
}

// JSONLogFunc is the method for printing JSON.
var JSONLogFunc = func(ctx context.Context, req, rsp any) string {
	reqJSON, _ := json.Marshal(req)
	rspJSON, _ := json.Marshal(rsp)
	return fmt.Sprintf("\nreq:%s\nrsp:%s", string(reqJSON), string(rspJSON))
}

// get log func by log type
func getLogFunc(t string) LogFunc {
	switch t {
	case "simple":
		return SimpleLogFunc
	case "prettyjson":
		return PrettyJSONLogFunc
	case "json":
		return JSONLogFunc
	default:
		return DefaultLogFunc
	}
}
