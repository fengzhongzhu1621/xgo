package resource_watch

import "errors"

// ResType is the resource type for general resource cache
type ResType string

const (
	// Host is the resource type for host cache
	Host ResType = "host"
	// ModuleHostRel is the resource type for host relation cache
	ModuleHostRel ResType = "host_relation"
	// Biz is the resource type for business cache
	Biz ResType = "biz"
	// Set is the resource type for set cache
	Set ResType = "set"
	// Module is the resource type for module cache
	Module ResType = "module"
	// Process is the resource type for process cache
	Process ResType = "process"
	// ProcessRelation is the resource type for process instance relation cache
	ProcessRelation ResType = "process_relation"
	// BizSet is the resource type for  cache
	BizSet ResType = "biz_set"
	// Plat is the resource type for cloud area cache
	Plat ResType = "plat"
	// Project is the resource type for project cache
	Project ResType = "project"
	// ObjectInstance is the resource type for common object instance cache, its sub resource specifies the object id
	ObjectInstance ResType = "object_instance"
	// MainlineInstance is the resource type for mainline instance cache, its sub resource specifies the object id
	MainlineInstance ResType = "mainline_instance"
	// InstAsst is the resource type for instance association cache, its sub resource specifies the associated object id
	InstAsst ResType = "inst_asst"
	// KubeCluster is the resource type for kube cluster cache
	KubeCluster ResType = "kube_cluster"
	// KubeNode is the resource type for kube node cache
	KubeNode ResType = "kube_node"
	// KubeNamespace is the resource type for kube namespace cache
	KubeNamespace ResType = "kube_namespace"
	// KubeWorkload is the resource type for kube workload cache,  its sub resource specifies the workload type
	KubeWorkload ResType = "kube_workload"
	// KubePod is the resource type for kube pod cache, its event detail is pod info with containers in it
	KubePod ResType = "kube_pod"
)

// SupportedResTypeMap is a map whose key is resource type that is supported by general resource cache
// not all resource types are supported now, add related logics if other resource type needs cache.
var SupportedResTypeMap = map[ResType]struct{}{
	Host:             {},
	Biz:              {},
	Set:              {},
	Module:           {},
	BizSet:           {},
	Plat:             {},
	ObjectInstance:   {},
	MainlineInstance: {},
}

// ResTypeHasSubResMap is a map of supported resource type -> whether it has sub resource
var ResTypeHasSubResMap = map[ResType]struct{}{
	ObjectInstance:   {},
	MainlineInstance: {},
	InstAsst:         {},
	KubeWorkload:     {},
}

// ValidateWithSubRes validate ResType with sub resource
func (r ResType) ValidateWithSubRes(subRes string) error {
	_, exists := SupportedResTypeMap[r]
	if !exists {
		return errors.New("param is invalid")
	}

	_, hasSubRes := ResTypeHasSubResMap[r]
	if (subRes != "") != hasSubRes {
		return errors.New("param is invalid")
	}

	return nil
}

// ResTypeNeedOidMap is a map whose key is resource type that needs oid to generate id key
var ResTypeNeedOidMap = map[ResType]struct{}{
	ModuleHostRel:   {},
	ProcessRelation: {},
}
