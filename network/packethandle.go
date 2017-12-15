package network

type PacketHandler interface {
	HandlePacket(packet interface{}, s Session)
}

type PacketHandlerImp struct {
}

func NewPacketHandlerImp() PacketHandler {
	return &PacketHandlerImp{}
}

func (self *PacketHandlerImp) HandlePacket(packet interface{}, s Session) {
	//s.SendPacket(msg)
	//fmt.Println("HandleEvent", msg)
}
