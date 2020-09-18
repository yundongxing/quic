package main

import (
    "github.com/andlabs/ui"
	_"github.com/andlabs/ui/winmanifest"
	"io/ioutil"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
)

var mainwin *ui.Window
var rs,entry *ui.MultilineEntry
var s  *FileServer
func makeBasicControlsPage() ui.Control{
	//整体模块
    var protocol  bool
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
    //左边整体布局
	group := ui.NewGroup("choices")
	group.SetMargined(true)
	box.Append(group, true)
	lbox :=ui.NewVerticalBox()
	lbox.SetPadded(true)
	group.SetChild(lbox)
    //左边第一部分
	lbox.Append(ui.NewLabel("Server Ip"),false)
	ip :=ui.NewEntry()
	ip.SetReadOnly(false)
	port :=ui.NewEntry()
	lbox.Append(ip, false)
	lbox.Append(ui.NewLabel("Server Port"),false)
	lbox.Append(port, false)
    //左边第二部分
	lbox.Append(ui.NewLabel("Please choose the protocol "),false)
	cbox := ui.NewEditableCombobox()
	lbox.Append(cbox, false)
	cbox.Append("QUIC Protocol")
	cbox.Append("TCP Protocol")
	//左边第三部分
	entry =ui.NewNonWrappingMultilineEntry()
	entry.Handle()
	lbox.Append(ui.NewLabel("the filename in current directory" ),false)
	lbox.Append(entry, true)
    //右边整体部局
	rgroup:=ui.NewGroup("results")
	rgroup.SetMargined(true)
	box.Append(rgroup, true)
	rbox := ui.NewVerticalBox()
	rbox.SetPadded(true)
	rgroup.SetChild(rbox)
	//右边第一部分
	button :=ui.NewButton("start server")
	button2 :=ui.NewButton("refer file")
	grid := ui.NewGrid()
	grid.SetPadded(true)
	rbox.Append(grid, false)
	grid.Append(button,0, 0, 1, 1,true, ui.AlignFill, true, ui.AlignFill)
	grid.Append(button2,1, 0, 1, 1,true, ui.AlignFill, true, ui.AlignFill)
	rbox.Append(ui.NewLabel("the file transmission result" ),false)
	rs = ui.NewNonWrappingMultilineEntry()
	rs.Handle()
	rbox.Append(rs,true)
	//监听按钮
	button.OnClicked(func(*ui.Button){
		//判断协议
		if (cbox.Text()=="QUIC Protocol"){
			protocol=true
		} else {
			protocol=false
		}
		//判断上传下载	
		if protocol{
			address :=ip.Text()+":"+port.Text()
			s=NewFileServer(address,generateTLSConfig(), nil)
		    go s.Run()
		} else {
		    go TCPServer(ip.Text(),port.Text())
		}
	})
	//查看目录文件按钮
	button2.OnClicked(func(*ui.Button) {
		file:=""
		files, _ := ioutil.ReadDir("./")
		for _, f := range files {
			file+=f.Name()+"\n"
		}
		entry.Append(file+"\n")
	})
    return box
}

func setupUI(){
	mainwin=ui.NewWindow("this is a quic transmission for server",640,480,true)
	
	mainwin.OnClosing(func(*ui.Window) bool{
		ui.Quit()
		return true
	})
<<<<<<< HEAD
	//ui.OnShouldQuit(func() bool{
	//	mainwin.Destroy()
	//	return true
	//})
=======
	ui.OnShouldQuit(func() bool{
		mainwin.Destroy()
		return true
	})
>>>>>>> e2db2f5... a file transformission of quic
	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("server controls", makeBasicControlsPage())
	tab.SetMargined(0, true)

	mainwin.Show()
}

func main(){
	ui.Main(setupUI)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
