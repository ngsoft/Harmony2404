package main

import (
	"bufio"
	"flirc/util"
	"os"
	"path"
	"strings"
)

type Keymap string
type Keycode int

var (
	KeyCodes = map[string]Keycode{
		"KEY_A":                30,
		"KEY_B":                48,
		"KEY_C":                46,
		"KEY_D":                32,
		"KEY_E":                18,
		"KEY_F":                33,
		"KEY_G":                34,
		"KEY_H":                35,
		"KEY_I":                23,
		"KEY_J":                36,
		"KEY_K":                37,
		"KEY_L":                38,
		"KEY_M":                50,
		"KEY_N":                49,
		"KEY_O":                24,
		"KEY_P":                25,
		"KEY_Q":                16,
		"KEY_R":                19,
		"KEY_S":                31,
		"KEY_T":                20,
		"KEY_U":                22,
		"KEY_V":                47,
		"KEY_W":                17,
		"KEY_X":                45,
		"KEY_Y":                21,
		"KEY_Z":                44,
		"KEY_1":                2,
		"KEY_2":                3,
		"KEY_3":                4,
		"KEY_4":                5,
		"KEY_5":                6,
		"KEY_6":                7,
		"KEY_7":                8,
		"KEY_8":                9,
		"KEY_9":                10,
		"KEY_0":                11,
		"KEY_ENTER":            28,
		"KEY_ESC":              1,
		"KEY_BACKSPACE":        14,
		"KEY_TAB":              15,
		"KEY_SPACE":            57,
		"KEY_MINUS":            12,
		"KEY_EQUAL":            13,
		"KEY_LEFTBRACE":        26,
		"KEY_RIGHTBRACE":       27,
		"KEY_BACKSLASH":        43,
		"KEY_SEMICOLON":        39,
		"KEY_APOSTROPHE":       40,
		"KEY_GRAVE":            41,
		"KEY_COMMA":            51,
		"KEY_DOT":              52,
		"KEY_SLASH":            53,
		"KEY_CAPSLOCK":         58,
		"KEY_F1":               59,
		"KEY_F2":               60,
		"KEY_F3":               61,
		"KEY_F4":               62,
		"KEY_F5":               63,
		"KEY_F6":               64,
		"KEY_F7":               65,
		"KEY_F8":               66,
		"KEY_F9":               67,
		"KEY_F10":              68,
		"KEY_F11":              87,
		"KEY_F12":              88,
		"KEY_SYSRQ":            99,
		"KEY_SCROLLLOCK":       70,
		"KEY_PAUSE":            119,
		"KEY_INSERT":           110,
		"KEY_HOME":             102,
		"KEY_PAGEUP":           104,
		"KEY_DELETE":           111,
		"KEY_END":              107,
		"KEY_PAGEDOWN":         109,
		"KEY_RIGHT":            106,
		"KEY_LEFT":             105,
		"KEY_DOWN":             108,
		"KEY_UP":               103,
		"KEY_NUMLOCK":          69,
		"KEY_KPSLASH":          98,
		"KEY_KPASTERISK":       55,
		"KEY_KPMINUS":          74,
		"KEY_KPPLUS":           78,
		"KEY_KPENTER":          96,
		"KEY_KP1":              79,
		"KEY_KP2":              80,
		"KEY_KP3":              81,
		"KEY_KP4":              75,
		"KEY_KP5":              76,
		"KEY_KP6":              77,
		"KEY_KP7":              71,
		"KEY_KP8":              72,
		"KEY_KP9":              73,
		"KEY_KP0":              82,
		"KEY_KPDOT":            83,
		"KEY_102ND":            86,
		"KEY_COMPOSE":          127,
		"KEY_POWER":            116,
		"KEY_KPEQUAL":          117,
		"KEY_F13":              183,
		"KEY_F14":              184,
		"KEY_F15":              185,
		"KEY_F16":              186,
		"KEY_F17":              187,
		"KEY_F18":              188,
		"KEY_F19":              189,
		"KEY_F20":              190,
		"KEY_F21":              191,
		"KEY_F22":              192,
		"KEY_F23":              193,
		"KEY_F24":              194,
		"KEY_OPEN":             134,
		"KEY_HELP":             138,
		"KEY_PROPS":            130,
		"KEY_FRONT":            132,
		"KEY_STOP":             128,
		"KEY_AGAIN":            129,
		"KEY_UNDO":             131,
		"KEY_CUT":              137,
		"KEY_COPY":             133,
		"KEY_PASTE":            135,
		"KEY_FIND":             136,
		"KEY_MUTE":             113,
		"KEY_VOLUMEUP":         115,
		"KEY_VOLUMEDOWN":       114,
		"KEY_KPCOMMA":          121,
		"KEY_RO":               89,
		"KEY_KATAKANAHIRAGANA": 93,
		"KEY_YEN":              124,
		"KEY_HENKAN":           92,
		"KEY_MUHENKAN":         94,
		"KEY_KPJPCOMMA":        95,
		"KEY_HANGEUL":          122,
		"KEY_HANJA":            123,
		"KEY_KATAKANA":         90,
		"KEY_HIRAGANA":         91,
		"KEY_ZENKAKUHANKAKU":   85,
		"KEY_KPLEFTPAREN":      179,
		"KEY_KPRIGHTPAREN":     180,
		"KEY_LEFTCTRL":         29,
		"KEY_LEFTSHIFT":        42,
		"KEY_LEFTALT":          56,
		"KEY_LEFTMETA":         125,
		"KEY_RIGHTCTRL":        97,
		"KEY_RIGHTSHIFT":       54,
		"KEY_RIGHTALT":         100,
		"KEY_RIGHTMETA":        126,
		"KEY_PLAYPAUSE":        164,
		"KEY_STOPCD":           166,
		"KEY_PREVIOUSSONG":     165,
		"KEY_NEXTSONG":         163,
		"KEY_EJECTCD":          161,
		"KEY_WWW":              150,
		"KEY_BACK":             158,
		"KEY_FORWARD":          159,
		"KEY_SCROLLUP":         177,
		"KEY_SCROLLDOWN":       178,
		"KEY_EDIT":             176,
		"KEY_SLEEP":            142,
		"KEY_COFFEE":           152,
		"KEY_REFRESH":          173,
		"KEY_CALC":             140,
	}
)

func LoadKeymaps() []KeyPair {

	var (
		ok       bool
		dirName  string
		fileName string
		result   = make([]KeyPair, 0)
	)

	if dirName, ok = util.FindPath(cfgDir); ok {
		fileName = path.Join(dirName, "flirc.keymaps")
		if f, err := os.Open(fileName); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {

				list := strings.Split(scanner.Text(), "=")
				if len(list) == 2 {
					key, value := strings.TrimSpace(list[0]), Keymap(strings.TrimSpace(list[1]))
					if code, ok := KeyCodes[key]; ok {
						result = append(result, KeyPair{code, value})
					}

				}
			}
		}
	}
	return result
}
