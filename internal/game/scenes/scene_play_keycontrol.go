package gamescene

type KeyControl struct {
	isLeftPressed  bool
	isDownPressed  bool
	isUpPressed    bool
	isRightPressed bool
}

func NewKeyControl() *KeyControl {
	return &KeyControl{}
}

func (k *KeyControl) Reset() {
	k.isLeftPressed = false
	k.isDownPressed = false
	k.isUpPressed = false
	k.isRightPressed = false
}

func (k *KeyControl) IsSomeKeyPressed() bool {
	return k.isLeftPressed ||
		k.isDownPressed ||
		k.isUpPressed ||
		k.isRightPressed
}

func (k *KeyControl) PressLeft() {
	k.isLeftPressed = true
}
func (k *KeyControl) PressDown() {
	k.isDownPressed = true
}
func (k *KeyControl) PressUp() {
	k.isUpPressed = true
}
func (k *KeyControl) PressRight() {
	k.isRightPressed = true
}
