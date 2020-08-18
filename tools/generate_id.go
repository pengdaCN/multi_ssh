package tools

import "sync"

var (
	generateID *IdSeed
)

type IdSeed struct {
	id int
	sync.RWMutex
}

func init() {
	generateID = new(IdSeed)
}

func (i *IdSeed) GetID() int {
	rst := i.id
	i.Lock()
	i.id++
	i.Unlock()
	return rst
}

func SetGenerateIDSeed(i *IdSeed) {
	generateID = i
}

func GetID() int {
	return generateID.GetID()
}
