package main

import (
	"strconv"
	"strings"
	"unicode"
)

func ProcessText(text string) string {
	words := strings.Fields(text)
	var out []string

	for _, w := range words {
		switch w {
		case "(hex)":
			convertBase(&out, 16)
		case "(bin)":
			convertBase(&out, 2)
		case "(up)":
			applyTransform(&out, strings.ToUpper)
		case "(low)":
			applyTransform(&out, strings.ToLower)
		case "(cap)":
			applyTransform(&out, capOne)
		default:
			out = append(out, w)
		}
	}

	result := strings.Join(out, " ")
	result = fixPunctuation(result)
	result = fixArticles(result)
	result = fixQuotes(result)

	return result
}

func convertBase(out *[]string, base int) {
	if len(*out) == 0 {
		return
	}
	last := (*out)[len(*out)-1]
	if n, err := strconv.ParseInt(last, base, 64); err == nil {
		(*out)[len(*out)-1] = strconv.FormatInt(n, 10)
	}
}

func applyTransform(out *[]string, fn func(string) string) {
	if len(*out) == 0 {
		return
	}
	(*out)[len(*out)-1] = fn((*out)[len(*out)-1])
}

func capOne(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

func fixPunctuation(s string) string {
	puncs := []string{".", ",", "!", "?", ":", ";"}

	for _, p := range puncs {
		s = strings.ReplaceAll(s, " "+p, p)
	}
	
	for _, p := range puncs {
		s = strings.ReplaceAll(s, p, p+" ")
	}

	return strings.Join(strings.Fields(s), " ")
}

func fixArticles(s string) string {
	words := strings.Fields(s)
	for i := 0; i < len(words)-1; i++ {
		if strings.EqualFold(words[i], "a") {
			first := rune(words[i+1][0])
			if strings.ContainsRune("aeiouhAEIOUH", first) {
				if words[i] == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}
	return strings.Join(words, " ")
}

func fixQuotes(s string) string {
	s = strings.ReplaceAll(s, "' ", "'")
	s = strings.ReplaceAll(s, " '", "'")
	return s
}
