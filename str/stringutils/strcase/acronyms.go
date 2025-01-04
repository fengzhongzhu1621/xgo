package strcase

var uppercaseAcronym = map[string]string{
	"ID": "id",
}

// ConfigureAcronym 设置缩写转换规则 allows you to add additional words which will be considered acronyms.
func ConfigureAcronym(key, val string) {
	uppercaseAcronym[key] = val
}
