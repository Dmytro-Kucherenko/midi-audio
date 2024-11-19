package playback

import (
	"os/exec"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Run(command ...string) error {
	cmd := exec.Command("playerctl", command...)

	return cmd.Run()
}

func (service *Service) Toggle() error {
	return service.Run("play-pause")
}

func (service *Service) Next() error {
	return service.Run("next")
}

func (service *Service) Prev() error {
	return service.Run("previous")
}
