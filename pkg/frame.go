package main

import (
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type ClickHouseFrame struct {
	RefId string
	Name string
	Fields []*ClickHouseField
}

func NewFrames(query *Query, response *Response) []*ClickHouseFrame {
  if query.SplitTS {
	return splitFrame(query.RefId, response.Meta, response.Data)
  } else {
	return []*ClickHouseFrame{NewFrame(query.RefId, query.RefId, response.Meta, response.Data)}
  }
}

func NewFrame(refId string, name string, fieldsMeta []*FieldMeta, data []map[string]interface{}) *ClickHouseFrame {
  // TODO use fieldsMeta to define fields size for better performance
  fields:= make([]*ClickHouseField, 0)

  for _, meta := range fieldsMeta {
	fields = append(fields, NewField(meta.Name, meta.Type))
  }

  frame := &ClickHouseFrame{
	RefId: refId,
	Name: name,
	Fields: fields,
  }

  frame.SetData(data)

  return frame
}

func (f *ClickHouseFrame) getField(name string) *ClickHouseField {
  for _, field := range f.Fields {
	if field != nil && field.Name == name {
	  return field
	}
  }

  return nil
}

func (f *ClickHouseFrame) AddRow(row map[string]interface{})  {
	for key, value := range row {
		field := f.getField(key)
		if field != nil {
			field.Append(value)
		} else {
			backend.Logger.Warn(fmt.Sprintf("Unable to find field with name '%s', values are not added", key))
		}
	}
}

func (f *ClickHouseFrame) SetData(data []map[string]interface{})  {
  for _, row := range data {
	f.AddRow(row)
  }
}

func (f *ClickHouseFrame) ToDataFrame() *data.Frame  {
	fields := make([]*data.Field, 0)

	for _, f := range f.Fields {
		fields = append(fields, f.toFrameField())
	}

	frame := &data.Frame{
		RefID: f.RefId,
		Name: f.Name,
		Fields: fields,
	}

	return frame
}
