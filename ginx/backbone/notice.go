package backbone

import (
	"github.com/fengzhongzhu1621/xgo/db/zookeeper/zkclient"
)

type noticeHandler struct {
	client   *zkclient.ZkClient
	addrport string
}
