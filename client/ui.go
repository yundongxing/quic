package main

import (
    "github.com/andlabs/ui"
	_"github.com/andlabs/ui/winmanifest"
	"io/ioutil"
//	"path/filepath"
	"strings"

)

var files []string //传输文件列表
var address string //传输IP地址
var protocol,test bool  //判别协议
var mainwin *ui.Window  //主界面
var urs,drs *ui.MultilineEntry  //上传结果和下载结果
var ip,port *ui.Entry  //传输IP和端口
var dbox *ui.EditableCombobox  //传输结果框

func makeBasicControlsPage() ui.Control{	
	//整体模块
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
	group := ui.NewGroup("choices")
	group.SetMargined(true)
	box.Append(group, true)
	lbox :=ui.NewVerticalBox()
	lbox.SetPadded(true)
	group.SetChild(lbox)
	
    //第一部分
	lbox.Append(ui.NewLabel("Server Ip"),false)
	ip=ui.NewEntry()
	lbox.Append(ip, false)
	port=ui.NewEntry()
	lbox.Append(ui.NewLabel("Server Port"),false)
	lbox.Append(port, false)
	
    //第二部分
	lbox.Append(ui.NewLabel("Please choose the protocol "),false)
	dbox = ui.NewEditableCombobox()
	lbox.Append(dbox, false)
	dbox.Append("QUIC Protocol")
	dbox.Append("TCP  Protocol")
	//第三部分
	sbutton :=ui.NewButton("refer file")
	grid := ui.NewGrid()
	grid.SetPadded(true)
	lbox.Append(grid, false)
	grid.Append(sbutton,0, 0, 1, 1,true, ui.AlignFill, false, ui.AlignFill)
	entry := ui.NewNonWrappingMultilineEntry()
	entry.Handle()
	lbox.Append(entry,true)
	sbutton.OnClicked(func(*ui.Button){
		file:="the filename in current directory\n"
		filess, _ := ioutil.ReadDir("./")
		for _,f :=range filess {
			file+=f.Name()+"\n"
		}
		entry.SetText(file)
	})
	return box
}
// 上传界面
func makeUploadPage() ui.Control{	
	var c *FileClient
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
	group:=ui.NewGroup("upload")
	group.SetMargined(true)
	box.Append(group, true)
	rbox := ui.NewVerticalBox()
	rbox.SetPadded(true)
	group.SetChild(rbox)
	
	//第一部分
	uentry := ui.NewMultilineEntry()
	rbox.Append(uentry, true)
    //第二部分
	grid := ui.NewGrid()
	grid.SetPadded(true)
	rbox.Append(grid, false)
	button := ui.NewButton("  select  ")
	sbutton := ui.NewButton(" upload  ")
	cbutton := ui.NewButton("  clear  ")
	// 选择文件按钮
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "you have not choose file"
		}
		uentry.Append(string(filename)+"\n")
		files=append(files,filename)
	})
	// 上传按钮
	sbutton.OnClicked(func(*ui.Button){
		address =ip.Text()+":"+port.Text()
		if dbox.Text()=="QUIC Protocol"{
			protocol=true
		}else{
			protocol=false
		}
		if protocol  {
			c=NewFileClient(address)
			count:=len(files)
			chquit := make( chan int,count)
			for _,file :=range files{
				//fileReader, size := ReadFile(file)
	           // defer fileReader.Close()
				//if size<1024*1024*1 {
                go c.Upload(file,chquit)
				//}else{
				//	go c.ChunkUpload(file)
				//}	   
			}
			/*
			go func(){
				for range(chquit){
					<-chquit
					count=count-1;
					if(count==0){
					   close(chquit)
					   c.Close()
					}
				}
			}()
			*/
		} else {
			//for _,file :=range files{
				go TCPClient(address,files,true)
			//}		
		}
	})
    // 清空按钮
	cbutton.OnClicked(func(*ui.Button) {
		files=files[0:0]
        uentry.SetText("")
	})
	grid.Append(button,0,0,1,1,true,1,true,1)
	grid.Append(sbutton,1,0,1,1,true,ui.AlignCenter,true,1)
	grid.Append(cbutton,2,0,1,1,true,3,true,ui.AlignFill)
    //第三部分
	urs=ui.NewNonWrappingMultilineEntry()
	urs.Handle()
	urs.SetText("")
	rbox.Append(ui.NewLabel("the file result" ),false)
	rbox.Append(urs, true)
    return box
}

func makeDownloadPage() ui.Control{
	var c *FileClient
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
	group:=ui.NewGroup("download")
	group.SetMargined(true)
	box.Append(group, true)
	rbox := ui.NewVerticalBox()
	rbox.SetPadded(true)
	group.SetChild(rbox)
	//第一部分
	uentry := ui.NewMultilineEntry()
	rbox.Append(uentry, true)
    //第二部分
	grid := ui.NewGrid()
	grid.SetPadded(true)
	rbox.Append(grid, false)
	dbutton := ui.NewButton("  download  ")
	cbutton := ui.NewButton("  clear  ")
	// 下载按钮
	dbutton.OnClicked(func(*ui.Button) {
		address =ip.Text()+":"+port.Text()
		if dbox.Text()=="QUIC Protocol"{
			protocol=true
		}else{
			protocol=false
		}
		if protocol  {
			file:=uentry.Text()
			files=strings.Split(file,",")
			c=NewFileClient(address)
			count:=len(files)
			chquit := make( chan int,count)
			for _,file :=range files{
				if(file!=""){
					//fileReader, size := ReadFile(file)
	            	//defer fileReader.Close()
				//	if size<1024*1024*1{
				    go c.Download(file,chquit)
				//	}else{
				//		go c.ChunkDownload(file,chquit)
				}	   
					//go c.Download(file,chquit)
			}	
			/*
			go func(){
				for range(chquit){
					<-chquit
					count=count-1;
					if(count==0){
					   close(chquit)
					   c.Close()
					}
				}
			}()*/
		} else {
			filee:=uentry.Text()
			files=strings.Split(filee,",")
			//for _,file :=range files{
			go TCPClient(address,files,false)
			//}		
		}
	})
	// 清空按钮
	cbutton.OnClicked(func(*ui.Button) {
		files=files[0:0]
		uentry.SetText("")
	})
	grid.Append(dbutton,0,0,1,1,true,2,true, 1)
	grid.Append(cbutton,1,0,1,1,true,5,true,1)
    //第三部分
	drs=ui.NewNonWrappingMultilineEntry()
	drs.Handle()
	drs.SetText("")
	rbox.Append(ui.NewLabel("the file result" ),false)
	rbox.Append(drs, true)
	return box
}

func setupUI(){
	mainwin=ui.NewWindow("this is a quic transmission for client",440,500,true)

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("client controls", makeBasicControlsPage())
	tab.SetMargined(0, true)

	tab.Append("uploads pages", makeUploadPage())
	tab.SetMargined(1, true)

	tab.Append("download pages", makeDownloadPage())
	tab.SetMargined(2, true)
	mainwin.Show()
	mainwin.OnClosing(func(*ui.Window) bool{
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool{
		mainwin.Destroy()
		return true
	})
}

func main(){
	ui.Main(setupUI)
}

