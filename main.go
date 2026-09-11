package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
)

const version = "1.0.0"

type analyzer struct {
	total       int
	passwords   map[string]int
	baseWords   map[string]int
	lengths     map[int]int
	charSets    map[string]int
	patterns    map[string]int
	suffixes    map[int]map[string]int
	firstCapNum int
	firstCapSym int
	oneToSix    int
	oneToEight  int
	overEight   int
}

func newAnalyzer() *analyzer {
	return &analyzer{
		passwords: make(map[string]int), baseWords: make(map[string]int), lengths: make(map[int]int),
		charSets: make(map[string]int), patterns: make(map[string]int), suffixes: make(map[int]map[string]int),
	}
}

func (a *analyzer) add(password string) {
	if password == "" {
		return
	}
	a.total++
	a.passwords[password]++
	a.lengths[len([]rune(password))]++
	length := len([]rune(password))
	if length <= 6 {
		a.oneToSix++
	}
	if length <= 8 {
		a.oneToEight++
	}
	if length > 8 {
		a.overEight++
	}
	if len(password) > 1 && unicode.IsUpper([]rune(password)[0]) {
		last := []rune(password)[len([]rune(password))-1]
		if unicode.IsDigit(last) {
			a.firstCapNum++
		}
		if unicode.IsPunct(last) {
			a.firstCapSym++
		}
	}
	base := strings.ToLower(password)
	base = strings.Trim(base, "0123456789!@#$%^&*()_+-=[]{}|;:'\",.<>/?\\`~")
	if len([]rune(base)) > 3 {
		a.baseWords[base]++
	}
	a.charSets[characterSet(password)]++
	a.patterns[maskPattern(password)]++
	for n := 1; n <= 5; n++ {
		runes := []rune(password)
		if len(runes) >= n && allDigits(runes[len(runes)-n:]) {
			if a.suffixes[n] == nil {
				a.suffixes[n] = make(map[string]int)
			}
			a.suffixes[n][string(runes[len(runes)-n:])]++
		}
	}
}

func allDigits(runes []rune) bool {
	for _, r := range runes {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func characterSet(s string) string {
	hasLower, hasUpper, hasDigit, hasSpecial := false, false, false, false
	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	parts := []string{}
	if hasLower {
		parts = append(parts, "lower")
	}
	if hasUpper {
		parts = append(parts, "upper")
	}
	if hasDigit {
		parts = append(parts, "number")
	}
	if hasSpecial {
		parts = append(parts, "special")
	}
	return strings.Join(parts, "+")
}

func maskPattern(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			b.WriteByte('l')
		case unicode.IsUpper(r):
			b.WriteByte('u')
		case unicode.IsDigit(r):
			b.WriteByte('d')
		default:
			b.WriteByte('s')
		}
	}
	return b.String()
}

type counted struct {
	name  string
	count int
}

func top(values map[string]int, limit int) []counted {
	items := make([]counted, 0, len(values))
	for name, count := range values {
		items = append(items, counted{name, count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].count == items[j].count {
			return items[i].name < items[j].name
		}
		return items[i].count > items[j].count
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items
}
func percent(count, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(count) * 100 / float64(total)
}

func eachVariant(word string, emit func(string)) {
	runes := []rune(word)
	var walk func(int)
	walk = func(pos int) {
		if pos == len(runes) {
			emit(string(runes))
			return
		}
		original := runes[pos]
		choices := []rune{original, unicode.ToLower(original), unicode.ToUpper(original)}
		if i := strings.IndexRune("aesioltbg", unicode.ToLower(original)); i >= 0 {
			choices = append(choices, rune("@35101789"[i]))
		}
		for i, r := range choices {
			if !strings.ContainsRune(string(choices[:i]), r) {
				runes[pos] = r
				walk(pos + 1)
			}
		}
		runes[pos] = original
	}
	walk(0)
}

func (a *analyzer) report(w io.Writer, topN int, markdown, numbers bool) {
	for _, item := range top(a.passwords, topN) {
		if !numbers {
			fmt.Fprintln(w, item.name)
			continue
		}
		eachVariant(item.name, func(word string) {
			fmt.Fprintln(w, word)
			for _, separator := range []string{"", ".", "@", "#", "$", "!", "%", "&", "*", "_", "-", "+", "="} {
				for _, suffix := range []struct{ width, limit int }{{3, 1000}, {4, 10000}} {
					for number := 0; number < suffix.limit; number++ {
						fmt.Fprintf(w, "%s%s%0*d\n", word, separator, suffix.width, number)
					}
				}
			}
		})
	}
}

func writeTop(w io.Writer, title string, values []counted, total int, markdown bool) {
	if title != "" {
		fmt.Fprintf(w, "%s\n", title)
	}
	for _, item := range values {
		if markdown {
			fmt.Fprintf(w, "- `%s`: %d (%.2f%%)\n", item.name, item.count, percent(item.count, total))
		} else {
			fmt.Fprintf(w, "%s = %d (%.2f%%)\n", item.name, item.count, percent(item.count, total))
		}
	}
	fmt.Fprintln(w)
}

func main() {
	topN := flag.Int("t", 10, "number of top results")
	output := flag.String("o", "", "write report to file")
	markdown := flag.Bool("m", false, "write Markdown output")
	numbers := flag.Bool("numbers", false, "add leetspeak variants and every 3- and 4-digit suffix, with no separator or one of .@#$!%&*_-+=")
	flag.Usage = func() { fmt.Fprintln(os.Stderr, "Usage: lapip [-t N] [-o FILE] [-m] [-numbers] FILE") }
	flag.Parse()
	if *topN <= 0 || flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	input, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer input.Close()
	a := newAnalyzer()
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		a.add(strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
	var out io.Writer = os.Stdout
	var file *os.File
	if *output != "" {
		file, err = os.Create(*output)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		defer file.Close()
		out = file
	}
	a.report(out, *topN, *markdown, *numbers)
}
