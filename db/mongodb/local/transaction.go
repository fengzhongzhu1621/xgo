package local

import (
	"context"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/db/db"
	"github.com/fengzhongzhu1621/xgo/db/redis"
	"github.com/fengzhongzhu1621/xgo/network/constant"
	log "github.com/sirupsen/logrus"
)

// CommitTransaction 提交事务
func (c *Mongo) CommitTransaction(ctx context.Context, cap *db.TxnCapable) error {
	rid := ctx.Value(constant.ContextRequestIDField)

	// check if txn number exists, if not, then no db operation with transaction is executed, committing will return an
	// error: "(NoSuchTransaction) Given transaction number 1 does not match any in-progress transactions. The active
	// transaction number is -1.". So we will return directly in this situation.
	txnNumber, err := c.tm.GetTxnNumber(cap.SessionID)
	if err != nil {
		if redis.IsNilErr(err) {
			log.Infof("commit transaction: %s but no transaction need to commit, *skip*, rid: %s", cap.SessionID, rid)
			return nil
		}
		return fmt.Errorf("get txn number failed, err: %v", err)
	}
	if txnNumber == 0 {
		log.Infof("commit transaction: %s but no transaction to commit, **skip**, rid: %s", cap.SessionID, rid)
		return nil
	}

	reloadSession, err := c.tm.PrepareTransaction(cap, c.dbc)
	if err != nil {
		log.Errorf("commit transaction, but prepare transaction failed, err: %v, rid: %v", err, rid)
		return err
	}
	// reset the transaction state, so that we can commit the transaction after start the
	// transaction immediately.
	if err := PrepareCommitOrAbort(reloadSession); err != nil {
		log.Errorf("reset the commit transaction state failed, err: %v, rid: %v", err, rid)
		return err
	}

	// we commit the transaction with the session id
	err = reloadSession.CommitTransaction(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %s failed, err: %v, rid: %v", cap.SessionID, err, rid)
	}

	err = c.tm.RemoveSessionKey(cap.SessionID)
	if err != nil {
		// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
		log.Errorf("commit transaction, but delete txn session: %s key failed, err: %v, rid: %v", cap.SessionID, err, rid)
		// do not return.
	}

	return nil
}

// AbortTransaction 取消事务
func (c *Mongo) AbortTransaction(ctx context.Context, cap *db.TxnCapable) (bool, error) {
	rid := ctx.Value(constant.ContextRequestIDField)
	reloadSession, err := c.tm.PrepareTransaction(cap, c.dbc)
	if err != nil {
		log.Errorf("abort transaction, but prepare transaction failed, err: %v, rid: %v", err, rid)
		return false, err
	}
	// reset the transaction state, so that we can abort the transaction after start the
	// transaction immediately.
	if err := PrepareCommitOrAbort(reloadSession); err != nil {
		log.Errorf("reset abort transaction state failed, err: %v, rid: %v", err, rid)
		return false, err
	}

	// we abort the transaction with the session id
	err = reloadSession.AbortTransaction(ctx)
	if err != nil {
		return false, fmt.Errorf("abort transaction: %s failed, err: %v, rid: %v", cap.SessionID, err, rid)
	}

	err = c.tm.RemoveSessionKey(cap.SessionID)
	if err != nil {
		// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
		log.Errorf("abort transaction, but delete txn session: %s key failed, err: %v, rid: %v", cap.SessionID, err, rid)
		// do not return.
	}

	errorType := c.tm.GetTxnError(sessionKey(cap.SessionID))
	switch errorType {
	// retry when the transaction error type is write conflict, which means the transaction conflicts with another one
	case WriteConflictType:
		return true, nil
	}

	return false, nil
}
