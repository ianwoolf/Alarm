package models

import (
	"fmt"

	m "github.com/lodastack/models"
)

func init() {
	Work = make(map[string]NsBlock)
}

var (
	SendBlock = "block"
	SendPoint = "point"
)

type BlockWork struct {
	NS          string
	Measurement string
	SendType    string

	block  int // Alarm ADD
	period int // Alarm ADD
	cached []*Point
}

type AlarmBlock struct {
	Alarm  m.Alarm // TODO
	ByNs   BlockWork
	ByHost map[string]BlockWork
}

func NewAlarmBlock(ns, sendType string, alarm m.Alarm) AlarmBlock {
	block := AlarmBlock{
		Alarm: alarm,
		ByNs: BlockWork{NS: ns,
			Measurement: alarm.Measurement,
			SendType:    sendType,
			// block:0,
			// period:0,
			// cached: make([]*Point, block),
		},
		ByHost: make(map[string]BlockWork),
	}
	return block
}

func (b *BlockWork) Log() error    { return nil }
func (b *BlockWork) QueryHistory() {}
func (b *BlockWork) Check() bool   { return true }
func (b *BlockWork) Send() error   { return nil }

func (ab *AlarmBlock) Event(point Point) error { return nil }

type NsBlock struct {
	AlarmBlocks map[string]AlarmBlock
	// Mu
}

var Work map[string]NsBlock

func Println() {
	for k, v := range Work {
		fmt.Println(k, v)
	}
}
