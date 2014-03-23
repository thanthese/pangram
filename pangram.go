package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"
var threshold = 24

type set map[rune]bool

func main() {
	if len(os.Args) < 1 {
		log.Fatal("pangram requires path argument")
	}
	path := os.Args[1]
	fmt.Println("Loading word list...")
	wordlist := loadWordList(path)

	fmt.Println("Building anagrams list...")
	singlesOnly := removeDoubles(wordlist)
	anagrams := buildAnagrams(singlesOnly)
	words := mapKeys(anagrams)
	fmt.Println("- ", len(wordlist), "list length")
	fmt.Println("- ", len(singlesOnly), "with unique letters")
	fmt.Println("- ", len(anagrams), "anagrams")

	fmt.Println("Finding pangrams...")
	used := set{}
	found := []string{}
	recur(used, found, words, anagrams)
}

func mapKeys(m map[string][]string) []string {
	ret := make([]string, 0, len(m))
	for k, _ := range m {
		ret = append(ret, k)
	}
	return ret
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

func prettyFinding(words []string, anagrams map[string][]string) {
	fmt.Print(runesCount(words))
	for _, w := range words {
		if len(w) > 0 {
			fmt.Print(" ", anagrams[w])
		}
	}
	fmt.Println()
}

func recur(used set, foundwords []string, potentials []string, anagrams map[string][]string) {
	if len(used) == 26 || len(potentials) == 0 {
		if len(used) >= threshold {
			prettyFinding(foundwords, anagrams)
		}
		return
	}

	for i, word := range potentials {

		// prepare new set
		u := copymap(used)
		for _, r := range word {
			u[r] = true
		}

		// prepare new foundwords
		fw := append(foundwords, word)

		// prepare (filter) new potentials
		rm := make([]string, 0, len(potentials))
		for j := i + 1; j < len(potentials); j++ {
			w := potentials[j]
			if len(w)+len(u) <= 26 && wordFits(u, w) {
				rm = append(rm, w)
			}
		}

		recur(u, fw, rm, anagrams)
	}
}

func runesCount(words []string) int {
	sum := 0
	for _, word := range words {
		sum += len(word)
	}
	return sum
}

func wordFits(used set, word string) bool {
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
