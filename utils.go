package myanimelist

import (
	"crypto/rand"
	"log"
)

func makeList(objects []string) map[string]struct{} {
	var list = make(map[string]struct{}, len(objects))
	if len(objects) < 1 {
		return list
	}
	for _, v := range objects {
		list[v] = struct{}{}
	}
	return list
}

func makeListInt(objects []int) map[int]struct{} {
	var list = make(map[int]struct{}, len(objects))
	if len(objects) < 1 {
		return list
	}
	for _, v := range objects {
		list[v] = struct{}{}
	}
	return list
}

// fixSorting is small internal helper function.
// Main purpose is to avoid creating duplicate variables
// such as SortAnimeListByStartTime && SortMangaListByStartTime
func fixSorting(sort string, kind string) string {
	var prefix string
	if kind == "anime" {
		prefix = "anime_"
	} else if kind == "manga" {
		prefix = "manga_"
	} else {
		// undefined kind
		return sort
	}
	if sort == SortListByTitle || sort == SortListByStartDate || sort == SortListByID {
		return prefix + sort
	}
	return sort
}

func randomString(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	charsAmount := len(chars)
	if charsAmount < 2 || charsAmount > 256 {
		log.Println("randomString() wrong length")
		return ""
	}
	maxNumber := 255 - (256 % charsAmount)
	resultStorage := make([]byte, length)
	bytesStorage := make([]byte, length+(length/4))
	i := 0
	for {
		if _, err := rand.Read(bytesStorage); err != nil {
			log.Println("randomString() error reading random bytes: ", err)
			return ""
		}
		for _, rb := range bytesStorage {
			randomInt := int(rb)
			if randomInt > maxNumber {
				// Skip this number to avoid modulo bias.
				continue
			}
			resultStorage[i] = chars[randomInt%charsAmount]
			i++
			if i == length {
				return string(resultStorage)
			}
		}
	}
}