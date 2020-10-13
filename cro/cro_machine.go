package cro

import (
	"io"
	"io/ioutil"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"multi_ssh/tools"
	"sort"
	"time"
)

type (
	CroMachine struct {
		terminals []m_terminal.Terminal
	}
	Task struct {
		// 主机选择条件
		hostsF string
		// 设置执行任务最长时间
		timeout time.Duration
		// 原始的主机信息
		rawHostsInfo string
		// 输出信息格式
		format string
		// 任务信息输出位置
		out    io.Writer
		byTerm []m_terminal.Terminal
	}
	userSlice []model.SHHUser
)

func (u userSlice) Less(v1, v2 int) bool {
	return u[v1].Line() < u[v2].Line()
}

func (u userSlice) Swap(v1, v2 int) {
	u[v1], u[v2] = u[v2], u[v1]
}

func (u userSlice) Len() int {
	return len(u)
}

func ReadF(filename string) *taskBuilder {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var us userSlice
	u, err := model.ReadLines(tools.ByteSlice2String(f))
	if err != nil {
		panic(err)
	}
	for _, v := range u {
		us = append(us, v)
	}
	sort.Sort(us)

	return nil
}

func ReadS(hostsInfo string) *taskBuilder {
	return &taskBuilder{
		rawHostsInfo: hostsInfo,
	}
}
