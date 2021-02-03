package client

import (
	"saaa-api/internal/core/config"

	"github.com/imroc/req"
	"github.com/sirupsen/logrus"
)

// Service client service interface
type Service interface {
	GetRequest(url string, header interface{}, param interface{}, v interface{}) error
	PostRequest(url string, header interface{}, param interface{}, body interface{}, v interface{}) error
}

type service struct {
	config *config.Configs
	result *config.ReturnResult
}

// NewService service client
func NewService(config *config.Configs, result *config.ReturnResult) Service {
	return &service{
		config: config,
		result: result,
	}
}

// GetRequest get request
func (s *service) GetRequest(url string, header interface{}, param interface{}, v interface{}) error {
	response, err := req.Get(url, header, param)
	if err != nil {
		logrus.Errorf("[getRequest] request get error: %s", err)
		return err
	}

	if response.Response().StatusCode >= 200 && response.Response().StatusCode < 300 {
		err = response.ToJSON(v)
		if err != nil {
			logrus.Errorf("[getRequest] convert json response from body to struct error: %s", err)
			return err
		}
		return nil
	}

	return s.result.InvalidGoogleToken
}

// PostRequest post request
func (s *service) PostRequest(url string, header interface{}, param interface{}, body interface{}, v interface{}) error {
	response, err := req.Post(url, header, param, req.BodyJSON(body))
	if err != nil {
		logrus.Errorf("[PostRequest] request post error: %s", err)
		return err
	}

	if response.Response().StatusCode >= 200 && response.Response().StatusCode < 300 {
		if v != nil {
			err = response.ToJSON(v)
			if err != nil {
				logrus.Errorf("[PostRequest] convert json response from body to struct error: %s", err)
				return err
			}
		}

		return nil
	}

	return s.result.Internal.BadRequest
}
