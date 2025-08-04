package json

import (
	"encoding/json"

	"github.com/apapsch/go-jsonmerge/v2"
)

// JSONMerge merges two JSON representation into a single object. `data` is the
// existing representation and `patch` is the new data to be merged in
// 合并两个 JSON 表示为一个单一对象。`data` 是现有的表示，`patch` 是要合并的新数据。
// 如果 `data` 或 `patch` 为 nil，则视为空对象 {}。
func JSONMerge(data, patch json.RawMessage) (json.RawMessage, error) {
	// 设置置了合并器的配置选项 CopyNonexistent 为 true。这意味着在合并过程中，
	// patch 中存在但 data 中不存在的属性会被复制到最终的合并结果中。
	// data 中存在但 patch 中不存在的属性会保留下来
	// 确保在合并过程中，patch 中新增的属性会被复制到最终的合并结果中，而不会因为 data 中不存在而丢失。
	merger := jsonmerge.Merger{
		CopyNonexistent: true,
	}

	// 如果 data 为 nil，使用空对象
	if data == nil {
		data = []byte(`{}`)
	}

	// 如果 patch 为 nil，使用空对象
	if patch == nil {
		patch = []byte(`{}`)
	}

	// 将 data 和 patch 进行合并，得到合并后的 JSON 数据
	merged, err := merger.MergeBytes(data, patch)
	if err != nil {
		return nil, err
	}
	return merged, nil
}
