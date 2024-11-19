package internal

import (
	"github.com/Dmytro-Kucherenko/users-sam/internal/common/helpers"
	"github.com/Dmytro-Kucherenko/users-sam/internal/common/types"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/audio"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/binds"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/playback"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/process"
)

func UseControlChange(
	configService *binds.Service,
	audioService *audio.Service,
	processService *process.Service,
) func(channel, control, velocity uint8) {
	return func(channel, control, velocity uint8) {
		bind := configService.Schema().Channels[channel].Controls[control].Change
		if bind == nil {
			return
		}

		switch bind.Action {
		case binds.VOLUME:
			volume := helpers.ConvertVelocityToVolume(velocity)
			if bind.Apps != nil {
				audioService.SetAppVolume(&audio.GetAppsFilters{
					Names: types.NewOptional(bind.Apps),
				}, volume)
			}

			if bind.Active {
				processID, name, err := processService.GetActive()

				if err == nil {
					audioService.SetAppVolume(&audio.GetAppsFilters{
						Names:   types.NewOptional([]string{name}),
						Indices: types.NewOptional([]uint32{processID}),
					}, volume)
				}
			}

			if bind.Output {
				audioService.SetOutputVolume(volume)
			}

			if bind.Input {
				audioService.SetInputVolume(volume)
			}
		}
	}
}

func UseNoteStart(
	configService *binds.Service,
	audioService *audio.Service,
	playbackService *playback.Service,
	processService *process.Service,
) func(channel, control uint8) {
	return func(channel, control uint8) {
		bind := configService.Schema().Channels[channel].Notes[control].Start
		if bind == nil {
			return
		}

		switch bind.Action {
		case binds.MUTE:
			if bind.Apps != nil {
				audioService.ToggleAppMute(&audio.GetAppsFilters{
					Names: types.NewOptional(bind.Apps),
				})
			}

			if bind.Active {
				processID, name, err := processService.GetActive()

				if err == nil {
					audioService.ToggleAppMute(&audio.GetAppsFilters{
						Names:   types.NewOptional([]string{name}),
						Indices: types.NewOptional([]uint32{processID}),
					})
				}
			}

			if bind.Output {
				audioService.ToggleOutputMute()
			}

			if bind.Input {
				audioService.ToggleInputMute()
			}
		case binds.TOGGLE:
			playbackService.Toggle()
		case binds.PREVIOUS:
			playbackService.Prev()
		case binds.NEXT:
			playbackService.Next()
		}
	}
}
