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

package manifest

import "strings"

type Manifest struct {
	Name      string `yaml:"name"`                // Print name
	Sudo      bool   `yaml:"sudo,omitempty"`      // Sudo -i or sudo su
	Sudo_user string `yaml:"sudo_user,omitempty"` // su username
	Copy      []struct {
		Values  []string `yaml:"values,flow,omitempty"` // через этот пункт надо прогнать указанные файлы
		Project string   `yaml:"project,omitempty"`     // Из проекта и файла формируем путь
		File    []string `yaml:"file,flow"`             // относительно домашней директории
		Target  string   `yaml:"target"`                // Куда копируем
		User    string   `yaml:"user,omitempty"`        // Права - тут пользователь
		Group   string   `yaml:"group,omitempty"`       // Тут группа
		Chmod   int      `yaml:"chmod, omitempty"`      // Тут выставляются права в виде цифр 644,777 и тд
	} `yaml:"copy,flow,omitempty"` //добавить варнинги, данный блок объявлен, но его длинна равна нулю
	Command []struct {
		Values  []string `yaml:"values,flow,omitempty"`
		Cmd     string   `yaml:"cmd,omitempty"`
		Project string   `yaml:"project,omitempty"`
		Script  string   `yaml:"script,omitempty"`
	} `yaml:"command,flow,omitempty"` //добавить варнинги, данный блок объявлен, но его длинна равна нулю
}

func (a *Manifest) GetFiles(manifest Manifest) (result []string) {
	for _, b := range manifest.Copy {
		for _, c := range b.File {
			result = append(result, b.Project+"/"+c)
		}
	}
	return result
}

func (a *Manifest) FormatString(specification []string, template string) (result string) {
	for _, b := range specification {
		spec := strings.SplitN(b, "=", 2)
		key := spec[0]
		value := spec[1]
		replacer := strings.NewReplacer("{{ "+key+" }}", value)
		result = replacer.Replace(template)
		template = result
	}
	return result
}
