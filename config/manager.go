package config

// https://mp.weixin.qq.com/s/1L6ZirlfO4YnN3OOv2wQSw

import (
	"context"
	"log"
	"sync"
	"time"
)

// Repository 配置文件存储库接口
type IRepository interface {
	// UpdateConfig 更新配置
	UpdateConfig(ctx context.Context, value any) error
	// 获取配置内容
	GetConfig(ctx context.Context, obj, defaultValue any) error
}

type ConfigManager[T any] struct {
	repo         IRepository   // 存储库
	scanInterval time.Duration // 定时扫描周期
	mx           sync.Mutex    // 防止配置使用时竞争的互斥锁
	dynCfg       T             // 动态应用配置
}

// NewConfigManager 创建配置管理器
func NewConfigManager[T any](
	repo IRepository,
	initCfg T,
	scanInterval time.Duration,
) *ConfigManager[T] {
	return &ConfigManager[T]{
		dynCfg:       initCfg,      // 动态应用配置
		repo:         repo,         // 存储库
		scanInterval: scanInterval, // 定时扫描周期
	}
}

// LoadConfig 从存储库加载最新的配置到缓存
func (m *ConfigManager[T]) LoadConfig(ctx context.Context) error {
	var cfg T
	if err := m.repo.GetConfig(ctx, &cfg, m.dynCfg); err != nil {
		return err
	}
	m.mx.Lock()
	defer m.mx.Unlock()
	m.dynCfg = cfg
	return nil
}

// GetConfig 从配置管理器获取缓存的动态配置
func (m *ConfigManager[T]) GetConfig() T {
	m.mx.Lock()
	defer m.mx.Unlock()
	return m.dynCfg
}

// Run 定时刷新缓存的动态配置
func (m *ConfigManager[T]) Run(ctx context.Context, wg *sync.WaitGroup) error {
	// 从存储库加载最新的配置到缓存
	if err := m.LoadConfig(ctx); err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(m.scanInterval):
				// 这个通道会在指定的时间间隔（m.scanInterval）后发送一个消息。
				// 当这个通道接收到消息时，表示已经到达了加载配置文件的时间点。
				//
				// 需要注意的是，每次循环都会创建一个新的time.After(m.scanInterval)通道，
				// 这是因为time.After函数返回的通道在发送完值之后就会被关闭，不能再次使用。
				// 所以，在每次循环中都需要重新创建一个新的通道。
				if err := m.LoadConfig(ctx); err != nil {
					log.Printf("error while loading the config: %v", err)
				}
			}
		}
	}()
	return nil
}
