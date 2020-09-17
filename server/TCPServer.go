package main

import (
//	"bufio"
	// "code.google.com/p/mahonia"
	"fmt"
	"io"
	"net"
	"os"
	"log"
	"strconv"
	"encoding/binary"
	"path/filepath"
	"time"
)


// TCPServer 开启服务器
func TCPServer(address string, port string) {
	ip := net.ParseIP(address)
	por, err := strconv.Atoi(port)
	addr := net.TCPAddr{ip, por, ""}
	listener, err := net.ListenTCP("tcp", &addr) //TCPListener listen
	if err != nil {
		log.Fatalf("Initialize error", err.Error())
	} else {
		result:="the server is listening\n"+"the address is "+address+" "+port+"\n"
		rs.SetText(result)
	}
	for {
		tcpcon, err := listener.AcceptTCP() //TCPConn client
		defer tcpcon.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
		result:=" tcp listen success\n"
		rs.Append(result)	
		tmp := make([]byte, 1)
		len, err := tcpcon.Read(tmp)
		if err != nil {
			log.Fatalf("read byte error: %v", err)
			return
		}
		if len != 1 {
			log.Fatalf("read byte len != 1")
			return
		}
		if(tmp[0]=='1'){
			lenBytes := make([]byte,2)
			readn, err := tcpcon.Read(lenBytes)
			if err != nil {
				log.Fatalf("read filename len error: %v", err)
			}
			if readn != 2 {
				log.Fatalf("readn != 2")
			}
			pathLen := binary.BigEndian.Uint16(lenBytes)
			path := make([]byte, pathLen, pathLen)
			readn, err = tcpcon.Read(path)
			if err != nil {
				log.Fatalf("read filename error: %v", err)
			}
			if readn != int(pathLen) {
				log.Fatalf("readn != filename len")
			}
			//接受文件名
		    filename :=filepath.Base(string(path))
			fi, err := os.Create(filename)
			defer fi.Close()
			if err != nil {
				log.Fatalf("file create error")
			}	
			//接受文件
			for {
				datas := make([]byte, 1024)
				wd, err := tcpcon.Read(datas)
				if err != nil {
					log.Fatalf("connection read error")
				}
				if string(datas[0:wd]) == "filerecvend" {
					break
				}
				_, err = fi.Write(datas[0:wd])
				if err != nil {
					log.Fatalf("file write error")
					break
				}
			}
			fiinfo, err := fi.Stat()
			recvbyte:=strconv.FormatInt(fiinfo.Size(),10)
			t := time.Now()
		    timess := t.Format("2006-01-02 15:04:05")
	        result=" in this upload :the time is "+timess+" the file name is "+filename+" and receive the "+recvbyte+" bytes\n"
			rs.Append(result)	
			continue;
		} else if(tmp[0]=='2'){
			//接受文件名
		    data := make([]byte, 1024)
			wc, err := tcpcon.Read(data)
			file:=string(data[0:wc])
			fi, err := os.Open(file)
			defer fi.Close()
			if err != nil {
				log.Fatalf("file open error")
			}
			//发送文件
			buff := make([]byte, 1024)
			for {
				n, err := fi.Read(buff)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if n == 0 {
					tcpcon.Write([]byte("filerecvend"))
					break
				}
				_, err = tcpcon.Write(buff)
				if err != nil {
					fmt.Println("write error")
				}
			}
			fiinfo, err := fi.Stat()
			sendbyte:=strconv.FormatInt(fiinfo.Size(),10)
			t := time.Now()
		    timess := t.Format("2006-01-02 15:04:05")
	        result=" in this dowmload :the time is "+timess+" the file name is "+string(data[0:wc])+" and send the "+sendbyte+" bytes\n"
			rs.Append(result)
			continue;
		}
	}
}
  
