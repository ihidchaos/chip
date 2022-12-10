package session

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
)

const InitialSyncValue uint32 = 0

type Synced struct {
	mMaxCounter uint32
	mWindow     *bitset.BitSet
}

type position uint8
type peerMessageCounterStatus uint8

const (
	kNotSynced peerMessageCounterStatus = iota
	kSyncInProcess
	kSynced
	kChallengeSize = 8
)

const (
	kBeforeWindow position = iota
	kInWindow
	kMaxCounter
	kFutureCounter
)

type SyncInProcess struct {
	mChallenge []byte
	mSize      int
}

type PeerMessageCounter struct {
	mState  peerMessageCounterStatus
	mSynced *Synced
}

func NewPeerMessageCounter() *PeerMessageCounter {
	return &PeerMessageCounter{
		mState: kNotSynced,
	}
}

func (c *PeerMessageCounter) VerifyUnencrypted(counter uint32) error {
	switch c.mState {
	case kNotSynced:
		c.SetCounter(counter)
		return nil
	case kSynced:
		pos := c.classifyWithRollover(counter)
		return c.VerifyPositionUnencrypted(pos, counter)

	default:
		return lib.ErrorInternal
	}
}

func (c *PeerMessageCounter) SetCounter(counter uint32) {
	c.reset()
	c.mState = kSynced
	c.mSynced = &Synced{
		mMaxCounter: counter,
		mWindow:     bitset.New(config.MessageCounterWindowSize),
	}
	c.mSynced.mWindow.ClearAll()
}

func (c *PeerMessageCounter) reset() {
	switch c.mState {
	case kNotSynced:
		break
	case kSyncInProcess:
		break
	case kSynced:
	}
}

func (c *PeerMessageCounter) CommitUnencrypted(counter uint32) {
	c.commitWithRollover(counter)
}

func (c *PeerMessageCounter) commitWithRollover(counter uint32) {
	//TODO
	c.classifyWithoutRollover(counter)
}

func (c *PeerMessageCounter) classifyWithRollover(counter uint32) position {
	return kBeforeWindow
}

func (c *PeerMessageCounter) VerifyPositionUnencrypted(pos position, counter uint32) error {
	return nil
}

func (c *PeerMessageCounter) CommitUnencryptedUnicast(counter uint32) {
	c.commitWithRollover(counter)
}

func (c *PeerMessageCounter) classifyWithoutRollover(counter uint32) position {
	if counter > c.mSynced.mMaxCounter {
		return kFutureCounter
	}
	return c.classifyNonFutureCounter(counter)
}

func (c *PeerMessageCounter) classifyNonFutureCounter(counter uint32) position {
	if counter == c.mSynced.mMaxCounter {
		return kMaxCounter
	}
	return kBeforeWindow
}

func (c *PeerMessageCounter) VerifyOrTrustFirstGroup(counter uint32) error {
	switch c.mState {
	case kNotSynced:
		c.SetCounter(counter)
		return nil
	case kSynced:
		return c.VerifyGroup(counter)
	default:
		return lib.ErrorInternal
	}
}

func (c *PeerMessageCounter) VerifyGroup(counter uint32) error {
	if c.mState != kSynced {
		return lib.IncorrectState
	}
	pos := c.classifyWithRollover(counter)
	return c.VerifyPositionEncrypted(pos, counter)
}

func (c *PeerMessageCounter) VerifyPositionEncrypted(pos position, counter uint32) error {
	switch pos {
	case kFutureCounter:
		return nil
	case kInWindow:
		offset := c.mSynced.mMaxCounter - counter
		if c.mSynced.mWindow.Test(uint(offset - 1)) {
			return lib.DuplicateMessageReceived
		}
		return nil
	default:
		return lib.DuplicateMessageReceived
	}
}

func (c *PeerMessageCounter) CommitGroup(counter uint32) {
	c.commitWithRollover(counter)
}

func (c *PeerMessageCounter) VerifyEncryptedUnicast(counter uint32) error {
	return nil
}
