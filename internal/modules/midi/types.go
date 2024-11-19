package midi

type stops map[string]func()

type Listeners struct {
	ControlChange func(channel, control, velocity uint8)
	NoteStart     func(channel, control uint8)
	NoteEnd       func(channel, control uint8)
}
