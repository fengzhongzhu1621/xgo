package serializer

type GetVersionContentSerializer struct {
	Version string `json:"version" form:"version" binding:"required"`
}
