# github.com/stretchr/testify/assert
##  常用断言函数
|函数|作用|示例|
|:----|:----|:----|
|Equal|相等|assert.Equal(t, 1, 1)|
|NotEqual|不相等|assert.Equal(t, "hello", "world")|
|Nil|为 nil|assert.Nil(t, err)|
|NotNil|不为 nil|assert.NotNil(t, err)|
|Empty|为空|assert.Empty(t, obj)|
|NotEmpty|不为空|assert.NotEmpty(t, obj)|
|NoError|没有错误|assert.NoError(t, err)|
|Error|是预期的错误|assert.Error(t, err)|
|Zero|为零值|assert.Zero(t, obj)|
|True|为 true|assert.True(t, myBool)|
|False|为 false|assert.False(t, myBool)|
|Len|为预期的长度，一般用于 string slice map 等|assert.Len(t, mySlice, 3)|
|Contains|包含，一般用于 string slice map 等|assert.Contains(t, ["Hello", "World"], "World")|
|NotContains|不包含|assert.NotContains(t, "Hello", "World")|
|Subset|是子集|assert.Subset(t, [1, 2, 3], [1, 2]|
|NotSubset|不是子集|assert.NotSubset(t, [1, 3, 4], [1, 2]|
|FileExists|文件存在|assert.FileExists(t, filePath)|
|DirExists|目录存在|assert.DirExists(t, dirPath)|


## github.com/stretchr/testify/require

require在不符合断言条件时，会中断当前运行


## github.com/stretchr/testify/suite

用于组织测试用例，将一些公共逻辑和操作统一调用

