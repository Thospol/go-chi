package guest

import "saaa-api/internal/core/config"

// Service service guest interface
type Service interface{}

type service struct {
	config *config.Configs
	result *config.ReturnResult
}

// NewService new service guest
func NewService(config *config.Configs, result *config.ReturnResult) Service {
	return &service{
		config: config,
		result: result,
	}
}
