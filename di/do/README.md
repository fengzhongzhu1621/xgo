# 简介
do 是 Go 语言中一个轻量级的依赖注入（Dependency Injection, DI）容器，由 samber 开发。它基于 Go 1.18+ 泛型实现，为 Go 提供了一个类型安全的 DI 方案。

do 库的设计理念是简化服务组件之间的依赖管理，取代手工创建依赖关系的繁琐工作，使不同组件之间松散耦合、更易测试与维护。

与反射型 DI 框架不同，do 在注册和解析依赖时不使用反射，因此性能开销很小。

* 服务注册：使用 do.Provide 系列函数将服务构造函数注册到容器中（默认懒加载，即按需单例创建）；也可以使用 ProvideTransient 注册每次调用都新建实例的工厂（瞬时模式）；或使用 ProvideValue/ProvideNamedValue 注册已经创建好的实例（急加载）。注册时可指定名称或匿名服务（推荐匿名，由框架自动命名）。
* 依赖解析：通过 do.Invoke[T](injector) 或 do.MustInvoke[T](injector) 获取指定类型的服务实例。容器会自动根据函数签名的参数解析依赖，并以依赖图的方式按顺序实例化各服务（默认单例）。服务加载顺序为调用顺序（先调用的服务会优先初始化）。
* 生命周期管理：do 支持生命周期钩子。服务只要实现特定接口，就会被框架在适当时机调用。比如实现 do.Healthcheckable 接口的服务可以通过 do.HealthCheck[T](injector) 或 injector.HealthCheck() 进行健康检查；实现 do.Shutdownable 接口的服务会在容器关闭时被回调，以便释放资源。容器关闭时会按照服务注册的 反初始化顺序（后注册的先关闭）依次调用这些 Shutdown 方法。
* 其它特性：支持服务命名、覆盖（Override*）和组合（do.Package），可以复制容器（injector.Clone()），并提供工具函数列出已注册或已实例化的服务列表。整个库非常轻量，无外部依赖，也无需生成代码。

# Shutdown
* 依赖注入的注册时，需要在Shutdown方法中正确清理资源
* 使用依赖注入时，如果需要同类型的多个实例时使用命名依赖
