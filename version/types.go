package version

type VersionLog struct {
	Version   string `json:"version"`
	ReleaseAt string `json:"release_at"`
	Content   string `json:"content"`
}

type ListVersionsResponse struct {
	ShowVersion bool         `json:"show_version"`
	VersionLogs []VersionLog `json:"version_logs"`
	LastVersion string       `json:"last_version"`
}

type GetVersionContentRequestParam struct {
	Language string `json:"language" form:"language"`
	Version  string `json:"version" form:"version" binding:"required"`
	Operator string `json:"operator" form:"operator" binding:"required"`
}

type IVersionsConfigService interface {
	GetVersionContent(getVersionContentParam *GetVersionContentRequestParam) (string, error)

	ListVersions(language string) (*ListVersionsResponse, error)
}
