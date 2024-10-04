/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云-主机登录管理(BlueKing-SAM) available.
 * Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/pkg/errors"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// panicLog 记录 panic 信息并打印堆栈跟踪
func panicLog(rval interface{}) {
	// 将当前的堆栈跟踪打印到标准错误输出（stderr）
	debug.PrintStack()
	// 转换为字符串形式。这样可以确保无论 rval 是什么类型，都能以字符串形式记录其内容。
	rvalStr := fmt.Sprint(rval)
	// 创建一个新的错误对象
	err := errors.New(rvalStr)
	log.WithError(err).Error(fmt.Sprintf("system error %s", debug.Stack()))
}

// isBrokenPipeError 判断错误是否是网络连接已经断开
func isBrokenPipeError(err interface{}) bool {
	// 尝试将传入的错误对象转换为*net.OpError类型。如果转换成功，说明错误与网络操作有关
	if netErr, ok := err.(*net.OpError); ok {
		// 尝试将net.OpError的Err字段转换为*os.SyscallError类型。如果转换成功，说明错误与操作系统底层调用有关。
		if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
			// 用于检查传入的错误是否是“broken pipe”错误或“connection reset by peer”错误。
			// 这两种错误通常表示网络连接已经断开，无法再进行通信。
			if strings.Contains(strings.ToLower(sysErr.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(sysErr.Error()), "connection reset by peer") {
				return true
			}
		}
	}
	return false
}

const sentryValuesKey = "sentry"

// Recovery ...
func Recovery(withSentry bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			// 确保在函数退出时执行panic恢复逻辑
			if err := recover(); err != nil {
				// 检查错误是否是由于网络连接断开引起的
				brokenPipe := isBrokenPipeError(err)

				// 记录 panic 信息和堆栈跟踪
				panicLog(err)

				if withSentry && !brokenPipe {
					// 克隆当前的Sentry Hub
					hub := sentry.CurrentHub().Clone()
					// 设置请求上下文
					hub.Scope().SetRequest(c.Request)
					c.Set(sentryValuesKey, hub)

					// 使用RecoverWithContext在特定的请求上下文中记录panic
					hub.RecoverWithContext(
						context.WithValue(c.Request.Context(), sentry.RequestContextKey, c.Request),
						err,
					)
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		c.Next()
	}
}
