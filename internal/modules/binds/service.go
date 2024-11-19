package binds

import (
	"encoding/json"
	"os"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	path   string
	schema *Schema
}

func NewService(path string) *Service {
	return &Service{path, nil}
}

func SchemaExample() *Schema {
	return &Schema{Ports: []string{}, Channels: map[uint8]ChannelBind{}}
}

func (service *Service) Schema() *Schema {
	return service.schema
}

func (service *Service) Parse() (*Schema, error) {
	file, err := os.Open(service.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var schema Schema
	if err := json.NewDecoder(file).Decode(&schema); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(schema); err != nil {
		return nil, err
	}

	service.schema = &schema

	return service.schema, nil
}
