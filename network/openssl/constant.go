package openssl

import "regexp"

var OpensslVersionRegex = regexp.MustCompile(`(?i)OpenSSH_([0-9]+\.[0-9]+)`)
var OpensslVersionRegexV2 = regexp.MustCompile(`(?i)OpenSSH_([0-9]+)\.([0-9]+)`)
