package fileutils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ms "xgo/encoding/mapstructure"
	"xgo/utils/testutil"
)

var yamlExampleWithExtras = []byte(`Existing: true
Bogus: true
`)

type testUnmarshalExtra struct {
	Existing bool
}

var tomlExample = []byte(`
title = "TOML Example"

[owner]
organization = "MongoDB"
Bio = "MongoDB Chief Developer Advocate & Hacker at Large"
dob = 1979-05-27T07:32:00Z # First class dates? Why not?`)

var dotenvExample = []byte(`
TITLE_DOTENV="DotEnv Example"
TYPE_DOTENV=donut
NAME_DOTENV=Cake`)

var jsonExample = []byte(`{
"id": "0001",
"type": "donut",
"name": "Cake",
"ppu": 0.55,
"batters": {
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)

var remoteExample = []byte(`{
"id":"0002",
"type":"cronut",
"newkey":"remote"
}`)

var iniExample = []byte(`; Package name
NAME        = ini
; Package version
VERSION     = v1
; Package import path
IMPORT_PATH = gopkg.in/%(NAME)s.%(VERSION)s

# Information about package author
# Bio can be written in multiple lines.
[author]
NAME   = Unknown  ; Succeeding comment
E-MAIL = fake@localhost
GITHUB = https://github.com/%(NAME)s
BIO    = """Gopher.
Coding addict.
Good man.
"""  # Succeeding comment`)

func initConfig(typ, config string) *Watcher {
	v := Reset()
	v.SetConfigType(typ)
	r := strings.NewReader(config)

	if err := v.unmarshalReader(r, v.config); err != nil {
		panic(err)
	}

	return v
}

func initJSON() *Watcher {
	v := Reset()
	v.SetConfigType("json")
	r := bytes.NewReader(jsonExample)

	v.unmarshalReader(r, v.config)

	return v
}

func initDotEnv() *Watcher {
	v := Reset()
	v.SetConfigType("env")
	r := bytes.NewReader(dotenvExample)

	v.unmarshalReader(r, v.config)

	return v
}

func initIni() *Watcher {
	v := Reset()
	v.SetConfigType("ini")
	r := bytes.NewReader(iniExample)

	v.unmarshalReader(r, v.config)

	return v
}

// make directories for testing
func initDirs(t *testing.T) (string, string, func()) {
	var (
		testDirs = []string{`a a`, `b`, `C_`}
		config   = `improbable`
	)

	if runtime.GOOS != "windows" {
		testDirs = append(testDirs, `d\d`)
	}

	root, err := ioutil.TempDir("", "")
	require.NoError(t, err, "Failed to create temporary directory")

	cleanup := true
	defer func() {
		if cleanup {
			os.Chdir("..")
			os.RemoveAll(root)
		}
	}()

	assert.Nil(t, err)

	err = os.Chdir(root)
	require.Nil(t, err)

	for _, dir := range testDirs {
		err = os.Mkdir(dir, 0o750)
		assert.Nil(t, err)

		err = ioutil.WriteFile(
			path.Join(dir, config+".toml"),
			[]byte("key = \"value is "+dir+"\"\n"),
			0o640)
		assert.Nil(t, err)
	}

	cleanup = false
	return root, config, func() {
		os.Chdir("..")
		os.RemoveAll(root)
	}
}

// stubs for PFlag Values
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string {
	return string(*s)
}

func TestUnmarshalExact(t *testing.T) {
	vip := New()
	target := &testUnmarshalExtra{}
	vip.SetConfigType("yaml")
	r := bytes.NewReader(yamlExampleWithExtras)
	vip.ReadConfig(r)
	err := vip.UnmarshalExact(target)
	if err == nil {
		t.Fatal("UnmarshalExact should error when populating a struct from a conf that contains unused fields")
	}
}

func TestOverrides(t *testing.T) {
	v := Reset()
	v.Set("age", 40)
	assert.Equal(t, 40, v.Get("age"))
}

func TestJSON(t *testing.T) {
	v := initJSON()
	assert.Equal(t, "0001", v.Get("id"))
}

func TestDotEnv(t *testing.T) {
	v := initDotEnv()
	assert.Equal(t, "DotEnv Example", v.Get("title_dotenv"))
}

func TestIni(t *testing.T) {
	v := initIni()
	assert.Equal(t, "ini", v.Get("default.name"))
}

func TestRemotePrecedence(t *testing.T) {
	v := initJSON()

	remote := bytes.NewReader(remoteExample)
	assert.Equal(t, "0001", v.Get("id"))
	v.unmarshalReader(remote, v.kvstore)
	assert.Equal(t, "0001", v.Get("id"))
	assert.NotEqual(t, "cronut", v.Get("type"))
	assert.Equal(t, "remote", v.Get("newkey"))
	v.Set("newkey", "newvalue")
	assert.NotEqual(t, "remote", v.Get("newkey"))
	assert.Equal(t, "newvalue", v.Get("newkey"))
	v.Set("newkey", "remote")
}

func TestAutoEnv(t *testing.T) {
	v := Reset()

	v.AutomaticEnv()

	testutil.Setenv(t, "FOO_BAR", "13")

	assert.Equal(t, "13", v.Get("foo_bar"))
}

func TestAutoEnvWithPrefix(t *testing.T) {
	v := Reset()

	v.AutomaticEnv()
	v.SetEnvPrefix("Baz")

	testutil.Setenv(t, "BAZ_BAR", "13")

	assert.Equal(t, "13", v.Get("bar"))
}

func TestUnmarshalWithDecoderOptions(t *testing.T) {
	v := Reset()
	v.Set("credentials", "{\"foo\":\"bar\"}")

	opt := ms.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		// Custom Decode Hook Function
		func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
			if rf != reflect.String || rt != reflect.Map {
				return data, nil
			}
			m := map[string]string{}
			raw := data.(string)
			if raw == "" {
				return m, nil
			}
			return m, json.Unmarshal([]byte(raw), &m)
		},
	))

	type config struct {
		Credentials map[string]string
	}

	var C config

	err := v.Unmarshal(&C, opt)
	if err != nil {
		t.Fatalf("unable to decode into struct, %v", err)
	}

	assert.Equal(t, &config{
		Credentials: map[string]string{"foo": "bar"},
	}, &C)
}

func TestSizeInBytes(t *testing.T) {
	input := map[string]uint{
		"":               0,
		"b":              0,
		"12 bytes":       0,
		"200000000000gb": 0,
		"12 b":           12,
		"43 MB":          43 * (1 << 20),
		"10mb":           10 * (1 << 20),
		"1gb":            1 << 30,
	}

	for str, expected := range input {
		assert.Equal(t, expected, ParseSizeInBytes(str), str)
	}
}

var jsonWriteExpected = []byte(`{
  "batters": {
    "batter": [
      {
        "type": "Regular"
      },
      {
        "type": "Chocolate"
      },
      {
        "type": "Blueberry"
      },
      {
        "type": "Devil's Food"
      }
    ]
  },
  "id": "0001",
  "name": "Cake",
  "ppu": 0.55,
  "type": "donut"
}`)

func TestWriteConfigTOML(t *testing.T) {
	fs := afero.NewMemMapFs()

	testCases := map[string]struct {
		configName string
		configType string
		fileName   string
		input      []byte
	}{
		"with file extension": {
			configName: "c",
			configType: "toml",
			fileName:   "c.toml",
			input:      tomlExample,
		},
		"without file extension": {
			configName: "c",
			configType: "toml",
			fileName:   "c",
			input:      tomlExample,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := New()
			v.SetFs(fs)
			v.SetConfigName(tc.configName)
			v.SetConfigType(tc.configType)
			err := v.ReadConfig(bytes.NewBuffer(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			if err := v.WriteConfigAs(tc.fileName); err != nil {
				t.Fatal(err)
			}

			// The TOML String method does not order the contents.
			// Therefore, we must read the generated file and compare the data.
			v2 := New()
			v2.SetFs(fs)
			v2.SetConfigName(tc.configName)
			v2.SetConfigType(tc.configType)
			v2.SetConfigFile(tc.fileName)
			err = v2.ReadInConfig()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, v.GetString("title"), v2.GetString("title"))
			assert.Equal(t, v.GetString("owner.bio"), v2.GetString("owner.bio"))
			assert.Equal(t, v.GetString("owner.dob"), v2.GetString("owner.dob"))
			assert.Equal(t, v.GetString("owner.organization"), v2.GetString("owner.organization"))
		})
	}
}

func TestWriteConfigDotEnv(t *testing.T) {
	fs := afero.NewMemMapFs()
	testCases := map[string]struct {
		configName string
		configType string
		fileName   string
		input      []byte
	}{
		"with file extension": {
			configName: "c",
			configType: "env",
			fileName:   "c.env",
			input:      dotenvExample,
		},
		"without file extension": {
			configName: "c",
			configType: "env",
			fileName:   "c",
			input:      dotenvExample,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := New()
			v.SetFs(fs)
			v.SetConfigName(tc.configName)
			v.SetConfigType(tc.configType)
			err := v.ReadConfig(bytes.NewBuffer(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			if err := v.WriteConfigAs(tc.fileName); err != nil {
				t.Fatal(err)
			}

			// The TOML String method does not order the contents.
			// Therefore, we must read the generated file and compare the data.
			v2 := New()
			v2.SetFs(fs)
			v2.SetConfigName(tc.configName)
			v2.SetConfigType(tc.configType)
			v2.SetConfigFile(tc.fileName)
			err = v2.ReadInConfig()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, v.GetString("title_dotenv"), v2.GetString("title_dotenv"))
			assert.Equal(t, v.GetString("type_dotenv"), v2.GetString("type_dotenv"))
			assert.Equal(t, v.GetString("kind_dotenv"), v2.GetString("kind_dotenv"))
		})
	}
}

func TestSafeWriteConfigWithMissingConfigPath(t *testing.T) {
	v := New()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigName("c")
	v.SetConfigType("yaml")
	require.EqualError(t, v.SafeWriteConfig(), "missing configuration for 'configPath'")
}

func TestSafeWriteConfigWithExistingFile(t *testing.T) {
	v := New()
	fs := afero.NewMemMapFs()
	fs.Create(testutil.AbsFilePath(t, "/test/c.yaml"))
	v.SetFs(fs)
	v.AddConfigPath("/test")
	v.SetConfigName("c")
	v.SetConfigType("yaml")
	err := v.SafeWriteConfig()
	require.Error(t, err)
	_, ok := err.(ConfigFileAlreadyExistsError)
	assert.True(t, ok, "Expected ConfigFileAlreadyExistsError")
}

func TestSafeWriteConfigAsWithExistingFile(t *testing.T) {
	v := New()
	fs := afero.NewMemMapFs()
	fs.Create("/test/c.yaml")
	v.SetFs(fs)
	err := v.SafeWriteConfigAs("/test/c.yaml")
	require.Error(t, err)
	_, ok := err.(ConfigFileAlreadyExistsError)
	assert.True(t, ok, "Expected ConfigFileAlreadyExistsError")
}

var yamlMergeExampleTgt = []byte(`
hello:
    pop: 37890
    largenum: 765432101234567
    num2pow63: 9223372036854775808
    universe: null
    world:
    - us
    - uk
    - fr
    - de
`)

var yamlMergeExampleSrc = []byte(`
hello:
    pop: 45000
    largenum: 7654321001234567
    universe:
    - mw
    - ad
    ints:
    - 1
    - 2
fu: bar
`)

var jsonMergeExampleTgt = []byte(`
{
	"hello": {
		"foo": null,
		"pop": 123456
	}
}
`)

var jsonMergeExampleSrc = []byte(`
{
	"hello": {
		"foo": "foo str",
		"pop": "pop str"
	}
}
`)

func TestMergeConfig(t *testing.T) {
	v := New()
	v.SetConfigType("yml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlMergeExampleTgt)); err != nil {
		t.Fatal(err)
	}

	if pop := v.GetInt("hello.pop"); pop != 37890 {
		t.Fatalf("pop != 37890, = %d", pop)
	}

	if pop := v.GetInt32("hello.pop"); pop != int32(37890) {
		t.Fatalf("pop != 37890, = %d", pop)
	}

	if pop := v.GetInt64("hello.largenum"); pop != int64(765432101234567) {
		t.Fatalf("int64 largenum != 765432101234567, = %d", pop)
	}

	if pop := v.GetUint("hello.pop"); pop != 37890 {
		t.Fatalf("uint pop != 37890, = %d", pop)
	}

	if pop := v.GetUint32("hello.pop"); pop != 37890 {
		t.Fatalf("uint32 pop != 37890, = %d", pop)
	}

	if pop := v.GetUint64("hello.num2pow63"); pop != 9223372036854775808 {
		t.Fatalf("uint64 num2pow63 != 9223372036854775808, = %d", pop)
	}

	if world := v.GetStringSlice("hello.world"); len(world) != 4 {
		t.Fatalf("len(world) != 4, = %d", len(world))
	}

	if fu := v.GetString("fu"); fu != "" {
		t.Fatalf("fu != \"\", = %s", fu)
	}

	if err := v.MergeConfig(bytes.NewBuffer(yamlMergeExampleSrc)); err != nil {
		t.Fatal(err)
	}

	if pop := v.GetInt("hello.pop"); pop != 45000 {
		t.Fatalf("pop != 45000, = %d", pop)
	}

	if pop := v.GetInt32("hello.pop"); pop != int32(45000) {
		t.Fatalf("pop != 45000, = %d", pop)
	}

	if pop := v.GetInt64("hello.largenum"); pop != int64(7654321001234567) {
		t.Fatalf("int64 largenum != 7654321001234567, = %d", pop)
	}

	if world := v.GetStringSlice("hello.world"); len(world) != 4 {
		t.Fatalf("len(world) != 4, = %d", len(world))
	}

	if universe := v.GetStringSlice("hello.universe"); len(universe) != 2 {
		t.Fatalf("len(universe) != 2, = %d", len(universe))
	}

	if ints := v.GetIntSlice("hello.ints"); len(ints) != 2 {
		t.Fatalf("len(ints) != 2, = %d", len(ints))
	}

	if fu := v.GetString("fu"); fu != "bar" {
		t.Fatalf("fu != \"bar\", = %s", fu)
	}
}

func TestMergeConfigOverrideType(t *testing.T) {
	v := New()
	v.SetConfigType("json")
	if err := v.ReadConfig(bytes.NewBuffer(jsonMergeExampleTgt)); err != nil {
		t.Fatal(err)
	}

	if err := v.MergeConfig(bytes.NewBuffer(jsonMergeExampleSrc)); err != nil {
		t.Fatal(err)
	}

	if pop := v.GetString("hello.pop"); pop != "pop str" {
		t.Fatalf("pop != \"pop str\", = %s", pop)
	}

	if foo := v.GetString("hello.foo"); foo != "foo str" {
		t.Fatalf("foo != \"foo str\", = %s", foo)
	}
}

func TestMergeConfigNoMerge(t *testing.T) {
	v := New()
	v.SetConfigType("yml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlMergeExampleTgt)); err != nil {
		t.Fatal(err)
	}

	if pop := v.GetInt("hello.pop"); pop != 37890 {
		t.Fatalf("pop != 37890, = %d", pop)
	}

	if world := v.GetStringSlice("hello.world"); len(world) != 4 {
		t.Fatalf("len(world) != 4, = %d", len(world))
	}

	if fu := v.GetString("fu"); fu != "" {
		t.Fatalf("fu != \"\", = %s", fu)
	}

	if err := v.ReadConfig(bytes.NewBuffer(yamlMergeExampleSrc)); err != nil {
		t.Fatal(err)
	}

	if pop := v.GetInt("hello.pop"); pop != 45000 {
		t.Fatalf("pop != 45000, = %d", pop)
	}

	if world := v.GetStringSlice("hello.world"); len(world) != 0 {
		t.Fatalf("len(world) != 0, = %d", len(world))
	}

	if universe := v.GetStringSlice("hello.universe"); len(universe) != 2 {
		t.Fatalf("len(universe) != 2, = %d", len(universe))
	}

	if ints := v.GetIntSlice("hello.ints"); len(ints) != 2 {
		t.Fatalf("len(ints) != 2, = %d", len(ints))
	}

	if fu := v.GetString("fu"); fu != "bar" {
		t.Fatalf("fu != \"bar\", = %s", fu)
	}
}

func TestMergeConfigMap(t *testing.T) {
	v := New()
	v.SetConfigType("yml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlMergeExampleTgt)); err != nil {
		t.Fatal(err)
	}

	assert := func(i int) {
		large := v.GetInt64("hello.largenum")
		pop := v.GetInt("hello.pop")
		if large != 765432101234567 {
			t.Fatal("Got large num:", large)
		}

		if pop != i {
			t.Fatal("Got pop:", pop)
		}
	}

	assert(37890)

	update := map[string]interface{}{
		"Hello": map[string]interface{}{
			"Pop": 1234,
		},
		"World": map[interface{}]interface{}{
			"Rock": 345,
		},
	}

	if err := v.MergeConfigMap(update); err != nil {
		t.Fatal(err)
	}

	if rock := v.GetInt("world.rock"); rock != 345 {
		t.Fatal("Got rock:", rock)
	}

	assert(1234)
}

func TestDotParameter(t *testing.T) {
	v := initJSON()
	// shoud take precedence over batters defined in jsonExample
	r := bytes.NewReader([]byte(`{ "batters.batter": [ { "type": "Small" } ] }`))
	v.unmarshalReader(r, v.config)

	actual := v.Get("batters.batter")
	expected := []interface{}{map[string]interface{}{"type": "Small"}}
	assert.Equal(t, expected, actual)
}

func TestParseNested(t *testing.T) {
	type duration struct {
		Delay time.Duration
	}

	type item struct {
		Name   string
		Delay  time.Duration
		Nested duration
	}

	config := `[[parent]]
	delay="100ms"
	[parent.nested]
	delay="200ms"
`
	initConfig("toml", config)

	var items []item
	err := v.UnmarshalKey("parent", &items)
	if err != nil {
		t.Fatalf("unable to decode into struct, %v", err)
	}

	assert.Equal(t, 1, len(items))
	assert.Equal(t, 100*time.Millisecond, items[0].Delay)
	assert.Equal(t, 200*time.Millisecond, items[0].Nested.Delay)
}

func newViperWithConfigFile(t *testing.T) (*Watcher, string, func()) {
	watchDir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	configFile := path.Join(watchDir, "config.yaml")
	err = ioutil.WriteFile(configFile, []byte("foo: bar\n"), 0o640)
	require.Nil(t, err)
	cleanup := func() {
		os.RemoveAll(watchDir)
	}
	v := New()
	v.SetConfigFile(configFile)
	err = v.ReadInConfig()
	require.Nil(t, err)
	require.Equal(t, "bar", v.Get("foo"))
	return v, configFile, cleanup
}

func newViperWithSymlinkedConfigFile(t *testing.T) (*Watcher, string, string, func()) {
	watchDir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	dataDir1 := path.Join(watchDir, "data1")
	err = os.Mkdir(dataDir1, 0o777)
	require.Nil(t, err)
	realConfigFile := path.Join(dataDir1, "config.yaml")
	t.Logf("Real config file location: %s\n", realConfigFile)
	err = ioutil.WriteFile(realConfigFile, []byte("foo: bar\n"), 0o640)
	require.Nil(t, err)
	cleanup := func() {
		os.RemoveAll(watchDir)
	}
	// now, symlink the tm `data1` dir to `data` in the baseDir
	os.Symlink(dataDir1, path.Join(watchDir, "data"))
	// and link the `<watchdir>/datadir1/config.yaml` to `<watchdir>/config.yaml`
	configFile := path.Join(watchDir, "config.yaml")
	os.Symlink(path.Join(watchDir, "data", "config.yaml"), configFile)
	t.Logf("Config file location: %s\n", path.Join(watchDir, "config.yaml"))
	// init Viper
	v := New()
	v.SetConfigFile(configFile)
	err = v.ReadInConfig()
	require.Nil(t, err)
	require.Equal(t, "bar", v.Get("foo"))
	return v, watchDir, configFile, cleanup
}

var yamlDeepNestedSlices = []byte(`TV:
- title: "The expanse"
  seasons:
  - first_released: "December 14, 2015"
    episodes:
    - title: "Dulcinea"
      air_date: "December 14, 2015"
    - title: "The Big Empty"
      air_date: "December 15, 2015"
    - title: "Remember the Cant"
      air_date: "December 22, 2015"
  - first_released: "February 1, 2017"
    episodes:
    - title: "Safe"
      air_date: "February 1, 2017"
    - title: "Doors & Corners"
      air_date: "February 1, 2017"
    - title: "Static"
      air_date: "February 8, 2017"
  episodes:
    - ["Dulcinea", "The Big Empty", "Remember the Cant"]
    - ["Safe", "Doors & Corners", "Static"]
`)

func TestSliceIndexAccess(t *testing.T) {
	v.SetConfigType("yaml")
	r := strings.NewReader(string(yamlDeepNestedSlices))

	err := v.unmarshalReader(r, v.config)
	require.NoError(t, err)

	assert.Equal(t, "The expanse", v.GetString("tv.0.title"))
	assert.Equal(t, "February 1, 2017", v.GetString("tv.0.seasons.1.first_released"))
	assert.Equal(t, "Static", v.GetString("tv.0.seasons.1.episodes.2.title"))
	assert.Equal(t, "December 15, 2015", v.GetString("tv.0.seasons.0.episodes.1.air_date"))

	// Test for index out of bounds
	assert.Equal(t, "", v.GetString("tv.0.seasons.2.first_released"))

	// Accessing multidimensional arrays
	assert.Equal(t, "Static", v.GetString("tv.0.episodes.1.2"))
}

func BenchmarkGetBool(b *testing.B) {
	key := "BenchmarkGetBool"
	v = New()
	v.Set(key, true)

	for i := 0; i < b.N; i++ {
		if !v.GetBool(key) {
			b.Fatal("GetBool returned false")
		}
	}
}

func BenchmarkGet(b *testing.B) {
	key := "BenchmarkGet"
	v = New()
	v.Set(key, true)

	for i := 0; i < b.N; i++ {
		if !v.Get(key).(bool) {
			b.Fatal("Get returned false")
		}
	}
}

// BenchmarkGetBoolFromMap is the "perfect result" for the above.
func BenchmarkGetBoolFromMap(b *testing.B) {
	m := make(map[string]bool)
	key := "BenchmarkGetBool"
	m[key] = true

	for i := 0; i < b.N; i++ {
		if !m[key] {
			b.Fatal("Map value was false")
		}
	}
}
