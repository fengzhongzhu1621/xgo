package flagutils

import (
	"strings"

	"github.com/spf13/pflag"
)

// WordSepNormalizeFunc 将命令行参数名称中的下划线统一替换为短横线，以便在不同平台和工具中保持一致性。
// e.g. my_flag_name -> my-flag-name
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}
