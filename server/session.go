package main

import (
	"context"
	"log"
//	"time"
//	"github.com/873314461/quic-file/common"
	"github.com/lucas-clemente/quic-go"
)

type SessionHandler struct {
	Ctx     context.Context
	Session quic.Session
	Streams map[int64]*StreamHandler
}
//  创建会话对象
func NewSessionHandler(session *quic.Session) *SessionHandler {
	return &SessionHandler{
		Ctx:     context.Background(),
		Session: *session,
		Streams: make(map[int64]*StreamHandler, 0),
	}
}
//  会话运行
func (h *SessionHandler) Run(){//(chquit chan bool) 
	for {
		stream, err := h.Session.AcceptStream(h.Ctx)
		if err != nil {
			if err.Error() != NoError {
				log.Fatalf("accept stream error: %v", err)
			}
			break
		}
		streamHandler := NewStreamHandler(&stream)
		go streamHandler.Run()
		continue
		//time.Sleep(1000*time.Second)
	}
}
