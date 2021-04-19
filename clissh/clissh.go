package clissh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

//см. https://habr.com/ru/post/215111/ и https://yanzay.com/post/go_ssh_client_example/
type SignerContainer struct {
	signers []ssh.Signer
}

func (t *SignerContainer) Key(i int) (key ssh.PublicKey, err error) {
	if i >= len(t.signers) {
		return
	}
	key = t.signers[i].PublicKey()
	return
}

func (t *SignerContainer) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	if i >= len(t.signers) {
		return
	}
	sig, err = t.signers[i].Sign(rand, data)
	return
}

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, _ := ioutil.ReadAll(fp)
	signer, _ = ssh.ParsePrivateKey(buf)
	return
}

func makeKeyring() ssh.ClientAuth {
	signers := []ssh.Signer{}
	keys := []string{os.Getenv("HOME") + "/.ssh/id_rsa", os.Getenv("HOME") + "/.ssh/id_dsa"}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return ssh.ClientAuthKeyring(&SignerContainer{signers})
}

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
	conn, _ := ssh.Dial("tcp", hostname+":22", config)
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	return hostname + ": " + stdoutBuf.String()
}

//всё надо собрать в единую функцию, которую потом будем брать из кейсов в main
func SSHconnect() { //(user string, config ssh.ClientConfig)
	comand := os.Args[1] //команда, которую исполняем на хостах
	hosts := os.Args[2:] //список хостов

	results := make(chan string, 10)        //результаты в буфер строк
	timeout := time.After(10 * time.Second) //в канал timeout придет сообщение(например 10)

	//инициализируем структуру с конфигурацией для пакета ssh. Функцию makeKeyring() напишем позднее
	config := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.ClientAuth{makeKeyring()},
	}

	//запустим по одной goroutine на хост, функцию executeCmd() напишем позднее
	for _, hostname := range hosts {
		go func(hostname string) {
			results <- executeCmd(cmd, hostname, config)
		}(hostname)
	}

	//соберем результаты со всех серверов, или напишем "Timed out", если общее время исполнения истекло
	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}
