package resource_watch

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo/ginx/apps/cmdb/resource_watch/watch"
)

var cursorTypeMap = map[ResType]watch.CursorType{
	Host:             watch.Host,
	ModuleHostRel:    watch.ModuleHostRelation,
	Biz:              watch.Biz,
	Set:              watch.Set,
	Module:           watch.Module,
	Process:          watch.Process,
	ProcessRelation:  watch.ProcessInstanceRelation,
	BizSet:           watch.BizSet,
	Plat:             watch.Plat,
	Project:          watch.Project,
	ObjectInstance:   watch.ObjectBase,
	MainlineInstance: watch.MainlineInstance,
	InstAsst:         watch.InstAsst,
	KubeCluster:      watch.KubeCluster,
	KubeNode:         watch.KubeNode,
	KubeNamespace:    watch.KubeNamespace,
	KubeWorkload:     watch.KubeWorkload,
	KubePod:          watch.KubePod,
}

// GetCursorTypeByResType get event watch cursor type by resource type
func GetCursorTypeByResType(res ResType) (watch.CursorType, error) {
	typ, exists := cursorTypeMap[res]
	if !exists {
		return watch.UnknownType, fmt.Errorf("resource type %s has no matching cursor type", res)
	}

	return typ, nil
}
