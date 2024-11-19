package history

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
)

type Service struct {
	path string
}

func NewService(path string) *Service {
	return &Service{path}
}

func ExampleSchema() *Schema {
	return &Schema{
		Status:   "inspected",
		Messages: []string{},
		Ports:    []string{},
		Apps:     []string{},
	}
}

func (service *Service) read() (*Schema, error) {
	file, err := os.Open(service.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var schema Schema
	if err = json.NewDecoder(file).Decode(&schema); err != nil {
		return nil, err
	}

	return &schema, nil
}

func (service *Service) write(schema *Schema) error {
	file, err := os.OpenFile(service.path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(schema)
}

func (service *Service) Add(ports []string, apps []string) (schema *Schema, err error) {
	schema, err = service.read()
	if err != nil {
		return
	}

	for _, port := range ports {
		if !slices.Contains(schema.Ports, port) {
			schema.Ports = append(schema.Ports, port)
		}
	}

	for _, name := range apps {
		if !slices.Contains(schema.Apps, name) {
			schema.Apps = append(schema.Apps, name)
		}
	}

	err = service.write(schema)
	if err != nil {
		return nil, err
	}

	return
}

func (service *Service) Refresh(status Status, message string) error {
	schema, err := service.read()
	if err != nil {
		return err
	}

	schema.Status = status
	if message != "" {
		schema.Messages = strings.Split(message, "\n")
	} else {
		schema.Messages = []string{}
	}

	return service.write(schema)
}
