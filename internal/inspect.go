package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/audio"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/binds"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/config"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/history"
	"github.com/Dmytro-Kucherenko/users-sam/internal/modules/midi"
)

func openOrCreateJSON(filePath string, value any) (created bool, err error) {
	dir := filepath.Dir(filePath)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return
	}

	_, err = os.Open(filePath)
	if err == nil {
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(value)

	created = true

	return
}

func Inspect() {
	bindsPath := config.BindsPath()
	historyPath := config.HistoryPath()

	created, err := openOrCreateJSON(bindsPath, binds.SchemaExample())
	if err != nil {
		panic(fmt.Sprintln("Failed to open or create binds file:", err.Error()))
	}

	if created {
		fmt.Println("Created binds file:", bindsPath)
	} else {
		fmt.Println("Opened binds file:", bindsPath)
	}

	created, err = openOrCreateJSON(historyPath, history.ExampleSchema())
	if err != nil {
		panic(fmt.Sprintln("Failed to open or create history file:", err.Error()))
	}

	if created {
		fmt.Println("Created history file:", bindsPath)
	} else {
		fmt.Println("Opened history file:", bindsPath)
	}

	historyService := history.NewService(config.HistoryPath())

	midiService := midi.NewService()
	defer midiService.Close()

	audioService := audio.NewService()
	defer audioService.Close()

	ports := midiService.GetPortNames()
	apps := audioService.GetAppNames()
	schema, err := historyService.Add(ports, apps)
	if err != nil {
		panic(fmt.Sprintln("Failed to fill history file:", err.Error()))
	}

	fmt.Printf("\nAvailable connections:\nPorts - %v\nApps - %v\n", strings.Join(schema.Ports, ", "), strings.Join(schema.Apps, ", "))
}
