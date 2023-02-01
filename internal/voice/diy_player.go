//go:build !darwin

package voice

func init() {
	player = &DiyPlayer{}
}

type DiyPlayer struct {
}

func (d *DiyPlayer) Play(voicePath string) {
	panic("This feature is not implemented yet")
}
