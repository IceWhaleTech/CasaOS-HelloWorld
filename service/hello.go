package service

type Hello struct{}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Ping() string {
	return "pong"
}
