package xgin

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
// 模块方法
func Default() *Engine {
	debugPrintWARNINGDefault()
	// 创建一个engine对象
	engine := New()
	// 加载中间件
	engine.Use(Logger(), Recovery())
	return engine
}

