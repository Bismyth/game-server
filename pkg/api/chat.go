package api

import (
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
)

type m_IChat struct {
	Message string `json:"message"`
}

type m_OChat struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

const pt_IChat IPacketType = "client_chat"
const pt_OChat OPacketType = "server_chat"

func handleChat(i HandlerInput) error {
	iChatMessage, err := hp[m_IChat](i.Packet)
	if err != nil {
		return err
	}

	userName, err := db.GetUserName(i.UserId)
	if err != nil {
		return err
	}
	oChatMessage := m_OChat{
		Message: iChatMessage.Message,
		Sender:  userName,
	}

	ids, err := db.GetAllUserIds()
	if err != nil {
		return fmt.Errorf("failed to get user ids for broadcast")
	}

	oPacket := mp(pt_OChat, oChatMessage)
	SendMany(i.C, ids, &oPacket)

	return nil
}
