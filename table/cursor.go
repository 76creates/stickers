package table

// TODO: create cursor struct that holds the cursor position and direction

// cursorDirection indicates the direction of the cursor movement, it starts at the up left direction
type cursorDirection uint8

const (
	cursorDirectionUpLeft    cursorDirection = 0b10
	cursorDirectionUpRight   cursorDirection = 0b11
	cursorDirectionDownLeft  cursorDirection = 0b01
	cursorDirectionDownRight cursorDirection = 0b00
)

func (cursor cursorDirection) setUp() cursorDirection {
	return cursor | 0b10
}

func (cursor cursorDirection) setDown() cursorDirection {
	return cursor &^ 0b10
}

func (cursor cursorDirection) setLeft() cursorDirection {
	return cursor &^ 0b01
}

func (cursor cursorDirection) setRight() cursorDirection {
	return cursor | 0b01
}

func (cursor cursorDirection) isDown() bool {
	return (cursor & (1 << 1)) == 0
}

func (cursor cursorDirection) isUp() bool {
	return (cursor & (1 << 1)) != 0
}

func (cursor cursorDirection) isLeft() bool {
	return (cursor & (1 << 0)) == 0
}

func (cursor cursorDirection) isRight() bool {
	return (cursor & (1 << 0)) != 0
}
