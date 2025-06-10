package convcommit_test

import (
	"testing"

	convcommit "github.com/chonla/conv-commit"
	"github.com/stretchr/testify/assert"
)

func TestParseCommitMessages(t *testing.T) {
	// Added new test cases to the end of the existing slices
	commitMessages := []string{
		"feat: allow provided config object to extend other configs",
		"feat(api): allow provided config object to extend other configs",
		"feat!: allow provided config object to extend other configs",
		"feat(api)!: allow provided config object to extend other configs",
		`feat(api): allow provided config object to extend other configs
    
Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.`,
		`feat(api): allow provided config object to extend other configs
    
Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.`,
		`feat(api): allow provided config object to extend other configs
    
Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Reviewed-by: Z
Refs: #123`,
		`feat(api): allow provided config object to extend other configs
    
Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

BREAKING CHANGE: ` + "`extends`" + ` key in config file is now used for extending other config files`,
		`feat(api): allow provided config object to extend other configs
    
Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

BREAKING CHANGE: ` + "`extends`" + ` key in config file is now used for extending other config files
Reviewed-by: Z
Refs: #123`,
		// --- NEW TEST CASES START HERE ---
		`fix(parser): handle new separator and repeated keys

This commit adds support for the ' #' separator and ensures
that repeated footer keys like 'Closes' are aggregated.

Closes: #101
Fixes #102
close #103`,
		`refactor(api-v2): simplify endpoint logic

Reviewed-by: Z`,
		`chore!: drop support for Node 12

BREAKING-CHANGE: Dropped support for Node 12.

This is the second line of the breaking change description.
All consumers must upgrade to Node 14 or higher.`,
	}
	expected := []*convcommit.Commit{
		{
			Type:   "feat",
			Title:  "allow provided config object to extend other configs",
			Footer: map[string][]string{},
		},
		{
			Type:   "feat",
			Title:  "allow provided config object to extend other configs",
			Scope:  "api",
			Footer: map[string][]string{},
		},
		{
			Type:             "feat",
			Title:            "allow provided config object to extend other configs",
			IsBreakingChange: true,
			Footer:           map[string][]string{},
		},
		{
			Type:             "feat",
			Title:            "allow provided config object to extend other configs",
			IsBreakingChange: true,
			Scope:            "api",
			Footer:           map[string][]string{},
		},
		{
			Type:  "feat",
			Title: "allow provided config object to extend other configs",
			Scope: "api",
			Body: `Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.`,
			Footer: map[string][]string{},
		},
		{
			Type:  "feat",
			Title: "allow provided config object to extend other configs",
			Scope: "api",
			Body: `Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.`,
			Footer: map[string][]string{},
		},
		{
			Type:  "feat",
			Title: "allow provided config object to extend other configs",
			Scope: "api",
			Body: `Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.`,
			Footer: map[string][]string{
				"Reviewed-by": {"Z"},
				"Refs":        {"#123"},
			},
		},
		{
			Type:             "feat",
			Title:            "allow provided config object to extend other configs",
			IsBreakingChange: true,
			Scope:            "api",
			Body: `Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.`,
			Footer: map[string][]string{
				"BREAKING CHANGE": {"`extends` key in config file is now used for extending other config files"},
			},
		},
		{
			Type:             "feat",
			Title:            "allow provided config object to extend other configs",
			IsBreakingChange: true,
			Scope:            "api",
			Body: `Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.`,
			Footer: map[string][]string{
				"BREAKING CHANGE": {"`extends` key in config file is now used for extending other config files"},
				"Reviewed-by":     {"Z"},
				"Refs":            {"#123"},
			},
		},
		// --- EXPECTED RESULTS FOR NEW TEST CASES ---
		{
			Type:  "fix",
			Title: "handle new separator and repeated keys",
			Scope: "parser",
			Body: `This commit adds support for the ' #' separator and ensures
that repeated footer keys like 'Closes' are aggregated.`,
			Footer: map[string][]string{
				"Closes": {"#101", "102", "103"},
			},
		},
		{
			Type:   "refactor",
			Title:  "simplify endpoint logic",
			Scope:  "api-v2",
			Footer: map[string][]string{"Reviewed-by": {"Z"}},
		},
		{
			Type:             "chore",
			Title:            "drop support for Node 12",
			IsBreakingChange: true,
			Footer: map[string][]string{
				"BREAKING CHANGE": {`Dropped support for Node 12.

This is the second line of the breaking change description.
All consumers must upgrade to Node 14 or higher.`},
			},
		},
	}

	for i, commitMessage := range commitMessages {
		result, err := convcommit.Parse(commitMessage)

		assert.NoError(t, err)
		assert.Equal(t, expected[i], result, "Failed on test case %d", i)
	}
}

func TestParse_ErrorCases(t *testing.T) {
	testCases := map[string]string{
		"Empty commit":           "",
		"Whitespace only commit": "   ",
		"No colon in header":     "feat a new feature",
		"Random text":            "this is not a commit message",
	}

	for name, commitMessage := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := convcommit.Parse(commitMessage)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}
