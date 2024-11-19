package audio

import (
	"slices"
	"strconv"
	"strings"

	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
)

type Service struct {
	client *pulse.Client
}

func NewService() *Service {
	client, err := pulse.NewClient()
	if err != nil {
		panic("Failed to initialize pulse client")
	}

	return &Service{client}
}

func (service *Service) Close() {
	service.client.Close()
}

func getAppName(app *proto.GetSinkInputInfoReply) string {
	return app.Properties["application.process.binary"].String()
}

func getAppPID(app *proto.GetSinkInputInfoReply) (PID uint32, err error) {
	value, err := strconv.ParseUint(app.Properties["application.process.id"].String(), 10, 32)
	if err != nil {
		return
	}

	return uint32(value), nil
}

func (service *Service) getApps(filters *GetAppsFilters) (apps []proto.GetSinkInputInfoReply, err error) {
	list := &proto.GetSinkInputInfoListReply{}
	err = service.client.RawRequest(&proto.GetSinkInputInfoList{}, list)
	if err != nil {
		return
	}

	for _, app := range *list {
		appPID, err := getAppPID(app)
		appName := getAppName(app)

		formatName := func(name string) string {
			return strings.ReplaceAll(strings.ToLower(name), " ", "-")
		}

		matchName := func(currentName string) bool {
			appNameFormatted := formatName(appName)
			currentNameFormatted := formatName(currentName)

			return strings.EqualFold(appNameFormatted, currentNameFormatted) ||
				strings.Contains(currentNameFormatted, appNameFormatted) ||
				strings.Contains(appNameFormatted, currentNameFormatted)
		}

		if (!filters.Names.Valid && !filters.Indices.Valid) || (filters.Names.Valid && slices.ContainsFunc(filters.Names.Value, matchName)) ||
			(filters.Indices.Valid && err == nil && slices.Contains(filters.Indices.Value, appPID)) {
			apps = append(apps, *app)

			continue
		}
	}

	return
}

func (service *Service) GetAppNames() []string {
	var names []string
	apps, err := service.getApps(&GetAppsFilters{})

	if err != nil {
		return names
	}

	for _, app := range apps {
		names = append(names, getAppName(&app))
	}

	return names
}

func (service *Service) SetAppVolume(filters *GetAppsFilters, volume uint32) error {
	apps, err := service.getApps(filters)
	if err != nil {
		return err
	}

	for _, app := range apps {
		service.client.RawRequest(&proto.SetSinkInputVolume{
			SinkInputIndex: app.SinkInputIndex,
			ChannelVolumes: proto.ChannelVolumes{volume, volume},
		}, nil)
	}

	return nil
}

func (service *Service) ToggleAppMute(filters *GetAppsFilters) error {
	apps, err := service.getApps(filters)
	if err != nil {
		return err
	}

	for _, app := range apps {
		service.client.RawRequest(&proto.SetSinkInputMute{
			SinkInputIndex: app.SinkInputIndex,
			Mute:           !app.Muted,
		}, nil)
	}

	return nil
}

func (service *Service) getDefaultOutput() (*proto.GetSinkInfoReply, error) {
	var output proto.GetSinkInfoReply

	err := service.client.RawRequest(&proto.GetSinkInfo{}, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (service *Service) SetOutputVolume(volume uint32) error {
	output, err := service.getDefaultOutput()
	if err != nil {
		return err
	}

	return service.client.RawRequest(&proto.SetSinkVolume{
		SinkIndex:      output.SinkIndex,
		ChannelVolumes: proto.ChannelVolumes{volume, volume},
	}, nil)
}

func (service *Service) ToggleOutputMute() error {
	output, err := service.getDefaultOutput()
	if err != nil {
		return err
	}

	return service.client.RawRequest(&proto.SetSinkMute{
		SinkIndex: output.SinkIndex,
		Mute:      !output.Mute,
	}, nil)
}

func (service *Service) getDefaultInput() (*proto.GetSourceInfoReply, error) {
	var input proto.GetSourceInfoReply

	err := service.client.RawRequest(&proto.GetSourceInfo{}, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func (service *Service) SetInputVolume(volume uint32) error {
	input, err := service.getDefaultInput()
	if err != nil {
		return err
	}

	return service.client.RawRequest(&proto.SetSourceVolume{
		SourceIndex:    input.SourceIndex,
		ChannelVolumes: proto.ChannelVolumes{volume, volume},
	}, nil)
}

func (service *Service) ToggleInputMute() error {
	input, err := service.getDefaultInput()
	if err != nil {
		return err
	}

	return service.client.RawRequest(&proto.SetSourceMute{
		SourceIndex: input.SourceIndex,
		Mute:        !input.Mute,
	}, nil)
}
