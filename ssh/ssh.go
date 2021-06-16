package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

//SSHSession - берём конфиг(собирается из пользователя и ключа, так что эти два параметра под вопросом), берём ключ для входа на сервера. ставим таймаут для команд,
//если такой вообще есть в кофниге, передаём из базы список хостов, и парсим сами команды (пихать судя по всему будем сплошным пайпом в стдин хостов)
//ну и пользователя ставим само собой
func SSHSession(ClientConfig ssh.ClientConfig, Key []byte, Timeout time.Duration, hosts []string, command []string, user string) {
	timeout := time.After(Timeout)
	result := make(chan string)
	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-result:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}
