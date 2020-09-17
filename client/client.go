package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"io"
	"log"
	"os"
	"strings"

	//	"path/filepath"
	"math"
	"strconv"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

type FileClient struct {
	Session quic.Session
	Ctx     context.Context
}

// 创建文件对象
func NewFileClient(address string) *FileClient {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(address, tlsConf, nil)
	if err != nil {
		log.Fatalf("connect server error: %v\n", err)
	}
	return &FileClient{
		Session: session,
		Ctx:     context.Background(),
	}
}

// 关闭会话
func (c *FileClient) Close() {
	time.Sleep(time.Second)
	time.Sleep(time.Second)
}

// 文件结构
type filedir struct {
	filenames string
	filesize  uint64
	filedata  map[int]string
}

func Newfiledir(filename string, size uint64) filedir {
	dada := make(map[int]string)
	return filedir{
		filenames: filename,
		filesize:  size,
		filedata:  dada,
	}
}

var clientfile [100]filedir
var cindex = 0
var fileChunk = 1024 * 1024 * 50

// 上传文件
func (c *FileClient) Upload(file string, quit chan int) {
	fileReader, size := ReadFile(file)
	defer fileReader.Close()
	stream, err := c.Session.OpenStreamSync(c.Ctx)
	if err != nil {
		log.Fatalf("open stream error: %v\n", err)
	}
	defer stream.Close()
	t := time.Now()
	times := t.Format("2006-01-02 15:04:05")
	result := " in this upload :the start time is " + times + " the file name is " + file + "\n"
	urs.Append(result)
	writer := bufio.NewWriter(stream)
	err = writer.WriteByte(byte(1))
	if err != nil {
		log.Fatalf("write op error: %v", err)
	}
	pathLenBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
	writen, err := writer.Write(pathLenBytes)
	if err != nil {
		log.Fatalf("write filename len error: %v", err)
	}
	if writen != 2 {
		log.Fatalf("filename len != 2")
	}
	writen, err = writer.WriteString(file)
	if err != nil {
		log.Fatalf("write filename error: %v", err)
	}
	if writen != len(file) {
		log.Fatalf("writen !=filename , %d, %d", writen, len(file))
	}
	dataLenBytes := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(dataLenBytes, size)
	writen, err = writer.Write(dataLenBytes)
	if err != nil {
		log.Fatalf("write file len error: %v", err)
	}
	if writen != 8 {
		log.Fatalf("file len != 8")
	}
	if int(size) < fileChunk {
		writeFileN, err := writer.ReadFrom(fileReader)
		if err != nil {
			log.Fatalf("write data error: %v", err)
		}
		if uint64(writeFileN) != size {
			log.Fatalf("write file n != file size")
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalf("writer flush error: %v", err)
		}
		sendbyte := strconv.FormatInt(writeFileN, 10)
		tt := time.Now()
		timess := tt.Format("2006-01-02 15:04:05")
		result = " in this upload :the end time is " + timess + " the file name is " + file + "\n"
		urs.Append(result)
		tsecond := t.Unix()
		ttsecond := tt.Unix()
		sub := strconv.FormatInt(ttsecond-tsecond, 10)
		result = " in this upload :the time is " + sub + " the file name is " + file + "and send the " + sendbyte + " bytes\n"
		urs.Append(result)
		quit <- 1
	} else {
		err = writer.Flush()
		end := make([]byte, 2)
		reader := bufio.NewReader(stream)
		reader.Read(end)
		pn := binary.BigEndian.Uint16(end)
		//log.Printf("the end is %d\n", pn)
		if pn == 999 {
			number := int(math.Ceil(float64(size) / float64(fileChunk)))
			for i := 1; i <= number; i++ {
				partSize := int(math.Ceil(math.Min(float64(fileChunk), float64(size-uint64((i-1)*fileChunk)))))
				partBuffer := make([]byte, partSize)
				fileReader.Read(partBuffer)
				//fmt.Printf(string(partBuffer))
				go c.ChunkUploaddata(file, i, partSize, partBuffer, quit)
			}
		}
	}

}

// 分快上传数据文件
func (c *FileClient) ChunkUploaddata(file string, i int, size int, concent []byte, quit chan int) {
	stream, err := c.Session.OpenStreamSync(c.Ctx)
	if err != nil {
		log.Fatalf("open stream error: %v\n", err)
	}
	writer := bufio.NewWriter(stream)
	err = writer.WriteByte(byte(3))
	pathLenBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
	writen, err := writer.Write(pathLenBytes)
	if err != nil {
		log.Fatalf("write filename len error: %v", err)
	}
	if writen != 2 {
		log.Fatalf("filename len != 2")
	}
	writen, err = writer.WriteString(file)
	if err != nil {
		log.Fatalf("write filename error: %v", err)
	}
	if writen != len(file) {
		log.Fatalf("writen !=filename , %d, %d", writen, len(file))
	}
	num := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(num, (uint16)(i))
	writen, err = writer.Write(num)
	if err != nil {
		log.Fatalf("number chunk transmission error\n")
	}
	if writen != 2 {
		log.Fatalf("file len != 2")
	}
	dataLenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLenBytes, (uint64)(size))
	writen, err = writer.Write(dataLenBytes)
	if writen != 8 {
		log.Fatalf("file len != 8")
	}
	s := strings.NewReader(string(concent))
	w := bufio.NewReader(s)
	//writen,err =writer.ReadFrom(s)
	w.WriteTo(writer)
	if err != nil {
		log.Fatalf("write data errors: %v", err)
	}
	// if writen!=size {
	//	log.Fatalf("write dada error\n")
	//}
	err = writer.Flush()
	defer stream.Close()
	//log.Printf("the data is "+string(data)+"\n\n\n")
	number := int(math.Ceil(float64(size) / float64(fileChunk)))
	if i == number {
		tt := time.Now()
		//timess:=tt.UTC().Format(time.UnixDate)
		timess := tt.Format("2006-01-02 15:04:05")
		result := " in this upload :the end time is " + timess + " the file name is " + file + "\n"
		urs.Append(result)
		quit <- 1
	}
}

// 读取文件
func ReadFile(file string) (*os.File, uint64) {
	fp, err := os.Open(file)
	if err != nil {
		log.Fatalf("open file error: %v\n", err)
	}
	fileInfo, err := fp.Stat()
	if err != nil {
		log.Fatalf("get file info error: %v\n", err)
	}
	return fp, uint64(fileInfo.Size())
}

// 下载文件
func (c *FileClient) Download(file string, quit chan int) {
	stream, err := c.Session.OpenStreamSync(c.Ctx)
	if err != nil {
		log.Fatalf("open stream error: %v", err)
	}
	defer stream.Close()
	t := time.Now()
	//times:=t.UTC().Format(time.UnixDate)
	times := t.Format("2006-01-02 15:04:05")
	result := " in this download :the start time is " + times + " the file name is " + file + "\n"
	drs.Append(result)
	writer := bufio.NewWriter(stream)
	err = writer.WriteByte(byte(2))
	if err != nil {
		log.Fatalf("write op error: %v", err)
	}
	pathLenBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
	writeN, err := writer.Write(pathLenBytes)
	if err != nil {
		log.Fatalf("write filename len error: %v", err)
	}
	if writeN != 2 {
		log.Fatalf("filename len != 2")
	}
	writeN, err = writer.WriteString(file)
	if err != nil {
		log.Fatalf("write filename error: %v", err)
	}
	if writeN != len(file) {
		log.Fatalf("writeN != filename len, %d, %d", writeN, len(file))
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("writer flush error: %v", err)
	}
	filelens := make([]byte, 8, 8)
	reader := bufio.NewReader(stream)
	reader.Read(filelens)
	if err != nil {
		log.Fatalf("read file size error: %v", err)
	}
	filelen := binary.BigEndian.Uint64(filelens)
	if int(filelen) < fileChunk {
		tmpFile, err := os.Create(file)
		defer tmpFile.Close()
		if err != nil {
			log.Fatalf("creat file error: %v", err)
		}
		recvN, err := io.Copy(tmpFile, stream)
		if err != nil {
			log.Fatalf("write file error: %v", err)
		}
		recvbyte := strconv.FormatInt(recvN, 10)
		tt := time.Now()
		timess := tt.Format("2006-01-02 15:04:05")
		result = " in this download :the end time is " + timess + " the file name is " + file + "\n"
		drs.Append(result)
		tsecond := t.Unix()
		ttsecond := tt.Unix()
		sub := strconv.FormatInt(ttsecond-tsecond, 10)
		result = " in this download :the time is " + sub + " the file name is " + file + " and receive the " + recvbyte + " bytes\n"
		drs.Append(result)
		quit <- 1
	} else {
		clientfile[cindex] = Newfiledir(file, filelen)
		cindex++
		number := int(math.Ceil(float64(filelen) / float64(fileChunk)))
		end := make([]byte, 2)
		reader.Read(end)
		pn := binary.BigEndian.Uint16(end)
		log.Printf("the end is %d\n", pn)
		if pn == 999 {
			for i := 1; i <= number; i++ {
				partSize := int(math.Min(float64(fileChunk), float64((filelen)-uint64((i-1)*fileChunk))))
				go c.ChunkDownloaddata(file, filelen, i, partSize, quit)
				//log.Printf("the size is %d\n",partSize)
			}
		}
	}
}

// 分快下载传输文件
func (c *FileClient) ChunkDownloaddata(file string, size uint64, i int, datalen int, quit chan int) {
	stream, err := c.Session.OpenStreamSync(c.Ctx)
	if err != nil {
		log.Fatalf("open stream error: %v", err)
	}
	defer stream.Close()
	writer := bufio.NewWriter(stream)
	err = writer.WriteByte(byte(4))
	pathLenBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
	writeN, err := writer.Write(pathLenBytes)
	if err != nil {
		log.Fatalf("write filename len error: %v", err)
	}
	if writeN != 2 {
		log.Fatalf("filename len != 2")
	}
	writeN, err = writer.WriteString(file)
	if err != nil {
		log.Fatalf("write filename error: %v", err)
	}
	if writeN != len(file) {
		log.Fatalf("writeN != filename len, %d, %d", writeN, len(file))
	}
	num := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(num, (uint16)(i))
	writen, err := writer.Write(num)
	if err != nil {
		log.Fatalf("number chunk transmission error\n")
	}
	if writen != 2 {
		log.Fatalf("file len != 2")
	}
	datalens := make([]byte, 8)
	binary.BigEndian.PutUint64(datalens, (uint64)(datalen))
	writen, err = writer.Write(datalens)
	if err != nil {
		log.Fatalf("number chunk transmission error\n")
	}
	if writen != 8 {
		log.Fatalf("file len != 2")
	}
	err = writer.Flush()
	data := make([]byte, datalen, datalen)
	//reader := bufio.NewReader(stream)
	tmpAbsPath := file + strconv.FormatInt(int64(i), 10)
	//reader.Read(data)
	f, err := os.Create(tmpAbsPath)
	defer f.Close()
	if err != nil {
		log.Fatalf("creat file error: %v", err)
	}
	writenn, err := io.Copy(f, stream)
	if err != nil {
		log.Fatalf("read file error: %v", err)
	}
	if writenn != int64(datalen) {
		log.Fatalf("write data error")
	}
	for j, _ := range clientfile {
		if file == clientfile[j].filenames {
			f, err := os.OpenFile(tmpAbsPath, os.O_RDONLY, 0666)
			if err != nil {
				log.Println("open file error :", err)
				return
			}
			f.Read(data)
			defer f.Close()
			clientfile[j].filedata[i] = string(data)
			//log.Printf(string(data))
			numbers := int(math.Ceil(float64(clientfile[j].filesize) / float64(fileChunk)))
			if len(clientfile[j].filedata) == numbers {
				log.Printf("the new data\n")
				f, err := os.Create(file)
				if err != nil {
					log.Println("create file error :", err)
					return
				}
				defer f.Close()
				var writes = 0
				for l := 1; l <= len(clientfile[j].filedata); l++ {
					//writen, err := tmpFile.WriteString(uploadfile.filedata[l])
					writen, err := f.WriteString(clientfile[j].filedata[l])
					if err != nil {
						log.Println(err)
						return
					}
					writes += writen
				}
				if writes == int(clientfile[j].filesize) {
					//time.Sleep(100*time.Second)
					for l := 1; l <= numbers; l++ {
						if l != i {
							tmpPath := clientfile[j].filenames + strconv.FormatInt(int64(l), 10)
							err = os.Remove(tmpPath)
							if err != nil {
								log.Println("delete file error :", err)
							} else {
								log.Println("delete file success\n")
							}
						}
					}
					tt := time.Now()
					//timess:=tt.UTC().Format(time.UnixDate)
					timess := tt.Format("2006-01-02 15:04:05")
					result := " in this upload :the end time is " + timess + " the file name is " + file + "\n"
					urs.Append(result)
					quit <- 1
				}
			}
		}
	}
}
