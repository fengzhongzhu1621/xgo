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
	"encoding/hex"

	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

// RequestID 用于处理HTTP请求并为其生成或提取一个唯一的X-Request-ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("Middleware: RequestID")

		// 尝试从HTTP请求头中获取X-Request-ID
		requestID := c.GetHeader(constant.RequestIDHeaderKey)
		// 如果该头不存在或其长度不是32（通常UUID的长度），则生成一个新的UUID，并将其转换为十六进制字符串
		if requestID == "" || len(requestID) != 32 {
			requestID = hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes())
		}

		// 设置Request ID，将生成的或提取的X-Request-ID存储在Gin的上下文中，以便后续的处理函数可以访问它
		utils.SetRequestID(c, requestID)

		// 将该ID设置回HTTP响应头中，这样客户端也可以接收到这个ID
		c.Writer.Header().Set(constant.RequestIDHeaderKey, requestID)

		c.Next()
	}
}
