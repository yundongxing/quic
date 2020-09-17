package main

import (
	"context"
	"encoding/binary"
	"math"
    "strings"
	//	"errors"
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
    "time"
	//	"github.com/873314461/quic-file/common"
	"github.com/lucas-clemente/quic-go"
)

var uploadfile [1000]filedir
var index = 0
var downfile [1000]filedir
var indexs = 0
var fileChunk =  1024*1024*50
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

type StreamHandler struct {
	Ctx    context.Context
	Stream quic.Stream
	Reader io.Reader
	Writer io.Writer
}

// 创建流对象
func NewStreamHandler(stream *quic.Stream) *StreamHandler {
	return &StreamHandler{
		Stream: *stream,
		Reader: io.Reader(*stream),
		Writer: io.Writer(*stream),
		Ctx:    context.Background(),
	}
}

// 流程序控制
func (h *StreamHandler) Run() {
	defer h.Stream.Close()
	tmp := make([]byte, 1, 1)
	len, err := h.Reader.Read(tmp)
	if err != nil {
		log.Fatalf("read byte error: %v", err)
		return
	}
	if len != 1 {
		log.Fatalf("read byte len != 1")
		return
	}
	ops := int(tmp[0])
	result := " quic listen success\n"
	rs.Append(result)
	//log.Printf(string(tmp))
	//if tmp[0] == 'c' {
	//	fmt.Printf("ops=3")
	//h.ChunkUpload()
	//}
	if ops == 1 {
	//	log.Printf("the ops is %d\n", ops)
		h.handlerUpload()
	} else if ops == 2 {
	//	log.Printf("the ops is %d\n", ops)
		h.handlerDownload()
	} else if ops == 3 {
	//	log.Printf("the ops is %d\n", ops)
		//h.ChunkUpload()
		h.ChunkUploaddata()
	} else if ops == 4 {
		//h.ChunkDownload()
	//	log.Printf("the ops is %d\n", ops)
		h.ChunkDownloaddata()
	}
	//defer h.Stream.Close()
}

// 文件上传
func (h *StreamHandler) handlerUpload() {
	lenBytes := make([]byte, 2, 2)
	readn, err := h.Reader.Read(lenBytes)
	if err != nil {
		log.Fatalf("read filename len error: %v", err)
	}
	if readn != 2 {
		log.Fatalf("readn != 2")
	}
	pathLen := binary.BigEndian.Uint16(lenBytes)
	path := make([]byte, pathLen, pathLen)
	readn, err = h.Reader.Read(path)
	if err != nil {
		log.Fatalf("read filename error: %v", err)
	}
	if readn != int(pathLen) {
		log.Fatalf("readn != filename len")
	}
	filename := filepath.Base(string(path))
	tmpAbsPath, err := filepath.Abs(filename + TempFileSuffix)
	if err != nil {
		log.Fatalf("get tmp abs path error: %v", err)
	}
	absPath, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("get abs path error: %v", err)
	}
	dataLenBytes := make([]byte, 8, 8)
	readn, err = h.Reader.Read(dataLenBytes)
	if err != nil {
		log.Fatalf("read file len error: %v", err)
	}
	if readn != 8 {
		log.Fatalf("readn != 8")
	}
	size := binary.BigEndian.Uint64(dataLenBytes)
	if int(size)<fileChunk {
		file, err := os.Create(tmpAbsPath)
		if err != nil {
			log.Fatalf("creat file error: %v", err)
		}
		writen, err := io.Copy(file, h.Reader)
		if err != nil {
			log.Fatalf("write file error: %v", err)
		}
		if size != uint64(writen) {
			log.Fatalf("data len != writen")
		}
		file.Close()
		err = os.Rename(tmpAbsPath, absPath)
		if err != nil {
			log.Fatalf("rename file error: %v", err)
		}
		writebyte := strconv.FormatInt(writen, 10)
		result := " in this upload :the file name is " + filename + " and receive the " + writebyte + " bytes\n"
		rs.Append(result)
	} else {
		uploadfile[index] = Newfiledir(filename, size)
		index++
		//defer h.Stream.Close()
		pathlen := make([]byte, 2)
		binary.BigEndian.PutUint16(pathlen, uint16(999))
		writen, err := h.Writer.Write(pathlen)
		if err != nil {
			log.Fatalf("write end error")
		}
		if writen != 2 {
			log.Fatalf("filename len != 2")
		}
	}
}

//文件分块数据上传
func (h *StreamHandler) ChunkUploaddata() {
	lenBytes := make([]byte, 2, 2)
	readn, err := h.Reader.Read(lenBytes)
	if err != nil {
		log.Fatalf("read filename len error: %v", err)
	}
	if readn != 2 {
		log.Fatalf("readn != 2")
	}
	pathLen := binary.BigEndian.Uint16(lenBytes)
	path := make([]byte, pathLen, pathLen)
	readn, err = h.Reader.Read(path)
	if err != nil {
		log.Fatalf("read filename error: %v", err)
	}
	if readn != int(pathLen) {
		log.Fatalf("readn != filename len")
	}
	filename := filepath.Base(string(path))
	num := make([]byte, 2, 2)
	readn, err = h.Reader.Read(num)
	if err != nil {
		log.Fatalf("read filename error: %v", err)
	}
	if readn != 2 {
		log.Fatalf("readn != 2")
	}
	number := (int)(binary.BigEndian.Uint16(num))
	lens := make([]byte, 8)
	readn, err = h.Reader.Read(lens)
	if err != nil {
		log.Fatalf("read file size len error: %v", err)
	}
	if readn != 8 {
		log.Fatalf("readn != 2")
	}
	datalen := binary.BigEndian.Uint64(lens)
//	con := make([]byte,datalen)
	log.Printf("the size is %d\n",datalen)
	tmpAbsPath:=filename+strconv.FormatInt(int64(number),10)
	file, err :=os.Create(tmpAbsPath)
	defer file.Close()
	if err != nil {
		log.Fatalf("creat file error: %v", err)
	}
	writen, err :=io.Copy(file,h.Reader)
	if err != nil {
		log.Fatalf("read file error: %v", err)
	}
	if writen != int64(datalen) {
		log.Fatalf("write data error")
	}
	for i, _ := range uploadfile {
		if filename == uploadfile[i].filenames {
			file ,err :=os.OpenFile(tmpAbsPath, os.O_RDONLY, 0666)
			if err != nil {
				log.Println("open file error :", err)
				return
			}
			con := make([]byte,datalen)
			file.Read(con)
			file.Close()
			uploadfile[i].filedata[number] = string(con)
			numbers := int(math.Ceil(float64(uploadfile[i].filesize) / float64(fileChunk)))
			if len(uploadfile[i].filedata) == numbers {
				//log.Printf("write over\n")
				f, err := os.OpenFile(uploadfile[i].filenames, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					log.Println("open file error :", err)
					return
				}
				defer f.Close()
				writes := 0
				write := bufio.NewWriter(f)
				for l := 1; l <= len(uploadfile[i].filedata); l++ {
					//writen, err := tmpFile.WriteString(uploadfile.filedata[l])
					writen, err := write.WriteString(uploadfile[i].filedata[l])
					if err != nil {
						log.Println(err)
						return
					}
					writes += writen
				}
				write.Flush()
				if writes == int(uploadfile[i].filesize) {
					//time.Sleep(100*time.Second)
					for j:=1;j<=numbers;j++{
						if j!=number{
                            tmpPath:=uploadfile[i].filenames+strconv.FormatInt(int64(j),10)
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
				result := " in this upload :the end time is " + timess + " the file name is " + filename + "\n"
				rs.Append(result)
				}
			}
		}
	}
	
}

//文件下载
func (h *StreamHandler) handlerDownload() {
	lenBytes := make([]byte, 2, 2)
	readn, err := h.Reader.Read(lenBytes)
	if err != nil {
		log.Fatalf("read filename len error: %v", err)
	}
	if readn != 2 {
		log.Fatalf("readn != 2")
	}
	pathLen := binary.BigEndian.Uint16(lenBytes)
	path := make([]byte, pathLen, pathLen)
	readn, err = h.Reader.Read(path)
	if err != nil {
		log.Fatalf("read filename error: %v", err)
	}
	if readn != int(pathLen) {
		log.Fatalf("readn != path len")
	}
	file, err := os.Open(string(path))
	if err != nil {
		log.Fatalf("open file[%s] error: %v", string(path), err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("get file[%s] info error: %v", string(path), err)
	}
	filelens := make([]byte, 8, 8)
	size := uint64(fileInfo.Size())
	binary.BigEndian.PutUint64(filelens, size)
	_, err = h.Writer.Write(filelens)
	if err != nil {
		log.Fatalf("file size write error")
	}
	if int(size) < fileChunk {
        sendN, err := io.Copy(h.Writer, file)
		if err != nil {
			log.Fatalf("send file[%s] error: %v", string(path), err)
		}
		if sendN != fileInfo.Size() {
			log.Fatalf("sendn != file size")
		}
		sendbyte := strconv.FormatInt(sendN, 10)
		result := " in this download : the file name is " + string(path) + " and send the " + sendbyte + " bytes\n"
		rs.Append(result)
	}else {
		filename := filepath.Base(string(path))
		downfile[indexs] = Newfiledir(filename, size)
		number := int(math.Ceil(float64(size) / float64(fileChunk)))
		for i := 1; i <= number; i++ {
			partSize := int(math.Min(float64(fileChunk), float64(size-uint64((i-1)*fileChunk))))
			partBuffer := make([]byte, partSize)
			file.Read(partBuffer)
			downfile[indexs].filedata[i] = string(partBuffer)
		}
		pathlen := make([]byte, 2)
		binary.BigEndian.PutUint16(pathlen, uint16(999))
		writen, err := h.Writer.Write(pathlen)
		if err != nil {
			log.Fatalf("write end error")
		}
		if writen != 2 {
			log.Fatalf("filename len != 2")
		}
	}
}

//文件分块数据下载
func (h *StreamHandler) ChunkDownloaddata() {
	lenBytes := make([]byte, 2, 2)
	readn, err := h.Reader.Read(lenBytes)
	if err != nil {
		log.Fatalf("read filename len error: %v", err)
	}
	if readn != 2 {
		log.Fatalf("readn != 2")
	}
	pathLen := binary.BigEndian.Uint16(lenBytes)
	path := make([]byte, pathLen, pathLen)
	readn, err = h.Reader.Read(path)
	if err != nil {
		log.Fatalf("read filename error: %v", err)
	}
	if readn != int(pathLen) {
		log.Fatalf("readn != path len")
	}
	filename := filepath.Base(string(path))
	nums := make([]byte, 2, 2)
	readn, err = h.Reader.Read(nums)
	if err != nil {
		log.Fatalf("read filename len error: %v", err)
	}
	if readn != 2 {
		log.Fatalf("readn != 2")
	}
	num := (int)(binary.BigEndian.Uint16(nums))
	datalen := make([]byte, 8)
	readn, err = h.Reader.Read(datalen)
	dataLen := binary.BigEndian.Uint64(datalen)
	//log.Printf("the len1 is %d\n", datalen)
	data := make([]byte, dataLen, dataLen)
	var numbers int
	for i, _ := range downfile {
		if filename == downfile[i].filenames {
			data = ([]byte)(downfile[i].filedata[num])
			numbers = int(math.Ceil(float64(downfile[i].filesize) / float64(fileChunk)))
		}
	}
    s := strings.NewReader(string(data))
    w := bufio.NewReader(s)
	//writen,err =writer.ReadFrom(s)
	w.WriteTo(h.Writer)	
	h.Stream.Close()
	if num==numbers{
        tt := time.Now()
		//timess:=tt.UTC().Format(time.UnixDate)
		timess := tt.Format("2006-01-02 15:04:05")
		result := " in this upload :the end time is " + timess + " the file name is " + filename + "\n"
		rs.Append(result)
	}
}
