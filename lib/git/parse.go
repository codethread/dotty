package git

import (
	"github.com/codethread/dotty/lib/fp"
	"io/ioutil"
	"regexp"
	"strings"
)

type Patterns struct {
	Patterns []Pattern
}

func (pats Patterns) Matches(f string) bool {
	_, notFound := fp.Find(func(pat Pattern) bool {
		return pat.Matches(f)

	}, pats.Patterns)

	return notFound == nil
}

func (self Patterns) Concat(pat Patterns) Patterns {
	self.Patterns = append(self.Patterns, pat.Patterns...)
	return self
}

type Pattern struct {
	Regexp  *regexp.Regexp
	Negated bool
}

func (pat Pattern) Matches(f string) bool {
	// TODO: handle negated patterns
	return pat.Regexp.MatchString(f)
}

// ParseGitIgnoreFile uses an ignore file as the input, parses the lines out of
// the file and converts them to a list of Regexp
//
// Function stolen and edited from [ignore.go](https://github.com/sabhiram/go-gitignore/blob/master/ignore.go)
func ParseGitIgnoreFile(fpath string) (Patterns, error) {
	bs, err := ioutil.ReadFile(fpath)

	if err != nil {
		return Patterns{}, err
	}

	s := strings.Split(string(bs), "\n")

	r := fp.FilterMap(func(l string) (Pattern, bool) {
		rg, negated := getPatternFromLine(l)

		if rg == nil {
			return Pattern{}, false
		}

		return Pattern{Regexp: rg, Negated: negated}, true
	}, s)

	return Patterns{Patterns: r}, nil
}

// This function pretty much attempts to mimic the parsing rules
// listed above at the start of this file
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
