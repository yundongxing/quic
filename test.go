package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/lucas-clemente/quic-go"
)

var (
	address = "127.0.0.1:8000"
)

func Client() {
	log.Println("client connect to", address)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-file"},
	}
	sess, err := quic.DialAddr(address, tlsConf, nil)
	if err != nil {
		log.Fatalln("DialAddr error:", err)
	}
	st, err := sess.OpenStreamSync()
	if err != nil {
		log.Fatalln("OpenStreamSync error:", err)
	}
	log.Println("Get Stream", st)
	stw := bufio.NewWriter(st)
	stw.WriteString("Hello\n")
	stw.Flush()
	str := bufio.NewReader(st)
	data, err := str.ReadString(byte('\n'))
	log.Println("echo", data)
	st1, err := sess.AcceptStream()
	if err != nil {
		log.Fatalln("AcceptStream error:", err)
	}
	log.Println("Get Stream", st1)
	n, err := io.CopyN(st1, st1, 6)
	if err != nil {
		log.Fatalln("Copy error:", err)
	}
	log.Println("echo", n, "bytes")
	time.Sleep(time.Second)
	st.Close()
	st1.Close()
	sess.Close()

}

func Server() {
	log.Println("server is listening at", address)
	l, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err != nil {
		log.Fatalln("ListenAddr error:", err)
	}
	sess, err := l.Accept()
	if err != nil {
		log.Fatal("Accept error", err)
	}
	st, err := sess.AcceptStream()
	if err != nil {
		log.Fatalln("AcceptStream error:", err)
	}
	log.Println("Get Stream", st)
	n, err := io.CopyN(st, st, 6)
	if err != nil {
		log.Fatalln("Copy error:", err)
	}
	log.Println("echo", n, "bytes")
	st1, err := sess.OpenStreamSync()
	if err != nil {
		log.Fatal("OpenStreamSync error", err)
	}
	log.Println("Get Stream", st1)
	stw := bufio.NewWriter(st1)
	stw.WriteString("Hello\n")
	stw.Flush()
	str := bufio.NewReader(st1)
	data, err := str.ReadString(byte('\n'))
	log.Println("echo", data)
	time.Sleep(time.Second)
	st.Close()
	st1.Close()
	sess.Close()
}

func main() {
	isServer := flag.Bool("s", false, "server mode")
	isClient := flag.Bool("c", false, "client mode")
	flag.Parse()

	if (*isServer && *isClient) || (!*isServer && !*isClient) {
		log.Fatalln("server or client?")
	}
	if *isServer {
		Server()
	}
	if *isClient {
		Client()
	}
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
		NextProtos:   []string{"quic-file"},
	}
}