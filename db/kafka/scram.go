package kafka

import (
	"crypto/sha512"

	"github.com/xdg-go/scram"
)

var SHA512 scram.HashGeneratorFcn = sha512.New

// XDGSCRAMClient is a SCRAM client.
type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

// Begin prepares the client for the SCRAM exchange with the server with a user name and a password.
func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.NewConversation()
	return nil
}

// Step steps client through the SCRAM exchange.
func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

// Done should return true when the SCRAM conversation is over.
func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
