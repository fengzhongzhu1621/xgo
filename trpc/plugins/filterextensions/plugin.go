package filterextensions

import (
	"fmt"

	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	// MethodFilters plugin 注册的 filter 的名字
	MethodFilters = "method_filters"
)

type (
	methodClientFilters map[string][]filter.ClientFilter
	methodServerFilters map[string][]filter.ServerFilter
)

type (
	serviceMethodClientFilters map[string]map[string][]filter.ClientFilter
	serviceMethodServerFilters map[string]map[string][]filter.ServerFilter
)

type serviceMethodFiltersPlugin struct {
	client serviceMethodClientFilters
	server serviceMethodServerFilters
}

// Type 返回插件类型。
func (p *serviceMethodFiltersPlugin) Type() string {
	return PluginType
}

// Setup 初始化插件。
func (p *serviceMethodFiltersPlugin) Setup(_ string, dec plugin.Decoder) error {
	// 获取插件配置
	var cfg cfg
	if err := dec.Decode(&cfg); err != nil {
		return err
	}

	// 加载客户端过滤器
	// 从 tRPC 全局注册的 filter 中寻找，没找到则报错。
	client, err := loadClientFilters(cfg.Client, filter.GetClient)
	if err != nil {
		return fmt.Errorf("failed to load client service method filters, err: %w", err)
	}
	p.client = client

	// 加载服务器端过滤器
	server, err := loadServerFilters(cfg.Server, filter.GetServer)
	if err != nil {
		return fmt.Errorf("failed to load server service method filters, err: %w", err)
	}
	p.server = server

	// 注册过滤器
	filter.Register(MethodFilters, newServerIntercept(p.server), newClientIntercept(p.client))
	return nil
}

func loadClientFilters(
	services []cfgService,
	filterLoader func(name string) filter.ClientFilter,
) (serviceMethodClientFilters, error) {
	// 接收一组过滤器名称，通过 filterLoader 加载每个过滤器
	loadMethodFilters := func(filterNames []string) ([]filter.ClientFilter, error) {
		filters := make([]filter.ClientFilter, 0, len(filterNames))
		// 遍历方法过滤器的名称
		for _, filterName := range filterNames {
			// 获得已经加载的过滤器
			f := filterLoader(filterName)
			if f == nil {
				return nil, fmt.Errorf("filter %s not registered", filterName)
			}
			filters = append(filters, f)
		}
		return filters, nil
	}

	// 用于存储所有服务的客户端过滤器
	smf := make(serviceMethodClientFilters, len(services))
	// 遍历所有服务
	for _, service := range services {
		// 用于存储该服务的所有方法的过滤器
		mf := make(methodClientFilters, len(service.Methods))
		// 遍历服务的每个方法
		for _, method := range service.Methods {
			// 加载该方法指定的多个过滤器
			f, err := loadMethodFilters(method.Filters)
			if err != nil {
				return nil, err
			}
			mf[method.Name] = f
		}
		smf[service.Name] = mf
	}
	return smf, nil
}

func loadServerFilters(
	services []cfgService,
	filterLoader func(name string) filter.ServerFilter,
) (serviceMethodServerFilters, error) {
	loadMethodFilters := func(filterNames []string) ([]filter.ServerFilter, error) {
		filters := make([]filter.ServerFilter, 0, len(filterNames))
		for _, filterName := range filterNames {
			f := filterLoader(filterName)
			if f == nil {
				return nil, fmt.Errorf("filter %s not registered", filterName)
			}
			filters = append(filters, f)
		}
		return filters, nil
	}

	smf := make(serviceMethodServerFilters, len(services))
	for _, service := range services {
		mf := make(methodServerFilters, len(service.Methods))
		for _, method := range service.Methods {
			f, err := loadMethodFilters(method.Filters)
			if err != nil {
				return nil, err
			}
			mf[method.Name] = f
		}
		smf[service.Name] = mf
	}
	return smf, nil
}
