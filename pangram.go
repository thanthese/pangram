package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// todo
//
// - bug: Total number of leaves seen changes from run to run. Why?! There
// shouldn't be anything random expect thread execution order, and those should
// be entirely independent.

const alphabet = "abcdefghijklmnopqrstuvwxyz"

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
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) < 2 {
		usage()
	}

	threshold, err := strconv.Atoi(os.Args[1])
	if err != nil || threshold < 0 {
		fmt.Println(err.Error())
		usage()
	}

	fmt.Println("Loading word list...")
	path := os.Args[2]
	wordlist := loadWordList(path)

	PrintPangrams(threshold, wordlist)
}

func PrintPangrams(threshold int, wordlist []string) {
	fmt.Println("Building anagrams list...")
	singlesOnly := removeDoubles(wordlist)
	anagrams := buildAnagrams(singlesOnly)
	words := mapKeys(anagrams)
	sorted := specialSort(words)
	fmt.Println("- ", len(wordlist), "words")
	fmt.Println("- ", len(singlesOnly), "with unique letters")
	fmt.Println("- ", len(anagrams), "anagrams")

	fmt.Println("Finding pangrams...")
	used := set{}
	found := []string{}
	out := make(chan []string)
	recur(used, found, sorted, out, nil)

	n := 0
	for words := range out {
		n++
		if runesCount(words) >= threshold {
			prettyFinding(words, anagrams)
		}
	}
	fmt.Printf("Total leaves seen: %d\n", n)
}

func recur(used set, foundwords []string, potentials []string,
	out chan []string, done chan int) {

	if len(used) == 26 || len(potentials) == 0 {
		ret := listcopy(foundwords)
		out <- ret
	}

	threads := 0
	d := make(chan int)
	for i, word := range potentials {

		if len(foundwords) == 0 {
			if !strings.ContainsAny(word, "q") {
				continue
			}
		}

		if len(foundwords) == 1 {
			if !strings.ContainsAny(word, "z") &&
				!strings.ContainsAny(foundwords[0], "z") {

				continue
			}
		}

		// prepare new set
		u := copymap(used)
		for _, r := range word {
			u[r] = true
		}

		// prepare new foundwords
		fw := listcopy(foundwords)
		fw = append(fw, word)

		// prepare (filter) new potentials
		ps := make([]string, 0, len(potentials))
		for j := i + 1; j < len(potentials); j++ {
			w := potentials[j]
			if len(w)+len(u) <= 26 && wordFits(u, w) {
				ps = append(ps, w)
			}
		}

		if len(foundwords) == 0 {
			threads++
			go recur(u, fw, ps, out, d)
		} else {
			recur(u, fw, ps, out, nil)
		}
	}

	if len(foundwords) == 0 {
		go func() {
			for i := 0; i < threads; i++ {
				<-d
				// fmt.Printf("%d threads remaining\n", threads-i)
			}
			close(d)
			close(out)
		}()
	} else if len(foundwords) == 1 {
		done <- 1
	}
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
	words := strings.Split(string(content), "\n")

	for _, word := range words {
		for _, r := range word {
			if !strings.ContainsRune(alphabet, r) {
				fmt.Printf("Error: Word \"%s\" contains invalid character '%s'\n", word, string(r))
				usage()
			}
		}
	}

	nonzeroWords := make([]string, 0, len(words))
	for _, w := range words {
		if len(w) > 0 {
			nonzeroWords = append(nonzeroWords, w)
		}
	}

	if len(nonzeroWords) < 10 {
		fmt.Printf("Error: file \"%s\" doesn't contain enough words (only %d).\n",
			path, len(nonzeroWords))
		usage()
	}

	return nonzeroWords
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
	sort.Strings(words)
	for _, w := range words {
		fmt.Print(" ", anagrams[w])
	}
	fmt.Println()
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
	for _, r := range word {
		if used[r] {
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

func listcopy(list []string) []string {
	c := make([]string, 0, len(list))
	for _, n := range list {
		c = append(c, n)
	}
	return c
}

func specialSort(words []string) []string {
	q, rest := separateBy(words, "q")
	z, rest2 := separateBy(rest, "z")

	t := append(q, z...)
	return append(t, rest2...)
}

func separateBy(words []string, r string) (has []string, hasnot []string) {
	for _, w := range words {
		if strings.Contains(w, r) {
			has = append(has, w)
		} else {
			hasnot = append(hasnot, w)
		}
	}
	return
}
