package email

type Sender interface {
	Send(to string, subject string, html string) error
}
