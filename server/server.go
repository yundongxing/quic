package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"log"
	"os"
//    "time"
	quic "github.com/lucas-clemente/quic-go"
)

// FileServer 文件服务端
type FileServer struct {
	Address    string
	TLSConfig  *tls.Config
	QuicConfig *quic.Config
	Sessions   map[int64]*quic.Session
	Listener   quic.Listener
	Ctx        context.Context
}

// NewFileServer 创建FileServer对象
func NewFileServer(address string, tlsConfig *tls.Config, quicConfig *quic.Config) *FileServer {
	return &FileServer{
		Address:    address,
		TLSConfig:  tlsConfig,
		QuicConfig: quicConfig,
		Sessions:   make(map[int64]*quic.Session, 0),
		Ctx:        context.Background(),
	}
}

// Run 启动服务端
func (s *FileServer) Run() error {
	var err error
	s.Listener, err = quic.ListenAddr(s.Address, s.TLSConfig, s.QuicConfig)
	result:="the server is listening\n"
	rs.Append(result)
	if err != nil {
		log.Print("listen error: %v\n", err)
	}
	for {
		sess, err := s.Listener.Accept(s.Ctx)
		if err != nil {
			log.Print("accept session error: %v\n", err)
			continue
		} 
		sessionHandler :=NewSessionHandler(&sess)
		go sessionHandler.Run()
		continue
	//	time.Sleep(time.Second*120)
	}
}
//写入文件
func writeFile(file string) (*bufio.Writer, error) {
	fp, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return bufio.NewWriter(fp), nil
}
