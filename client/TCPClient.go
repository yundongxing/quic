package main

import (
	"strconv"
	"log"
	"io"
	"net"
	"os"
	"time"
	"encoding/binary"
//	"path/filepath"
//	"github.com/andlabs/ui"
//	_"github.com/andlabs/ui/winmanifest"
)


// TCPClient 开启客户端
<<<<<<< HEAD
func TCPClient(address string,files []string, test bool)  {
	if  test {
		for _,file :=range files{
			conn, err := net.Dial("tcp",address)
			defer conn.Close()
			if err != nil {
				log.Fatalf("connect server fail！")
			}
			t := time.Now()
			times := t.Format("2006-01-02 15:04:05")
			result:=" in this tcp upload :the start time is "+times+" the file name is "+file+"\n"
			urs.Append(result)
			fi, err := os.Open(file)
			defer fi.Close()
			if err != nil {
				panic(err)
			}
			fiinfo, err := fi.Stat()
			sendByte:= strconv.FormatInt(fiinfo.Size(),10)
			//tmp := make([]byte, 1)
			tmp:= []byte{'1'}
			_,err = conn.Write([]byte(tmp))
			if err != nil {
				log.Fatalf("write op error: %v", err)
			}
			//发送文件名
			pathLenBytes := make([]byte,2)
			binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
			writen, err := conn.Write(pathLenBytes)
			if err != nil {
				log.Fatalf("write filename len error: %v", err)
			}
			if writen != 2 {
				log.Fatalf("filename len != 2")
			}
			_, err = conn.Write([]byte(file))
			if err != nil {
				log.Fatalf("name send error")
			}
			//发送文件
			buff := make([]byte, 1024)
			for {
				n, err := fi.Read(buff)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if n == 0 {
					conn.Write([]byte("filerecvend"))
					break
				}
				_, err = conn.Write(buff)
				if err != nil {
					log.Fatalf(err.Error())
				}
			}
			tt := time.Now()
			timess := tt.Format("2006-01-02 15:04:05")
			result =" in this tcp upload :the end time is "+timess+" the file name is "+file+"\n"
			urs.Append(result)
			tsecond :=t.Unix()
			ttsecond :=tt.Unix()
			sub :=strconv.FormatInt(ttsecond-tsecond,10)
			result=" in this tcp upload :the address is"+address+"the time is:"+sub+" the file name is "+file+" and send the "+sendByte+" bytes\n"
			urs.Append(result)
		}	
	} else {
		for _,file :=range files{
			conn, err := net.Dial("tcp",address)
			defer conn.Close()
			if err != nil {
				log.Fatalf("connect server fail！")
			}
			t := time.Now()
			//times:=t.UTC().Format(time.UnixDate)
			times := t.Format("2006-01-02 15:04:05")
			result:=" in this download :the start time is "+times+" the file name is "+file+"\n"
			drs.Append(result)
			tmps:= []byte{'2'}
			_,err = conn.Write([]byte(tmps))
			if err != nil {
				log.Fatalf("write op error: %v", err)
			}
			//发送文件名
			_, err = conn.Write([] byte(file))
			if err != nil {
				log.Fatalf("name send error")
			}
			//接受文件
			fi, err := os.Create(file)
			defer fi.Close()
			if err != nil {
				log.Fatalf("file create error")
			}
			for {
				data := make([]byte, 1024)
				wd, err := conn.Read(data)
				if err != nil {
					log.Fatalf("connection read error")
				}
				if string(data[0:wd]) == "filerecvend" {
					break
				}
				_, err = fi.Write(data[0:wd])
				if err != nil {
					log.Fatalf("file write error")
					break
				}
			}	
			fiinfo, err := fi.Stat()
			recvbyte:=strconv.FormatInt(fiinfo.Size(),10)
			tt := time.Now()
			timess := tt.Format("2006-01-02 15:04:05")
			result =" in this download :the end time is "+timess+" the file name is "+file+"\n"
			drs.Append(result)
			tsecond :=t.Unix()
			ttsecond :=tt.Unix()
			sub :=strconv.FormatInt(ttsecond-tsecond,10)
			result=" in this tcp download :the time is: "+sub+" the file name is "+file+" and receive the "+recvbyte+" bytes\n"
			drs.Append(result)
		}
=======
func TCPClient(address,file string, test bool)  {
	if test {
		conn, err := net.Dial("tcp",address)
		defer conn.Close()
		if err != nil {
			log.Fatalf("connect server fail！")
		}
		t := time.Now()
		//times:=t.UTC().Format(time.UnixDate)
		times := t.Format("2006-01-02 15:04:05")
		result:=" in this tcp upload :the start time is "+times+" the file name is "+file+"\n"
		urs.Append(result)
		fi, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		fiinfo, err := fi.Stat()
		sendByte:= strconv.FormatInt(fiinfo.Size(),10)
		//tmp := make([]byte, 1)
		tmp:= []byte{'1'}
		_,err = conn.Write([]byte(tmp))
		if err != nil {
			log.Fatalf("write op error: %v", err)
		}
		//发送文件名
		pathLenBytes := make([]byte,2)
		binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
		writen, err := conn.Write(pathLenBytes)
		if err != nil {
			log.Fatalf("write filename len error: %v", err)
		}
		if writen != 2 {
			log.Fatalf("filename len != 2")
		}
		_, err = conn.Write([]byte(file))
		if err != nil {
			log.Fatalf("name send error")
		}
		//发送文件
		buff := make([]byte, 1024)
		for {
			n, err := fi.Read(buff)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				conn.Write([]byte("filerecvend"))
				break
			}
			_, err = conn.Write(buff)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
		fi.Close()
		tt := time.Now()
		timess := t.Format("2006-01-02 15:04:05")
		result =" in this tcp upload :the end time is "+timess+" the file name is "+file+"\n"
		urs.Append(result)
		tsecond :=t.Unix()
		ttsecond :=tt.Unix()
		sub :=strconv.FormatInt(ttsecond-tsecond,10)
		result=" in this tcp upload :the address is"+address+"the time is:"+sub+" the file name is "+file+" and send the "+sendByte+" bytes\n"
		urs.Append(result)	
	} else {
	//	tmp := make([]byte, 1, 1)
		//string tmp="2"
		conn, err := net.Dial("tcp",address)
		defer conn.Close()
		if err != nil {
			log.Fatalf("connect server fail！")
		}
		t := time.Now()
		//times:=t.UTC().Format(time.UnixDate)
		times := t.Format("2006-01-02 15:04:05")
		result:=" in this download :the start time is "+times+" the file name is "+file+"\n"
		drs.Append(result)
		tmps:= []byte{'2'}
		_,err = conn.Write([]byte(tmps))
		if err != nil {
			log.Fatalf("write op error: %v", err)
		}
		//发送文件名
		_, err = conn.Write([] byte(file))
		if err != nil {
			log.Fatalf("name send error")
		}
        //接受文件
		fi, err := os.Create(file)
		if err != nil {
		    log.Fatalf("file create error")
		}
		defer fi.Close()
		for {
			data := make([]byte, 1024)
			wd, err := conn.Read(data)
			if err != nil {
				log.Fatalf("connection read error")
			}
			if string(data[0:wd]) == "filerecvend" {
				break
			}
			_, err = fi.Write(data[0:wd])
			if err != nil {
				log.Fatalf("file write error")
				break
			}
		}	
		//fi, err := os.Open(file)
		fiinfo, err := fi.Stat()
		recvbyte:=strconv.FormatInt(fiinfo.Size(),10)
		//fi.Close()
		tt := time.Now()
		timess := t.Format("2006-01-02 15:04:05")
		result =" in this download :the end time is "+timess+" the file name is "+file+"\n"
		drs.Append(result)
		tsecond :=t.Unix()
		ttsecond :=tt.Unix()
		sub :=strconv.FormatInt(ttsecond-tsecond,10)
	    result=" in this tcp download :the time is: "+sub+" the file name is "+file+" and receive the "+recvbyte+" bytes\n"
		drs.Append(result)
>>>>>>> e2db2f5... a file transformission of quic
	}
	//ui.MsgBox(mainwin,
	//	"congratulation ,operation completed",
	//	"please first close this window")
}