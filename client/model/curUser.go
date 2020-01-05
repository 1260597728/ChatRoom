package model

import (
	"net"
	"sggStudyGo/chatroom/common/message"
)

// 客户端多处需要用到，因此做成全局的
type CurUser struct {
	Conn net.Conn
	message.User
}

