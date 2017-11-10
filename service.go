package main

var (
	EXPrefix = "BULUSI.EX."
	QPrefix  = "BULUSI.QU."
)

type Service struct {
	name       string
	id         int //Id must be unique for one service
	interfaces []*ServiceInterface
	rpc        *RPC
}

func NewService(sname string, id int) *Service {
	return &Service{
		name: sname,
		id:   id,
		rpc:  NewRPC(sname),
	}
}

func (this *Service) AddInterface(si *ServiceInterface) {
	this.interfaces = append(this.interfaces, si)
}

func (this *Service) Start() {
	for iface, _ := range this.interfaces {
		this.rpc.AddQueue(QPrefix+this.name+"."+string(id), iface.name)
	}
}
