package midi

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type Service struct {
	stops stops
}

func NewService() *Service {
	return &Service{stops: stops{}}
}

func (service *Service) listen(in *drivers.In, listeners *Listeners) (stop func(), err error) {
	return midi.ListenTo(*in, func(message midi.Message, timestampms int32) {
		var channel, control, velocity uint8

		if ok := message.GetControlChange(&channel, &control, &velocity); ok && listeners.ControlChange != nil {
			listeners.ControlChange(channel, control, velocity)
		}

		if ok := message.GetNoteStart(&channel, &control, &velocity); ok && listeners.NoteStart != nil {
			listeners.NoteStart(channel, control)
		}

		if ok := message.GetNoteEnd(&channel, &control); ok && listeners.NoteEnd != nil {
			listeners.NoteEnd(channel, control)
		}
	}, midi.UseSysEx())
}

func (service *Service) CheckRunning(port string) bool {
	_, ok := service.stops[port]

	return ok
}

func (service *Service) Stop(port string) {
	if stop, ok := service.stops[port]; ok {
		stop()
		delete(service.stops, port)
	}
}

func (service *Service) Close() {
	for _, stop := range service.stops {
		stop()
	}

	midi.CloseDriver()
}

func (service *Service) Listen(ports []string, listeners *Listeners, overwrite bool) {
	for _, port := range ports {
		running := service.CheckRunning(port)
		in, err := midi.FindInPort(port)
		if err == nil && running && overwrite {
			continue
		}

		if err != nil {
			if running {
				service.Stop(port)
			}

			continue
		}

		stop, err := service.listen(&in, listeners)
		if err == nil {
			service.stops[port] = stop
		}
	}
}

func (service *Service) GetPortNames() []string {
	inPorts := midi.GetInPorts()

	var ports []string
	for _, inPort := range inPorts {
		ports = append(ports, inPort.String())
	}

	return ports
}
