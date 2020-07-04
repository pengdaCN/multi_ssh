package main

import (
	"multi_ssh/model"
	"time"
)

type pack struct {
	time     time.Time
	userInfo model.SHHUser
	msg      []byte
}
