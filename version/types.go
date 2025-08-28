package version

type VersionLog struct {
	// 版本号
	Version string `json:"version"`
	// 发布时间
	ReleaseAt string `json:"release_at"`
	// 版本内容
	Content string `json:"content"`
}

// Revision
const (
	RevisionEnterprise = "enterprise"
	RevisionCommunity  = "community"
	RevisionOpenSource = "opensource"
)
