package local

import (
	"context"
	"encoding/base64"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// SessionInfo session information for mongo distributed transactions
type SessionInfo struct {
	TxnNubmer int64
	SessionID string
}

// ReloadSession is used to reset a created session's session id, so that we can
// put all the business operation
func ReloadSession(sess mongo.Session, info *SessionInfo) error {
	xsess, ok := sess.(mongo.XSession)
	if !ok {
		return errors.New("the session is not type XSession")
	}
	clientSession := xsess.ClientSession()

	sessionIDBytes, err := base64.StdEncoding.DecodeString(info.SessionID)
	if err != nil {
		return err
	}
	idx, idDoc := bsoncore.AppendDocumentStart(nil)
	idDoc = bsoncore.AppendBinaryElement(idDoc, "id", session.UUIDSubtype, sessionIDBytes[:])
	idDoc, _ = bsoncore.AppendDocumentEnd(idDoc, idx)

	clientSession.Server.SessionID = idDoc
	clientSession.SessionID = idDoc
	// i.didCommitAfterStart=false
	if info.TxnNubmer > 1 {
		// when the txnNumber is large than 1, it means that it's not the first transaction in
		// this session, we do not need to create a new transaction with this txnNumber and mongodb does
		// not allow this, so we need to change the session status from Starting to InProgressing.
		// set state to InProgressing in a same session id, then we can use the same
		// transaction number as a transaction in a single transaction session.
		// otherwise a error like this will be occured as follows:
		// (NoSuchTransaction) Given transaction number 2 does not match any in-progress transactions.
		// The active transaction number is 1
		clientSession.TransactionState = session.InProgress
	}
	return nil
}

// PrepareCommitOrAbort set state to InProgress, so that we can commit with other
// operation directly. otherwise mongodriver will do a false commit
func PrepareCommitOrAbort(sess mongo.Session) error {
	xsess, ok := sess.(mongo.XSession)
	if !ok {
		return errors.New("the session is not type XSession")
	}
	clientSession := xsess.ClientSession()

	clientSession.TransactionState = session.InProgress

	return nil
}

// ContextWithSession set the session into context if context includes session info
func ContextWithSession(ctx context.Context, sess mongo.Session) mongo.SessionContext {
	return mongo.NewSessionContext(ctx, sess)
}
