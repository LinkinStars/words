//go:build darwin

package voice

import (
	"os/exec"

	"github.com/LinkinStars/go-scaffold/logger"
)

func init() {
	player = &AfPlayer{}
}

type AfPlayer struct {
}

func (a *AfPlayer) Play(voicePath string) {
	err := exec.Command("afplay", voicePath).Run()
	if err != nil {
		logger.Error(err.Error())
	}
}
