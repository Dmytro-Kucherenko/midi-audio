package internal

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/audio"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/binds"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/config"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/history"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/midi"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/playback"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/process"
)

func Init() {
	bindsService := binds.NewService(config.BindsPath())
	historyService := history.NewService(config.HistoryPath())
	playbackService := playback.NewService()
	processService := process.NewService()

	midiService := midi.NewService()
	defer midiService.Close()

	audioService := audio.NewService()
	defer audioService.Close()

	handleControlChange := UseControlChange(bindsService, audioService, processService)
	handleNoteStart := UseNoteStart(bindsService, audioService, playbackService, processService)

	go func() {
		for {
			time.Sleep(time.Second)

			ports := midiService.GetPortNames()
			apps := audioService.GetAppNames()
			historyService.Add(ports, apps)

			schema, err := bindsService.Parse()
			if err != nil {
				historyService.Refresh(history.SkipBinds, err.Error())
				continue
			}

			historyService.Refresh(history.Listening, "")
			midiService.Listen(schema.Ports, &midi.Listeners{
				ControlChange: handleControlChange,
				NoteStart:     handleNoteStart,
			}, true)
		}
	}()

	end, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-end.Done()
}
