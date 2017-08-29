package sms

type Sender interface {
	Send(prefix string, userId uint, phone string) error
	Verify(prefix string, userId uint, phone, code string) bool
}

type Vendor interface {
	Send(phone, code string) error
}
