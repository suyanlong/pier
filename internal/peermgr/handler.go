package peermgr

import (
	"fmt"

	"github.com/link33/sidercar/model/pb"
)

func (swarm *Swarm) RegisterMultiMsgHandler(messageTypes []pb.Message_Type, handler MessageHandler) error {
	for _, typ := range messageTypes {
		if err := swarm.RegisterMsgHandler(typ, handler); err != nil {
			return err
		}
	}
	return nil
}

func (swarm *Swarm) RegisterMsgHandler(messageType pb.Message_Type, handler MessageHandler) error {
	if handler == nil {
		return fmt.Errorf("register msg handler: empty handler")
	}
	for msgType := range pb.Message_Type_name {
		if msgType == int32(messageType) {
			swarm.msgHandlers.Store(messageType, handler)
			return nil
		}
	}

	return fmt.Errorf("register msg handler: invalid message type")
}

func (swarm *Swarm) RegisterConnectHandler(handler ConnectHandler) error {
	swarm.lock.Lock()
	defer swarm.lock.Unlock()

	swarm.connectHandlers = append(swarm.connectHandlers, handler)

	return nil
}
