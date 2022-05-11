package viper

import (
	"fmt"
	"strings"
	"testing"

	"xgo/cast"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type viper struct {
	pflags map[string]IFlagValue
}

func (vp *viper) bindFlagValues(flags IFlagValueSet) (err error) {
	flags.VisitAll(func(flag IFlagValue) {
		if err = vp.bindFlagValue(flag.Name(), flag); err != nil {
			return
		}
	})
	return nil
}

func (vp *viper) bindFlagValue(key string, flag IFlagValue) error {
	if flag == nil {
		return fmt.Errorf("flag for %q is nil", key)
	}
	vp.pflags[strings.ToLower(key)] = flag
	return nil
}

func (vp *viper) Get(key string) interface{} {
	lcaseKey := strings.ToLower(key)
	if flag, exists := vp.pflags[lcaseKey]; exists {
		switch flag.ValueType() {
		case "int", "int8", "int16", "int32", "int64":
			return cast.ToInt(flag.ValueString())
		case "bool":
			return cast.ToBool(flag.ValueString())
		default:
			return flag.ValueString()
		}
	}

	return nil
}

func newViper() *viper {
	vp := new(viper)
	vp.pflags = make(map[string]IFlagValue)
	return vp
}

// stubs for PFlag Values
type stringValueTest string

func newStringValueTest(val string, p *string) *stringValueTest {
	*p = val
	return (*stringValueTest)(p)
}

func (s *stringValueTest) Set(val string) error {
	*s = stringValueTest(val)
	return nil
}

func (s *stringValueTest) Type() string {
	return "string"
}

func (s *stringValueTest) String() string {
	return string(*s)
}

func TestBindFlagValueSet(t *testing.T) {
	// 发生错误后继续解析，CommandLine就是使用这个选项
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

	testValues := map[string]*string{
		"host":     nil,
		"port":     nil,
		"endpoint": nil,
	}

	mutatedTestValues := map[string]string{
		"host":     "localhost",
		"port":     "6060",
		"endpoint": "/public",
	}

	// 添加命令行参数默认值
	for name := range testValues {
		testValues[name] = flagSet.String(name, "", "test")
	}

	// 将命令行参数添加到指定对象的 pflags 参数
	flagValueSet := PFlagValueSet{flagSet}
	vp := newViper()
	err := vp.bindFlagValues(flagValueSet)
	if err != nil {
		t.Fatalf("error binding flag set, %v", err)
	}

	// 修改命令行参数的默认值
	flagSet.VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(mutatedTestValues[flag.Name])
		flag.Changed = true
	})

	for name, expected := range mutatedTestValues {
		assert.Equal(t, vp.Get(name), expected)
	}
}

func TestBindFlagValue(t *testing.T) {
	testString := "testing"
	testValue := newStringValueTest(testString, &testString)

	flag := &pflag.Flag{
		Name:    "testflag",
		Value:   testValue,
		Changed: false,
	}

	flagValue := PFlagValue{flag}
	vp := newViper()
	vp.bindFlagValue("testvalue", flagValue)

	assert.Equal(t, testString, vp.Get("testvalue"))

	flag.Value.Set("testing_mutate")
	flag.Changed = true // hack for pflag usage

	assert.Equal(t, "testing_mutate", vp.Get("testvalue"))
}
