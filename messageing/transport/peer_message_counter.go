package transport

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
)

type Synced struct {
	mMaxCounter uint32
	mWindow     *bitset.BitSet
}

type SyncInProcess struct {
	mChallenge []byte
	mSize      int
}

const (
	sNotSynced uint8 = iota
	sSyncInProcess
	sSynced
	kChallengeSize = 8
)

type PeerMessageCounter struct {
	mState  uint8
	mSynced *Synced
}

func NewPeerMessageCounter() *PeerMessageCounter {
	return &PeerMessageCounter{}
}

func (c *PeerMessageCounter) VerifyUnencrypted(counter uint32) error {
	switch c.mState {
	case sNotSynced:
		c.SetCounter(counter)
		return nil
	default:
		return lib.MatterErrorInternal
	}
}

func (c *PeerMessageCounter) SetCounter(counter uint32) {
	c.reset()
	c.mState = sSynced
	c.mSynced = &Synced{
		mMaxCounter: counter,
		mWindow:     bitset.New(config.ChipConfigMessageCounterWindowSize),
	}
	c.mSynced.mWindow.ClearAll()
}

func (c *PeerMessageCounter) reset() {
	switch c.mState {
	case sNotSynced:
		break
	case sSyncInProcess:
		break
	case sSynced:
	}
}

func (c *PeerMessageCounter) CommitUnencrypted(counter uint32) {
	c.commitWithRollover(counter)
}

func (c *PeerMessageCounter) commitWithRollover(counter uint32) {
	//TODO
}
