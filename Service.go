package airtable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

const (
	apiName         string = "Airtable"
	apiURL          string = "https://api.airtable.com/v0"
	DateTimeLayout  string = "2006-01-02T15:04:05.000Z"
	defaultPageSize int64  = 100
)

// type
//
type Service struct {
	apiKey      string
	httpService *go_http.Service
}

// Response represents highest level of exactonline api response
//
type Response struct {
	Data     *json.RawMessage `json:"data"`
	NextPage *NextPage        `json:"next_page"`
}

// NextPage contains info for batched data retrieval
//
type NextPage struct {
	Offset string `json:"offset"`
	Path   string `json:"path"`
	URI    string `json:"uri"`
}

type ServiceConfig struct {
	APIKey string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.APIKey == "" {
		return nil, errortools.ErrorMessage("Service APIKey not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiKey:      serviceConfig.APIKey,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", service.apiKey))
	(*requestConfig).NonDefaultHeaders = &header

	errorResponse := ErrorResponse{}
	if utilities.IsNil(requestConfig.ErrorModel) {
		// add error model
		(*requestConfig).ErrorModel = &errorResponse
	}

	request, response, e := service.httpService.HTTPRequest(requestConfig)
	if len(errorResponse.Errors) > 0 {
		messages := []string{}
		for _, message := range errorResponse.Errors {
			messages = append(messages, message.Message)
		}
		e.SetMessage(strings.Join(messages, "\n"))
	}

	return request, response, e
}

func (service *Service) url(baseID string, tableName string, params string) string {
	return fmt.Sprintf("%s/%s/%s?%s", apiURL, baseID, tableName, params)
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.apiKey
}

func (service *Service) APICallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}
