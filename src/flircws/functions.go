package main

import (
	"bufio"
	"os"
	"path"
	"strings"
)

func LoadKeymaps() []Keymap {

	result := make([]Keymap, 0)

	for _, pth := range etc {
		pwd, _ := os.Getwd()
		pth = path.Join(pwd, pth, cfgDir, "flirc.keymaps")

		f, err := os.Open(pth)
		if err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {

				list := strings.Split(scanner.Text(), "=")
				if len(list) == 2 {

					result = append(result, Keymap(strings.TrimSpace(list[1])))
				}

			}

		}

		// fmt.Println(err.Error())

	}
	return result
}
