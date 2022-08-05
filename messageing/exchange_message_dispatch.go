package messageing

type ExchangeMessageDispatch struct {
}

func (d ExchangeMessageDispatch) IsEncryptionRequired() bool {
	return false
}
