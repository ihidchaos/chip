package session

import "math/rand"

const kMessageCounterRandomInitMask uint32 = 0x0FFFFFFF

type MessageCounterType uint8

const (
	GlobalUnencrypted = iota
	GlobalEncrypted
	CounterTypeSession
)

type MessageCounterBase interface {
	AdvanceAndConsume() (uint32, error)
	Type() MessageCounterType
}

const kMessageCounterMax uint32 = 0xFFFFFFFF

type GlobalUnencryptedMessageCounter struct {
	LastUsedValue uint32
}

func NewGlobalUnencryptedMessageCounter() *GlobalUnencryptedMessageCounter {
	return &GlobalUnencryptedMessageCounter{
		LastUsedValue: rand.Uint32() & kMessageCounterRandomInitMask,
	}
}

func (g GlobalUnencryptedMessageCounter) AdvanceAndConsume() (uint32, error) {
	g.LastUsedValue = g.LastUsedValue + 1
	return g.LastUsedValue, nil
}

func (g GlobalUnencryptedMessageCounter) Type() MessageCounterType {
	return GlobalUnencrypted
}

type LocalMessageCounter struct {
	LastUsedValue uint32
}

func NewLocalMessageCounter() *LocalMessageCounter {
	return &LocalMessageCounter{
		LastUsedValue: rand.Uint32() & kMessageCounterRandomInitMask,
	}
}

func (l LocalMessageCounter) AdvanceAndConsume() (uint32, error) {
	if l.LastUsedValue == kMessageCounterMax {
		return 0, ErrorMessageCounterExhausted
	}
	l.LastUsedValue = l.LastUsedValue + 1
	return l.LastUsedValue, nil
}

func (l LocalMessageCounter) Type() MessageCounterType {
	return CounterTypeSession
}
