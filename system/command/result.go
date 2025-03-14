package command

// NewCmdResult returns a Cmd initialised with val and err for testing.
func NewCmdResult(val interface{}, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewSliceResult returns a SliceCmd initialised with val and err for testing.
func NewSliceResult(val []interface{}, err error) *SliceCmd {
	var cmd SliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}
