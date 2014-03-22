package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"
var threshold = 21

type set map[rune]bool

func main() {
	wordlist := loadWordList("/Users/thanthese/pangram/simplewords.txt")
	singlesOnly := removeDoubles(wordlist)
	anagrams := buildAnagrams(singlesOnly)
	_ = anagrams
	used := set{}
	found := []string{}
	recur(used, found, singlesOnly)
}

func loadWordList(path string) []string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}

func buildAnagrams(words []string) map[string][]string {
	anagrams := map[string][]string{}
	for _, word := range words {
		w := sortWord(word)
		if anagrams[w] == nil {
			anagrams[w] = []string{word}
		} else {
			anagrams[w] = append(anagrams[w], word)
		}
	}
	return anagrams
}

func sortWord(word string) string {
	runes := []string{}
	for _, r := range word {
		runes = append(runes, string(r))
	}
	sort.Strings(runes)
	return strings.Join(runes, "")
}

func removeDoubles(wordList []string) []string {
	var ret []string
	for _, word := range wordList {
		if !containsDoubles(word) {
			ret = append(ret, word)
		}
	}
	return ret
}

func containsDoubles(word string) bool {
	used := make(set)
	for _, r := range word {
		if used[r] {
			return true
		} else {
			used[r] = true
		}
	}
	return false
}

func recur(used set, foundwords []string, remainingwords []string) {
	if len(used) >= threshold {
		fmt.Println(runesCount(foundwords), foundwords)
		return
	}

	if len(remainingwords) == 0 {
		return
	}

	for i, word := range remainingwords {
		if isValidAddition(used, word) {
			u := copymap(used)
			for _, c := range word {
				u[c] = true
			}
			fw := append(foundwords, word)
			recur(u, fw, remainingwords[i+1:])
		}
	}
}

func runesCount(words []string) int {
	sum := 0
	for _, word := range words {
		sum += len(word)
	}
	return sum
}

func isValidAddition(used set, word string) bool {
	for _, c := range word {
		if used[c] {
			return false
		}
	}
	return true
}

func copymap(m set) set {
	n := make(set)
	for k, v := range m {
		n[k] = v
	}
	return n
}
