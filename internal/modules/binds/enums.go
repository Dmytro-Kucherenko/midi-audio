package binds

type ChangeAction string

const (
	VOLUME ChangeAction = "volume"
)

type ClickAction string

const (
	MUTE     ClickAction = "mute"
	TOGGLE   ClickAction = "toggle"
	NEXT     ClickAction = "next"
	PREVIOUS ClickAction = "prev"
)
