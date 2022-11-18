package models

type ServiceRouterConfig struct {
	BaseUrl string
	Router  map[string]string
}

func (cfg ServiceRouterConfig) GetPath(method string) string {
	return cfg.BaseUrl + cfg.Router[method]
}
