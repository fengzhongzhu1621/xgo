package output

// output name, default support console and file.
const (
	OutputConsole = "console"
	OutputFile    = "file"
)

type WriteType string

const (
	WriterTypeOs   WriteType = "os"
	WriterTypeFile WriteType = "file"
)
