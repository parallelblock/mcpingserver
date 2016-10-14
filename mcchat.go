package main

import(
    "regexp"
)

type TextSelection struct {
    Text string `json:"text"`
    Bold bool `json:"bold,omitempty"`
    Italic bool `json:"italic,omitempty"`
    Underlined bool `json:"underlined,omitempty"`
    Strikethrough bool `json:"strikethrough,omitempty"`
    Obfuscated bool `json:"obfuscated,omitempty"`
    Color string `json:"color,omitempty"`
    Extra []interface{} `json:"extra,omitempty"`
}

var COLOR_MAP = map[rune]string{
    '0': "black",
    '1': "dark_blue",
    '2': "dark_green",
    '3': "dark_aqua",
    '4': "dark_red",
    '5': "dark_purple",
    '6': "gold",
    '7': "gray",
    '8': "dark_gray",
    '9': "blue",
    'a': "green",
    'b': "aqua",
    'c': "red",
    'd': "light_purple",
    'e': "yellow",
    'f': "white",
}

func resetTextModifiers(a *TextSelection) {
    a.Bold = false
    a.Italic = false
    a.Underlined = false
    a.Strikethrough = false
    a.Obfuscated = false
}

func createNewChild(parent *TextSelection) (*TextSelection) {
    child := &TextSelection{"", false, false, false, false, false, "", []interface{}{}}
    parent.Extra = append(parent.Extra, child)
    return child
}


func isColorChar(r rune) bool {
    return (r >= '0' && r <= '9') ||
        (r >= 'a' && r <= 'f') || r == 'r'
}

var translateColorCodeRegex = regexp.MustCompile("/&([0-9a-fA-Fk-oK-OrR])/")

func TranslateColorCodes(input string) string {
    return string(translateColorCodeRegex.ReplaceAll([]byte(input), []byte("§${1}")))
}

func ConvertMCChat(input string) *TextSelection {
    parent := &TextSelection{"", false, false, false, false, false, "", []interface{}{}}
    working := createNewChild(parent)

    escape := false
    colorcode := false
    startedtext := false
    startingptr := 0

    for idx, r := range input {
        if escape {
            escape = false
            // if its n, newline it
            // otherwise, take out /
            working.Text = input[startingptr:idx - 1]
            if r == 'n' {
                working.Text += string('\u000a')
            } else {
                working.Text += string(r)
            }
            startedtext = true
            startingptr = idx + 1
            colorcode = false
            continue
        }

        if r == '\\' {
            escape = true
            continue
        }

        if colorcode {
            colorcode = false
            if r >= 'k' && r <= 'o' {
                if startedtext {
                    working.Text += input[startingptr:idx-1]
                    working = createNewChild(parent)
                    startedtext = false
                }
                if r == 'k' {
                    working.Obfuscated = true
                } else if r == 'l' {
                    working.Bold = true
                } else if r == 'm' {
                    working.Strikethrough = true
                } else if r == 'n' {
                    working.Underlined = true
                } else {
                    working.Italic = true
                }
                continue
            } else if isColorChar(r) || r == 'r' {
                if startedtext {
                    working.Text += input[startingptr:idx-1]
                    working = createNewChild(parent)
                    startedtext = false
                    if r != 'r' {
                        working.Color = COLOR_MAP[r]
                    }
                } else {
                    resetTextModifiers(working)
                    if r != 'r' {
                        working.Color = COLOR_MAP[r]
                    }
                }
                continue
            }
        }

        if r == '§' {
            colorcode = true
            continue
        }

        if !startedtext {
            startedtext = true
            startingptr = idx
            colorcode = false
        }
    }
    if startedtext {
        working.Text += input[startingptr:]
    } else {
        parent.Extra = parent.Extra[:len(parent.Extra)-1]
    }
    return parent
}
