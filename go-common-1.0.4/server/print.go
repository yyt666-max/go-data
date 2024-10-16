package server

type ServerDetail struct {
}

func (i *imlServiceBuilder) Detail() *ServerDetail {
	return &ServerDetail{}
}
