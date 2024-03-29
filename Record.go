package airtable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	a_types "github.com/leapforce-libraries/go_airtable/types"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Records struct {
	Records []Record `json:"records"`
	Offset  string   `json:"offset"`
}

type Record struct {
	Id          string                     `json:"id"`
	Fields      map[string]json.RawMessage `json:"-"`
	FieldsRaw   json.RawMessage            `json:"fields"`
	CreatedTime *a_types.DateTimeString    `json:"createdTime,omitempty"`
}

func (r *Record) UnmarshalJSON(data []byte) error {
	type R Record
	if err := json.Unmarshal(data, (*R)(r)); err != nil {
		return err
	}

	if err := json.Unmarshal(r.FieldsRaw, &r.Fields); err != nil {
		return err
	}

	return nil
}

type UpdateRecordsConfig struct {
	BaseId    string
	TableName string
	Records   []Record
}

// UpdateRecords updates up to 10 records
//
func (service *Service) UpdateRecords(config *UpdateRecordsConfig) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodPatch,
		Url:    service.url(config.BaseId, config.TableName, ""),
		BodyModel: struct {
			Records []Record `json:"records"`
		}{config.Records},
	}
	_, _, e := service.httpRequest(&requestConfig)

	return e
}

type GetRecordsConfig struct {
	BaseId          string
	TableName       string
	Fields          *[]string
	FilterByFormula *string
	MaxRecords      *int64
	PageSize        *int64
	Sort            *[]struct {
		Field     string
		Direction string
	}
	View       *string
	CellFormat *string
	TimeZone   *string
	UserLocale *string
}

// GetRecords returns all records
//
func (service *Service) GetRecords(config *GetRecordsConfig) (*[]Record, *errortools.Error) {
	records := []Record{}

	pageSize := defaultPageSize
	params := url.Values{}

	if config != nil {
		if config.PageSize != nil {
			pageSize = *config.PageSize
		}

		if config.Fields != nil {
			for _, field := range *config.Fields {
				params.Add("fields[]", field)
			}
		}
		if config.FilterByFormula != nil {
			params.Set("filterByFormula", *config.FilterByFormula)
		}
		if config.MaxRecords != nil {
			params.Set("maxRecords", fmt.Sprintf("%v", *config.MaxRecords))
		}
		if config.Sort != nil {
			for i, sort := range *config.Sort {
				params.Add(fmt.Sprintf("sort[%v][field]", i), sort.Field)
				params.Add(fmt.Sprintf("sort[%v][direction]", i), sort.Direction)
			}
		}
		if config.View != nil {
			params.Set("view", *config.View)
		}
		if config.CellFormat != nil {
			params.Set("cellFormat", *config.CellFormat)
		}
		if config.TimeZone != nil {
			params.Set("timeZone", *config.TimeZone)
		}
		if config.UserLocale != nil {
			params.Set("userLocale", *config.UserLocale)
		}
	}

	params.Set("pageSize", fmt.Sprintf("%v", pageSize))

	for {
		_records := Records{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(config.BaseId, config.TableName, params.Encode()),
			ResponseModel: &_records,
		}
		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		records = append(records, _records.Records...)

		if _records.Offset == "" {
			break
		}

		params.Set("offset", _records.Offset)
	}

	return &records, nil
}
