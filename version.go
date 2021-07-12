package xgo

import (
	"fmt"
	"runtime"
)

// Version is the current package version.
const version = "0.0.1"

// This value can be set for releases at build time using:
//   go {build|run} -ldflags "-X xgo.version 1.2.3 -X xgo.buildInfo timestamp-@githubuser-platform".
// If unset, Version() shall return "DEVBUILD".
// var version string = "DEVBUILD"
var buildInfo = "--"

func Version() string {
	return fmt.Sprintf("%s (%s %s-%s) %s", version, runtime.Version(), runtime.GOOS, runtime.GOARCH, buildInfo)
}
