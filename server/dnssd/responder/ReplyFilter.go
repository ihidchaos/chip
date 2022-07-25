package responder

type ReplyFilter interface {
	Accept(t uint16, c uint16, n string) bool
}
