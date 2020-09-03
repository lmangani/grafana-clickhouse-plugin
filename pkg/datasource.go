package main

import (
  "github.com/grafana/grafana-plugin-sdk-go/backend"
  "github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
  "github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
  "github.com/grafana/grafana-plugin-sdk-go/data"
  "golang.org/x/net/context"
)

func GetDatasourceServeOpts() datasource.ServeOpts {
	ds := &ClickHouseDatasource{
		im: datasource.NewInstanceManager(NewDatasourceSettings),
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}


type ClickHouseDatasource struct {
	im instancemgmt.InstanceManager
}

func (ds *ClickHouseDatasource) getClient(ctx backend.PluginContext) (*ClickHouseClient, error)  {
	im, err := ds.im.Get(ctx); if err != nil {
		return nil, err
	}

	return &ClickHouseClient{
		settings: im.(*DatasourceSettings),
	}, nil
}

func (ds *ClickHouseDatasource) query(ctx backend.PluginContext, query *Query) backend.DataResponse  {

	onErr := func(err error) backend.DataResponse {
		backend.Logger.Error("Datasource query error: " + err.Error())
		return backend.DataResponse{Error: err}
	}

	client, err := ds.getClient(ctx); if err != nil { return onErr(err) }

	res, err := client.Query(query.Format()); if err != nil { return onErr(err) }

	frames := res.ToFrames(query.RefId, query.SplitTs)

	response := backend.DataResponse{
		Frames: make([]*data.Frame, len(frames)),
	}

	for i, frame := range frames {
	  if frame != nil {
		response.Frames[i] = frame.ToDataFrame()
	  }
	}

	return response
}

func (ds *ClickHouseDatasource) QueryData(
	ctx context.Context,
	req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {

	onErr := func(err error) (*backend.QueryDataResponse, error) {
		backend.Logger.Error("QueryData error: " + err.Error())
		return nil, err
	}

	response := backend.NewQueryDataResponse()

	for _, query := range req.Queries {
		var q = Query{}
		err := parseJson(query.JSON, &q); if err != nil { return onErr(err) }

		response.Responses[q.RefId] = ds.query(req.PluginContext, &q)
	}

	return response, nil
}

func (ds *ClickHouseDatasource) CheckHealth(
	ctx context.Context,
	req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	onErr := func(err error) (*backend.CheckHealthResult, error) {
		backend.Logger.Error("HealthCheck error: " + err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, err
	}

	client, err := ds.getClient(req.PluginContext); if err != nil {
		return onErr(err)
	}
	_, err = client.Query(DefaultQuery); if err != nil {
		return onErr(err)
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "OK",
	}, nil
}
