package process

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (service *Service) GetActive() (ID uint32, name string, err error) {
	root, err := exec.Command("xprop", "-root").Output()
	if err != nil {
		return
	}

	activeSchema := regexp.MustCompile(`_NET_ACTIVE_WINDOW\(WINDOW\): window id # (0x[0-9a-fA-F]+)`)
	activeSubmatch := activeSchema.FindStringSubmatch(string(root))
	if len(activeSubmatch) < 2 {
		err = errors.New("root active window ID was not found ")

		return
	}

	info, err := exec.Command("xprop", "-id", activeSubmatch[1]).Output()
	if err != nil {
		return
	}

	IDSchema := regexp.MustCompile(`_NET_WM_PID\(CARDINAL\) = (\d+)`)
	IDSubmatch := IDSchema.FindStringSubmatch(string(info))
	if len(activeSubmatch) < 2 {
		err = errors.New("active window process ID was not found ")

		return
	}

	processID, err := strconv.ParseUint(IDSubmatch[1], 10, 32)
	if err != nil {
		return
	}

	nameSchema := regexp.MustCompile(`WM_CLASS\(STRING\) = "([^"]+)"`)
	nameSubmatch := nameSchema.FindStringSubmatch(string(info))
	if len(activeSubmatch) < 2 {
		err = errors.New("active window process name was not found ")

		return
	}

	return uint32(processID), nameSubmatch[1], nil
}
