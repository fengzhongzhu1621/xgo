package maps

import (
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var deploySpec = map[string]interface{}{
	"testKey":              "testValue",
	"replicas":             3,
	"revisionHistoryLimit": 10,
	"intKey4SetItem":       8,
	"selector": map[string]interface{}{
		"matchLabels": map[string]interface{}{
			"app": "nginx",
		},
	},
	"strategy": map[string]interface{}{
		"rollingUpdate": map[string]interface{}{
			"maxSurge":       "25%",
			"maxUnavailable": "25%",
		},
		"type": "RollingUpdate",
	},
	"template": map[string]interface{}{
		"metadata": map[string]interface{}{
			"creationTimestamp": nil,
			"int64Key4GetInt64": int64(10),
			"labels": map[string]interface{}{
				"app":           "nginx",
				"strKey4GetStr": "value",
			},
		},
		"spec": map[string]interface{}{
			"boolKey4GetBool": true,
			"containers": []map[string]interface{}{
				{
					"image":           "nginx:latest",
					"imagePullPolicy": "IfNotPresent",
					"name":            "nginx",
					"ports": map[string]interface{}{
						"containerPort": 80,
						"protocol":      "TCP",
					},
					"resources": map[string]interface{}{},
				},
			},
			"dnsPolicy":                     "ClusterFirst",
			"restartPolicy":                 "Always",
			"schedulerName":                 "default-scheduler",
			"securityContext":               map[string]interface{}{},
			"terminationGracePeriodSeconds": 30,
		},
		"interfaceList": []interface{}{
			map[string]interface{}{"key": "value"},
			"key-value",
		},
	},
}

// paths 为以 '.' 连接的字符串
func TestGetItems(t *testing.T) {
	// depth 1，val type int
	ret, _ := GetItems(deploySpec, "replicas")
	assert.Equal(t, 3, ret)

	// depth 2, val type string
	ret, _ = GetItems(deploySpec, "strategy.type")
	assert.Equal(t, "RollingUpdate", ret)

	// depth 3, val type string
	ret, _ = GetItems(deploySpec, "template.spec.restartPolicy")
	assert.Equal(t, "Always", ret)
}

// paths 为 []string，成功的情况
func TestGetItemsSuccessCase(t *testing.T) {
	// depth 1，val type int
	ret, _ := GetItems(deploySpec, []string{"replicas"})
	assert.Equal(t, 3, ret)

	// depth 2，val type map[string]interface{}
	r, _ := GetItems(deploySpec, []string{"selector", "matchLabels"})
	_, ok := r.(map[string]interface{})
	assert.Equal(t, true, ok)

	// depth 2, val type string
	ret, _ = GetItems(deploySpec, []string{"strategy", "type"})
	assert.Equal(t, "RollingUpdate", ret)

	// depth 3, val type nil
	ret, _ = GetItems(deploySpec, []string{"template", "metadata", "creationTimestamp"})
	assert.Nil(t, ret)

	// depth 3, val type string
	ret, _ = GetItems(deploySpec, []string{"template", "spec", "restartPolicy"})
	assert.Equal(t, "Always", ret)
}

// paths 为 []string 或 其他，失败的情况
func TestGetItemsFailCase(t *testing.T) {
	// invalid paths type error
	_, err := GetItems(deploySpec, 0)
	assert.True(t, errors.Is(err, ErrInvalidPathType))

	// not paths error
	_, err = GetItems(deploySpec, []string{})
	assert.NotNil(t, err)

	// not map[string]interface{} type error
	_, err = GetItems(deploySpec, []string{"replicas", "testKey"})
	assert.NotNil(t, err)

	_, err = GetItems(deploySpec, []string{"template", "spec", "containers", "image"})
	assert.NotNil(t, err)

	// key not exist
	_, err = GetItems(deploySpec, []string{"templateKey", "spec"})
	assert.NotNil(t, err)

	_, err = GetItems(deploySpec, []string{"selector", "spec"})
	assert.NotNil(t, err)

	// paths type error
	_, err = GetItems(deploySpec, []int{123, 456})
	assert.NotNil(t, err)

	_, err = GetItems(deploySpec, 123)
	assert.NotNil(t, err)
}

func TestGet(t *testing.T) {
	ret := Get(deploySpec, []string{"replicas"}, 1)
	assert.Equal(t, 3, ret)

	ret = Get(deploySpec, []string{}, nil)
	assert.Nil(t, ret)

	ret = Get(deploySpec, "container.name", "defaultName")
	assert.Equal(t, "defaultName", ret)
}

func TestGetBool(t *testing.T) {
	assert.True(t, GetBool(deploySpec, "template.spec.boolKey4GetBool"))
	assert.False(t, GetBool(deploySpec, "template.spec.notExistsKey"))
}

func TestGetInt64(t *testing.T) {
	assert.Equal(t, int64(10), GetInt64(deploySpec, "template.metadata.int64Key4GetInt64"))
	assert.Equal(t, int64(0), GetInt64(deploySpec, "template.spec.notExistsKey"))
}

func TestGetStr(t *testing.T) {
	assert.Equal(t, "value", GetStr(deploySpec, "template.metadata.labels.strKey4GetStr"))
	assert.Equal(t, "default-scheduler", GetStr(deploySpec, "template.spec.schedulerName"))
	assert.Equal(t, "", GetStr(deploySpec, "template.spec.notExistsKey"))
}

func TestGetList(t *testing.T) {
	assert.Equal(
		t, []interface{}{map[string]interface{}{"key": "value"}, "key-value"},
		GetList(deploySpec, "template.interfaceList"),
	)
	assert.Equal(t, []interface{}{}, GetList(deploySpec, "template.spec.notExistsKey"))
}

func TestGetMap(t *testing.T) {
	assert.Equal(
		t,
		map[string]interface{}{"app": "nginx"},
		GetMap(deploySpec, "selector.matchLabels"),
	)
	assert.Equal(t, map[string]interface{}{}, GetMap(deploySpec, "template.spec.notExistsKey"))
}

// 根据 value 值查找 key 值
func TestFindKey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1, ok1 := lo.FindKey(map[string]int{"foo": 1, "bar": 2, "baz": 3}, 2)
	is.Equal("bar", result1)
	is.True(ok1)

	result2, ok2 := lo.FindKey(map[string]int{"foo": 1, "bar": 2, "baz": 3}, 42)
	is.Equal("", result2)
	is.False(ok2)

	type test struct {
		foobar string
	}

	result3, ok3 := lo.FindKey(
		map[string]test{"foo": {"foo"}, "bar": {"bar"}, "baz": {"baz"}},
		test{"foo"},
	)
	is.Equal("foo", result3)
	is.True(ok3)

	result4, ok4 := lo.FindKey(
		map[string]test{"foo": {"foo"}, "bar": {"bar"}, "baz": {"baz"}},
		test{"hello world"},
	)
	is.Equal("", result4)
	is.False(ok4)
}

// 根据条件匹配查询 key值
func TestFindKeyBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1, ok1 := lo.FindKeyBy(
		map[string]int{"foo": 1, "bar": 2, "baz": 3},
		func(k string, v int) bool {
			return k == "foo"
		},
	)
	is.Equal("foo", result1)
	is.True(ok1)

	result2, ok2 := lo.FindKeyBy(
		map[string]int{"foo": 1, "bar": 2, "baz": 3},
		func(k string, v int) bool {
			return false
		},
	)
	is.Equal("", result2)
	is.False(ok2)
}
