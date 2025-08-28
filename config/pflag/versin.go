package pflag

import "github.com/spf13/pflag"

// AddCommonFlags add common flags that is needed by all modules
func AddCommonFlags(cmdline *pflag.FlagSet) *bool {
	version := cmdline.Bool("version", false, "show version information")
	return version
}
