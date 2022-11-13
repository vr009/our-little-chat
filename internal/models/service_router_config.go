package models

import (
	"path/filepath"
)

type ServiceRouterConfig struct {
	BaseUrl string
	Router  map[string]string
}

func (cfg ServiceRouterConfig) GetPath(method string) string {
	return filepath.Join(cfg.BaseUrl, cfg.Router[method])
}
