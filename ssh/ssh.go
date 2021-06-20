// Copyright (C) 2019-2021 Volkov Roman Olegovich.
//
// This file is part of Incirrate.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software Foundation,
// Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

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
