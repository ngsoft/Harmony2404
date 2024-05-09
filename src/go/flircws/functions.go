package main

import (
	"bufio"
	"flirc/util"
	"os"
	"path"
	"strings"
)

func LoadKeymaps() []Keymap {

	var (
		ok       bool
		dirName  string
		fileName string
		result   = make([]Keymap, 0)
	)

	if dirName, ok = util.FindPath(cfgDir); ok {
		fileName = path.Join(dirName, "flirc.keymaps")
		if f, err := os.Open(fileName); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {

				list := strings.Split(scanner.Text(), "=")
				if len(list) == 2 {
					result = append(result, Keymap(strings.TrimSpace(list[1])))
				}
			}
		}
	}
	return result
}
