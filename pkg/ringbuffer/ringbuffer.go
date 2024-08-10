package ringbuffer

import (
	"errors"
)

type Buffer struct {
	frontPos int32
	rearPos  int32

	defaultSize int32
	Buffer      []byte
}

func NewBuffer(size int32) *Buffer {
	return &Buffer{
		frontPos:    0,
		rearPos:     0,
		Buffer:      make([]byte, size, size),
		defaultSize: size,
	}
}

func (c *Buffer) Enqueue(data *[]byte, size int32) int32 {
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

func (c *Buffer) Dequeue(data *[]byte, size int32) (int32, error) {
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

func (c *Buffer) Peek(data *[]byte, size int32) (int32, error) {
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

func (c *Buffer) DirectDequeueSize() int32 {
	var tmpRearPos int32 = c.rearPos
	var tmpFrontPos int32 = c.frontPos

	if tmpRearPos >= tmpFrontPos {
		return tmpRearPos - tmpFrontPos
	} else {
		return c.defaultSize - tmpFrontPos
	}
}

/*
Dequeue 할 수 있는 사이즈
*/
func (c *Buffer) GetUseSize() int32 {
	tmpRear := c.rearPos
	tmpFront := c.frontPos
	if tmpRear >= tmpFront {
		return tmpRear - tmpFront
	}
	return (c.defaultSize - tmpFront) - tmpRear
}

func (c *Buffer) GetRearPos() []byte {
	return c.Buffer[c.rearPos:]
}

func (c *Buffer) GetFrontPos() []byte {
	return c.Buffer[c.frontPos:]
}

func (c *Buffer) MoveRearPos(move int32) {
	c.rearPos += move
}

func (c *Buffer) Clear() {
	c.frontPos = 0
	c.rearPos = 0
}
