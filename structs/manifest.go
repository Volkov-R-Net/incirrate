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
