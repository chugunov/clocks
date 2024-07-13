package main

import (
	"fmt"
)

type EventType int

const (
	local EventType = iota
	sent
	recv
)

type Event struct {
	t         EventType
	timestamp int64
	src, dst  int
}

func (e Event) String() string {
	return fmt.Sprintf("%d", e.timestamp)
}

type Message struct {
	timestamp int64
}

type Process struct {
	ID      int
	Clock   Clock
	MsgChan []chan Message
	events  []Event
}

func NewProcess(ID int, numProcesses int) *Process {
	process := new(Process)
	process.ID = ID
	process.Clock = Clock{ID: process.ID}
	for i := 0; i < numProcesses; i++ {
		msgChan := make(chan Message, 10)
		process.MsgChan = append(process.MsgChan, msgChan)
	}
	return process
}

func (p *Process) log(event Event) {
	p.events = append(p.events, event)
}

func (p *Process) Local() {
	timestamp := p.Clock.tick(0)
	p.log(Event{t: local, timestamp: timestamp})
}

func (p *Process) Send(dst *Process) {
	timestamp := p.Clock.tick(0)
	dst.MsgChan[p.ID] <- Message{timestamp}
	p.log(Event{t: sent, src: p.ID, dst: dst.ID, timestamp: timestamp})
}

func (p *Process) Recv(src int) {
	msg := <-p.MsgChan[src]
	timestamp := p.Clock.tick(msg.timestamp)
	p.log(Event{t: recv, src: src, dst: p.ID, timestamp: timestamp})
}
