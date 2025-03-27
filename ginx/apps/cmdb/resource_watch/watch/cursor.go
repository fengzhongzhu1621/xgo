package watch

import (
	"github.com/fengzhongzhu1621/xgo/ginx/apps/cmdb/resource_watch/stream"
)

// NoEventCursor TODO
// this cursor means there is no event occurs.
// we send this cursor to our the watcher and if we
// received a NoEventCursor, then we need to fetch event
// from the head cursor
var NoEventCursor string

type CursorType string

const (
	// NoEvent TODO
	NoEvent CursorType = "no_event"
	// UnknownType TODO
	UnknownType CursorType = "unknown"
	// Host TODO
	Host CursorType = "host"
	// ModuleHostRelation TODO
	ModuleHostRelation CursorType = "host_relation"
	// Biz TODO
	Biz CursorType = "biz"
	// Set TODO
	Set CursorType = "set"
	// Module TODO
	Module CursorType = "module"
	// ObjectBase TODO
	ObjectBase CursorType = "object_instance"
	// Process TODO
	Process CursorType = "process"
	// ProcessInstanceRelation TODO
	ProcessInstanceRelation CursorType = "process_instance_relation"
	// BizSet TODO
	BizSet CursorType = "biz_set"
	// HostIdentifier TODO
	// a mixed event type, which contains host, host relation, process events etc.
	// and finally converted to host identifier event.
	HostIdentifier CursorType = "host_identifier"
	// MainlineInstance specified for mainline instance event watch, filtered from object instance events
	MainlineInstance CursorType = "mainline_instance"
	// InstAsst specified for instance association event watch
	InstAsst CursorType = "inst_asst"
	// BizSetRelation a mixed event type containing biz set & biz events, which are converted to their relation events
	BizSetRelation CursorType = "biz_set_relation"
	// Plat cloud area event cursor type
	Plat CursorType = "plat"
	// Project project event cursor type
	Project CursorType = "project"
	// kube related cursor types
	// KubeCluster cursor type
	KubeCluster CursorType = "kube_cluster"
	// KubeNode cursor type
	KubeNode CursorType = "kube_node"
	// KubeNamespace cursor type
	KubeNamespace CursorType = "kube_namespace"
	// KubeWorkload cursor type, including all workloads(e.g. deployment) with their type specified in sub-resource
	KubeWorkload CursorType = "kube_workload"
	// KubePod cursor type, its event detail is pod info with containers in it
	KubePod CursorType = "kube_pod"
)

var (
	cursorTypeIntMap = map[CursorType]int{
		NoEvent:                 1,
		Host:                    2,
		ModuleHostRelation:      3,
		Biz:                     4,
		Set:                     5,
		Module:                  6,
		ObjectBase:              8,
		Process:                 9,
		ProcessInstanceRelation: 10,
		HostIdentifier:          11,
		MainlineInstance:        12,
		InstAsst:                13,
		BizSet:                  14,
		BizSetRelation:          15,
		Plat:                    16,
		KubeCluster:             17,
		KubeNode:                18,
		KubeNamespace:           19,
		KubeWorkload:            20,
		KubePod:                 21,
		Project:                 22,
	}

	intCursorTypeMap = make(map[int]CursorType)
)

func (ct CursorType) ToInt() int {
	intVal, exists := cursorTypeIntMap[ct]
	if !exists {
		return -1
	}

	return intVal
}

// ParseInt TODO
func (ct *CursorType) ParseInt(typ int) {
	cursorType, exists := intCursorTypeMap[typ]
	if !exists {
		*ct = UnknownType
		return
	}

	*ct = cursorType
}

// ListCursorTypes returns all support CursorTypes.
func ListCursorTypes() []CursorType {
	return []CursorType{Host, ModuleHostRelation, Biz, Set, Module, ObjectBase, Process, ProcessInstanceRelation,
		HostIdentifier, MainlineInstance, InstAsst, BizSet, BizSetRelation, Plat, KubeCluster, KubeNode, KubeNamespace,
		KubeWorkload, KubePod, Project}
}

type Cursor struct {
	Type        CursorType
	ClusterTime stream.TimeStamp
	// a random hex string to avoid the caller to generated a self-defined cursor.
	Oid  string
	Oper stream.OperType
	// UniqKey is an optional key which is used to ensure that a cursor is unique among a same resources(
	// as is Type field).
	// This key is used when the upper fields can not describe a unique cursor. such as the common object instance
	// event, all these instance event all have a same cursor type, as is object instance. but the instance id is
	// unique which can be used as a unique key, and is REENTRANT when a retry operation happens which is
	// **VERY IMPORTANT**.
	UniqKey string
}
