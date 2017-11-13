package main

var (
	EXPrefix = "BULUSI.EX."
	QPrefix  = "BULUSI.QU."
)

type Service struct {
	name       string
	id         int //Id must be unique for one service
	interfaces []*ServicePort
	rpc        *RPC
}

func NewService(sname, host string, id int) *Service {
	return &Service{
		name: sname,
		id:   id,
		rpc:  NewRPC(sname, host),
	}
}

func (this *Service) AddInterface(si *ServicePort) {
	this.interfaces = append(this.interfaces, si)
}

func (this *Service) Start() {
	for _, iface := range this.interfaces {
		this.rpc.AddQueue(QPrefix+this.name+"."+string(this.id), iface.name)
	}
}
