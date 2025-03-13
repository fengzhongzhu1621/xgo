package validator

import "regexp"

const (
	userPattern = `^(\d|[a-zA-Z])([a-zA-Z0-9\@.,_-])*$`
)

var (
	userRegexp = regexp.MustCompile(userPattern)
)

func IsUser(sInput string) bool {
	return userRegexp.MatchString(sInput)
}
