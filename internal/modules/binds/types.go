package binds

type Bind struct {
	Apps   []string `json:"apps" validate:"dive,required"`
	Active bool     `json:"active"`
	Output bool     `json:"output"`
	Input  bool     `json:"input"`
}

type ControlChangeBind struct {
	Bind
	Action ChangeAction `json:"action" validate:"required"`
}

type ControlBind struct {
	Change *ControlChangeBind `json:"change"`
}

type NoteStartBind struct {
	Bind
	Action ClickAction `json:"action" validate:"required"`
}

type NoteEndBind struct {
	Bind
	Action ClickAction `json:"action" validate:"required"`
}

type NotelBind struct {
	Start *NoteStartBind `json:"start"`
	End   *NoteEndBind   `json:"end"`
}

type ChannelBind struct {
	Controls map[uint8]ControlBind `json:"controls" validate:"required,dive,required"`
	Notes    map[uint8]NotelBind   `json:"notes" validate:"required,dive,required"`
}

type Schema struct {
	Ports    []string              `json:"ports" validate:"required,dive,required"`
	Channels map[uint8]ChannelBind `json:"channels" validate:"required,dive,required"`
}
