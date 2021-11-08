package airtable

import (
	"fmt"
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
	ID          string                 `json:"id"`
	Fields      map[string]string      `json:"fields"`
	CreatedTime a_types.DateTimeString `json:"createdTime"`
}

type GetRecordsConfig struct {
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
func (service *Service) GetRecords(config *GetRecordsConfig) ([]Record, *errortools.Error) {
	records := []Record{}

	pageSize := defaultPageSize

	if config.PageSize != nil {
		pageSize = *config.PageSize
	}

	params := url.Values{}
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
	params.Set("pageSize", fmt.Sprintf("%v", pageSize))
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

	for {
		_records := Records{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("?%s", params.Encode())),
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

	return records, nil
}
