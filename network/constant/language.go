package constant

const (
	HTTPCookieLanguage = "blueking_language"
	HTTPHeadLanguage   = "X-BkApi-Blueking-Language"
)

type LanguageType string

const (
	Chinese LanguageType = "zh-cn"
	English LanguageType = "en"
)
