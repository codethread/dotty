package git

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// Patterns is a collection of Regexp `Pattern` which supports negation
type Patterns struct {
	Patterns []pattern
}

func New(rs []*regexp.Regexp) Patterns {
	var patterns []pattern
	for i, rs := range rs {
		patterns[i] = pattern{regexp: rs, negated: false}
	}
	return Patterns{Patterns: patterns}
}

func Empty() Patterns {
	return New([]*regexp.Regexp{})
}

func (pats Patterns) Matches(f string) bool {
	for _, p := range pats.Patterns {
		if p.Matches(f) {
			return true
		}
	}
	return false
}

func (self Patterns) Concat(pat Patterns) Patterns {
	self.Patterns = append(self.Patterns, pat.Patterns...)
	return self
}

type pattern struct {
	regexp  *regexp.Regexp
	negated bool
}

func (pat pattern) Matches(f string) bool {
	// TODO: handle negated patterns
	return pat.regexp.MatchString(f)
}

// ParseFile uses an ignore file as the input, parses the lines out of
// the file and converts them to a list of Regexp
func ParseFile(fpath string) (Patterns, error) {
	bs, err := ioutil.ReadFile(fpath)

	if err != nil {
		return Patterns{}, err
	}

	var pats []pattern

	for _, l := range strings.Split(string(bs), "\n") {
		rg, negated := getPatternFromLine(l)

		if rg != nil {
			pats = append(pats, pattern{regexp: rg, negated: negated})
		}
	}

	return Patterns{Patterns: pats}, nil
}

// This function pretty much attempts to mimic the parsing rules
// listed above at the start of this file
//
// Function stolen and edited from [ignore.go](https://github.com/sabhiram/go-gitignore/blob/master/ignore.go)
func getPatternFromLine(line string) (*regexp.Regexp, bool) {
	// Trim OS-specific carriage returns.
	line = strings.TrimRight(line, "\r")

	// Strip comments [Rule 2]
	if strings.HasPrefix(line, `#`) {
		return nil, false
	}

	// Trim string [Rule 3]
	// TODO: Handle [Rule 3], when the " " is escaped with a \
	line = strings.Trim(line, " ")

	// Exit for no-ops and return nil which will prevent us from
	// appending a pattern against this line
	if line == "" {
		return nil, false
	}

	// TODO: Handle [Rule 4] which negates the match for patterns leading with "!"
	negatePattern := false
	if line[0] == '!' {
		negatePattern = true
		line = line[1:]
	}

	// Handle [Rule 2, 4], when # or ! is escaped with a \
	// Handle [Rule 4] once we tag negatePattern, strip the leading ! char
	if regexp.MustCompile(`^(\#|\!)`).MatchString(line) {
		line = line[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// Handle escaping the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.Replace(line, "?", `\?`, -1)

	line = strings.Replace(line, magicStar, "*", -1)

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "(|/.*)$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern, negatePattern
}
