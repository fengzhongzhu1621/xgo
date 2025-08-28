package ip

import (
	"fmt"
)

const (
	// PatternIP regular pattern for ip
	PatternIP = `^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(
1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$`
	// PatternMultipleIP regular pattern for Multiple ip
	PatternMultipleIP = `^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(
1\\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)(,(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.
((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d))*$`
	// PatternPort regular pattern for port range
	PatternPort = `(([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))`
)

// PatternMultiplePortRange regular pattern for multiple port range
var PatternMultiplePortRange = fmt.Sprintf(`^((%s-%s)|(%s))(,((%s)|(%s-%s)))*$`,
	PatternPort, PatternPort, PatternPort, PatternPort, PatternPort, PatternPort)
