/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云-主机登录管理(BlueKing-SAM) available.
 * Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package router

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	cfg := config.GetGlobalConfig()

	router := gin.Default()
	rootGroup := router.Group("/")
	RegisterRouter(cfg, router)
	RegisterRouterGroup(cfg, rootGroup)

	assert.NotNil(t, router)
	assert.NotNil(t, rootGroup)
}
