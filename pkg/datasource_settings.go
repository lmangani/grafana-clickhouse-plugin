package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

func NewDatasourceSettings(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var dsSettings = DatasourceSettings{}
	err := parseJson(setting.JSONData, &dsSettings); if err != nil {return nil, err}

	return &dsSettings, nil
}

type SecureSettings struct {
	Password		string			`json:"password"`
}

type DatasourceSettings struct {
	Host			string 			`json:"host"`
	Port 			int64 			`json:"port"`
	Username 		string 			`json:"username"`
	Secure			SecureSettings	`json:"secureJsonData"`
}

func (s *DatasourceSettings) Dispose() {}
