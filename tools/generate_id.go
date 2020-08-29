package tools

import "sync/atomic"

var (
	generateID *IdSeed
)

type IdSeed struct {
	id int32
}

func init() {
	generateID = new(IdSeed)
}

func (i *IdSeed) GetID() int {
	return int(atomic.AddInt32(&i.id, 1))
}

func SetGenerateIDSeed(i *IdSeed) {
	generateID = i
}

func GetID() int {
	return generateID.GetID()
}
