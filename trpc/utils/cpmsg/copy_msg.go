package cpmsg

import (
	"fmt"

	copyutils "github.com/fengzhongzhu1621/xgo/copier"
	"trpc.group/trpc-go/trpc-go/codec"
)

// CopyMsg is not a common Message copy util. It's specially customized for slime.
func CopyMsg(dst, src codec.Msg) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("CopyMsg paniced, this usually means slime may not support your protocol: %v", r)
		}
	}()

	// 原始头信息保存
	oriReqHead := dst.ClientReqHead()
	oriRspHead := dst.ClientRspHead()

	// 基础消息拷贝
	// copy src Msg to dst.
	codec.CopyMsg(dst, src)

	if oriReqHead != nil {
		// 当目标消息已有请求头时，使用浅拷贝（ShallowCopy）将源请求头数据回写到原始指针
		// This must be copying back to user.
		// There will be no data access to src ReqHead. A shallow copy is just ok.
		dst.WithClientReqHead(oriReqHead)
		if err := copyutils.ShallowCopy(oriReqHead, src.ClientReqHead()); err != nil {
			return fmt.Errorf("failed to shallow copy back to original ClientReqHead, err: %w", err)
		}
	} else if src.ClientReqHead() != nil {
		// 当目标消息无请求头但源消息有时，使用深拷贝（DeepCopy）避免数据竞争
		// This must be copying to a new created msg. DeepCopy should be used to avoid data races.
		head, err := copyutils.DeepCopy(src.ClientReqHead())
		if err != nil {
			return fmt.Errorf("failed to deepcopy ClientReqHead, err: %w", err)
		}
		dst.WithClientReqHead(head)
	}

	if oriRspHead != nil {
		// 浅拷贝回写
		dst.WithClientRspHead(oriRspHead)
		if err := copyutils.ShallowCopy(oriRspHead, src.ClientRspHead()); err != nil {
			return fmt.Errorf("failed to shallow copy back to original ClientRspHead, err: %w", err)
		}
	} else if src.ClientRspHead() != nil {
		// 深拷贝新建
		head, err := copyutils.DeepCopy(src.ClientRspHead())
		if err != nil {
			return fmt.Errorf("failed to deepcopy ClientRspHead, err: %w", err)
		}
		dst.WithClientRspHead(head)
	}

	return nil
}
