package chat

type UChat interface {
	IsBusy() bool
	Start(addr string) error
	End() error
	Send(msg string) error
	Recive() (string, error)
}
