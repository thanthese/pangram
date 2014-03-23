package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// todo
//
// - There seems to be a bug, or at least something I don't understand, where
// sets of words are printed twice.
//
// - I don't like how recur() prints. It would be better if it returned its
// results (preferrably via channel, so the user can still watch results come
// in in real time) and the caller dealt with printing.
//
// - It would be nice if recur was executing in multiple goroutines in
// parallel, for speed and whatnot.
//
// - None of my results contain a 'q' or 'z'. Maybe I can *force* the first two
// words to have those letters. Then I could use a larger dictionary, but skip
// the cost of a ton of computation.

var alphabet = "abcdefghijklmnopqrstuvwxyz"

var threshold int

// we'll fake a set with a map
type set map[rune]bool

func usage() {
	fmt.Println("pangram <threshold> <dict path>")
	fmt.Println()
	fmt.Println("- threshold: display words if they use >= threshold unique letters")
	fmt.Println("- dict path: path to list of words. One word per line, lower case, a-z only.")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	th, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		usage()
	}
	threshold = th

	fmt.Println("Loading word list...")
	path := os.Args[2]
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
		fmt.Println(err.Error())
		usage()
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

// Filter wordlist so that no words contain a letter more than once. For
// example, "house" would be allowed, but "example" would be filtered out.
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
		ps := make([]string, 0, len(potentials))
		for j := i + 1; j < len(potentials); j++ {
			w := potentials[j]
			if len(w)+len(u) <= 26 && wordFits(u, w) {
				ps = append(ps, w)
			}
		}

		recur(u, fw, ps, anagrams)
	}
}

// Return the sum of lengths of all words in the list.
func runesCount(words []string) int {
	sum := 0
	for _, word := range words {
		sum += len(word)
	}
	return sum
}

// Return true if the word does not contain any letters that have already been
// checked off by the used set. That is, the word only contains still allowable
// letters.
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
