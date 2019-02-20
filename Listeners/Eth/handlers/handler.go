package handlers

type IHandler interface {
}

type handler struct {
}

func New() IHandler {
	return &handler{}
}
