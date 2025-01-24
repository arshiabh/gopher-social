package mail

type Client interface {
	Send(string, string) error
}
