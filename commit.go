package convcommit

import (
	"errors"
	"regexp"
	"strings"
)

type Commit struct {
	Type             string
	Scope            string
	Title            string
	Body             string
	Footer           map[string][]string
	IsBreakingChange bool
}

// Regexes are compiled once at package level, which is already good for performance.
var headerRegex = regexp.MustCompile(`^(?P<type>\w+)(?:\((?P<scope>[\w-]+)\))?(?P<breaking>!)?: (?P<subject>.+)$`)
var footerBlockIdentifierRegex = regexp.MustCompile(`^((?:BREAKING-CHANGE|BREAKING CHANGE)|[\w-]+)(: | #)`)
var footerLineRegex = regexp.MustCompile(`^(?P<token>[\w-]+|BREAKING CHANGE|BREAKING-CHANGE)(?P<separator>: | #)(?P<value>.*)$`)

// --- PERFORMANCE OPTIMIZATION ---
// Pre-calculate the indices for the named capture groups to avoid map lookups in the loop.
var (
	footerTokenIndex = -1
	footerValueIndex = -1
)

func init() {
	for i, name := range footerLineRegex.SubexpNames() {
		switch name {
		case "token":
			footerTokenIndex = i
		case "value":
			footerValueIndex = i
		}
	}
}

func Parse(commitMessage string) (*Commit, error) {
	parts := strings.SplitN(strings.TrimSpace(commitMessage), "\n", 2)
	if len(parts) == 0 || parts[0] == "" {
		return nil, errors.New("convcommit: empty or invalid commit message")
	}

	header := parts[0]
	rest := ""
	if len(parts) > 1 {
		rest = parts[1]
	}

	matches := headerRegex.FindStringSubmatch(header)
	if len(matches) == 0 {
		return nil, errors.New("convcommit: could not parse commit header")
	}
	headerParts := make(map[string]string)
	for i, name := range headerRegex.SubexpNames() {
		if i > 0 && name != "" {
			headerParts[name] = matches[i]
		}
	}

	bodyAndFooterLines := strings.Split(rest, "\n")
	var bodyLines []string
	var footerLines []string
	footerStartIndex := -1

	for i, line := range bodyAndFooterLines {
		if footerBlockIdentifierRegex.MatchString(line) {
			footerStartIndex = i
			break
		}
	}

	if footerStartIndex != -1 {
		bodyLines = bodyAndFooterLines[0:footerStartIndex]
		footerLines = bodyAndFooterLines[footerStartIndex:]
	} else {
		bodyLines = bodyAndFooterLines
	}
	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))

	footerMap := make(map[string][]string)
	var lastToken string

	for _, line := range footerLines {
		lineMatches := footerLineRegex.FindStringSubmatch(line)

		// --- PERFORMANCE OPTIMIZATION ---
		// The previous code created a map for every line. This new version uses
		// pre-calculated indices, avoiding map allocations and lookups inside this hot loop.
		if len(lineMatches) > 0 {
			// Directly access matches using pre-calculated indices.
			token := lineMatches[footerTokenIndex]
			value := lineMatches[footerValueIndex]

			switch strings.ToLower(token) {
			case "breaking-change":
				token = "BREAKING CHANGE"
			case "fix", "fixes", "close", "closes":
				token = "Closes"
			}

			footerMap[token] = append(footerMap[token], value)
			lastToken = token
		} else {
			// This part for multi-line footers still uses `+=`. While not perfectly optimal,
			// it prioritizes readability, as true multi-line footers are rare.
			// Optimizing this further would require a more complex implementation (e.g., using strings.Builder).
			if lastToken != "" && len(footerMap[lastToken]) > 0 {
				lastElementIndex := len(footerMap[lastToken]) - 1
				footerMap[lastToken][lastElementIndex] += "\n" + line
			}
		}
	}

	isBreakingChange := headerParts["breaking"] == "!" || len(footerMap["BREAKING CHANGE"]) > 0

	newCommit := &Commit{
		Type:             headerParts["type"],
		Title:            headerParts["subject"],
		IsBreakingChange: isBreakingChange,
		Body:             body,
		Footer:           footerMap,
		Scope:            headerParts["scope"],
	}

	return newCommit, nil
}
