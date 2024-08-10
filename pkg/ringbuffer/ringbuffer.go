package ringbuffer

import (
	"errors"
)

type RingBuffer struct {
	frontPos int32
	rearPos  int32

	defaultSize int32

	Buffer []byte
}

func NewRingBuffer(size int32) *RingBuffer {
	return &RingBuffer{
		frontPos:    0,
		rearPos:     0,
		Buffer:      make([]byte, size, size),
		defaultSize: size,
	}
}

func (c *RingBuffer) Enqueue(data *[]byte, size int32) int32 {
	var tmpRearPos int32 = c.rearPos
	var tmpFrontPos int32 = c.frontPos
	var ret_val int32 = 0

	for size > 0 {
		if ((tmpRearPos + 1) % c.defaultSize) == tmpFrontPos {
			break
		}

		//(*c.Buffer)[tmpRearPos] = *data[]
		c.Buffer[tmpRearPos] = (*data)[ret_val]
		tmpRearPos = (tmpRearPos + 1) % c.defaultSize
		ret_val++
		size--
	}
	c.rearPos = tmpRearPos
	return ret_val
}

func (c *RingBuffer) Dequeue(data *[]byte, size int32) (int32, error) {
	var tmpRearPos int32 = c.rearPos
	var tmpFrontPos int32 = c.frontPos
	var orgFrontPos int32 = c.frontPos
	var orgSize int32 = size
	circleCheck := false

	loopCount := size
	var retCount int32 = 0
	for loopCount > 0 {

		if tmpFrontPos == tmpRearPos {
			return 0, errors.New("over size!")
		}

		tmpFrontPos = (tmpFrontPos + 1) % c.defaultSize

		if tmpFrontPos == 0 {
			circleCheck = true
		}
		loopCount--
		retCount++
	}

	if circleCheck {
		var endSpace int32 = c.defaultSize - orgFrontPos
		(*data) = append((*data), c.Buffer[orgFrontPos:c.defaultSize]...)
		(*data) = append((*data), c.Buffer[:(orgSize-endSpace)]...)
	} else {
		(*data) = append((*data), c.Buffer[orgFrontPos:orgFrontPos+size]...)
	}

	c.frontPos = tmpFrontPos
	return retCount, nil
}

func (c *RingBuffer) Peek(data *[]byte, size int32) (int32, error) {
	var tmpRearPos int32 = c.rearPos
	var tmpFrontPos int32 = c.frontPos
	var orgFrontPos int32 = c.frontPos
	var orgSize int32 = size
	circleCheck := false

	loopCount := size
	var retCount int32 = 0
	for loopCount > 0 {

		if tmpFrontPos == tmpRearPos {
			return 0, errors.New("over size!")
		}

		tmpFrontPos = (tmpFrontPos + 1) % c.defaultSize

		if tmpFrontPos == 0 {
			circleCheck = true
		}
		loopCount--
		retCount++
	}

	if circleCheck {
		var endSpace int32 = c.defaultSize - orgFrontPos
		(*data) = append((*data), c.Buffer[orgFrontPos:endSpace]...)
		(*data) = append((*data), c.Buffer[:(orgSize-endSpace)]...)
	} else {
		(*data) = append((*data), c.Buffer[orgFrontPos:size]...)
	}

	return retCount, nil
}

func (c *RingBuffer) DirectDequeueSize() int32 {
	var tmpRearPos int32 = c.rearPos
	var tmpFrontPos int32 = c.frontPos

	if tmpRearPos >= tmpFrontPos {
		return tmpRearPos - tmpFrontPos
	} else {
		return c.defaultSize - tmpFrontPos
	}
}

func (c *RingBuffer) GetRearPos() []byte {
	return c.Buffer[c.rearPos:]
}

func (c *RingBuffer) GetFrontPos() []byte {
	return c.Buffer[c.frontPos:]
}

func (c *RingBuffer) MoveRearPos(move int32) {
	c.rearPos += move
}

func (c *RingBuffer) Clear() {
	c.frontPos = 0
	c.rearPos = 0
}
