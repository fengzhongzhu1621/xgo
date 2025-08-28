package upload

// FileType 文件类型
type FileType string

const (
	// FileTypeXlsx 文件格式xlsx
	FileTypeXlsx FileType = "xlsx"
	// FileTypeXls 文件格式xls
	FileTypeXls FileType = "xls"
	// FileTypeZip 文件格式zip
	FileTypeZip FileType = "zip"
	// FileTypeYaml 文件格式yaml
	FileTypeYaml FileType = "yaml"
)
