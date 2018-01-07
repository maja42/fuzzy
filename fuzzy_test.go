package fuzzy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var haystack = []string{}

type matchesTestCase struct {
	str string

	score   int
	matches []int
	matched bool
}

func TestSimpleMatches(t *testing.T) {
	testCases := []matchesTestCase{
		matchesTestCase{
			str:     "abcABC",
			score:   firstRuneBonus + 2*sequentialBonus + 3*unmatchedRunePenalty,
			matches: []int{0, 1, 2},
			matched: true,
		},
		matchesTestCase{
			str:     "ABCabc",
			score:   firstRuneBonus + 2*sequentialBonus + 3*unmatchedRunePenalty,
			matches: []int{0, 1, 2},
			matched: true,
		},
		matchesTestCase{
			str:     "aXbXcXXX",
			score:   firstRuneBonus + 5*unmatchedRunePenalty,
			matches: []int{0, 2, 4},
			matched: true,
		},
		matchesTestCase{
			str:     "XaXbcX",
			score:   leadingRunePenalty + sequentialBonus + 3*unmatchedRunePenalty,
			matches: []int{1, 3, 4},
			matched: true,
		},
		matchesTestCase{
			str:     "XXXXabc",
			score:   max(4*leadingRunePenalty, maxLeadingRunePenalty) + 2*sequentialBonus + 4*unmatchedRunePenalty,
			matches: []int{4, 5, 6},
			matched: true,
		},
		matchesTestCase{
			str:     "aabbccXabcXabcX",
			score:   firstRuneBonus + sequentialBonus + 12*unmatchedRunePenalty,
			matches: []int{0, 3, 4},
			matched: true,
		},
		matchesTestCase{
			str:     "Xaabc",
			score:   2*leadingRunePenalty + 2*sequentialBonus + 2*unmatchedRunePenalty,
			matches: []int{2, 3, 4},
			matched: true,
		},
		matchesTestCase{
			str:     "\u0061bc",
			score:   firstRuneBonus + 2*sequentialBonus,
			matches: []int{0, 1, 2},
			matched: true,
		},
		// non-matching:
		matchesTestCase{
			str: "ab",
		},
		matchesTestCase{
			str: "bc",
		},
		matchesTestCase{
			str: "ac",
		},
		matchesTestCase{
			str: "acb",
		},
		matchesTestCase{
			str: "bca",
		},
		matchesTestCase{
			str: "äbc",
		},
	}

	patterns := []string{"abc", "ABC", "Abc", "aBC"}
	for _, pattern := range patterns {
		for _, testCase := range testCases {
			validateMatchesTestCase(t, pattern, testCase)
		}
	}
}

func TestSpecialCharacters(t *testing.T) {
	testCases := []matchesTestCase{
		matchesTestCase{
			str:     "äöüß", // ß is both an uppercase and lowercase character, but we don't want a camelCase bonus for it
			score:   firstRuneBonus + 3*sequentialBonus,
			matches: []int{0, 1, 2, 3},
			matched: true,
		},
		matchesTestCase{
			str:     "ÄÖÜß",
			score:   firstRuneBonus + 3*sequentialBonus,
			matches: []int{0, 1, 2, 3},
			matched: true,
		},
	}

	patterns := []string{"äöüß", "ÄÖÜß", "Äöüß", "äÖÜß"}
	for _, pattern := range patterns {
		for _, testCase := range testCases {
			validateMatchesTestCase(t, pattern, testCase)
		}
	}
}

func TestLeadingAndTrailingSpecialCharacters(t *testing.T) {
	testCase := matchesTestCase{
		str:     "äöüßXäöüß",
		score:   max(maxLeadingRunePenalty, 4*leadingRunePenalty) + 8*unmatchedRunePenalty,
		matches: []int{4},
		matched: true,
	}

	patterns := []string{"x", "X"}
	for _, pattern := range patterns {
		validateMatchesTestCase(t, pattern, testCase)
	}
}

func TestManyMatchesCauseRecursionLimit(t *testing.T) {
	testCase := matchesTestCase{
		str: "XXababababababababab" + // 10x2
			"abababababababababab" + // 10x2
			"abababababababababaB", //10x2
		// --> there would be a CamelCase bonus at the end
		// but due to the recursion limit, we shouldn't reach it
		score:   2*leadingRunePenalty + sequentialBonus + 58*unmatchedRunePenalty,
		matches: []int{2, 3},
		matched: true,
	}
	validateMatchesTestCase(t, "ab", testCase)
}

func TestCamelCaseBonus(t *testing.T) {
	testCase := matchesTestCase{
		str:     "XyababaBab", // Xy is used to avoid leading-letter penalities
		score:   firstRuneBonus + sequentialBonus + camelCaseBonus + 7*unmatchedRunePenalty,
		matches: []int{0, 6, 7}, // use the SequentialBonus and CamelCase bonus - don't take the first match
		matched: true,
	}
	validateMatchesTestCase(t, "Xab", testCase)

	testCase = matchesTestCase{
		str:     "thisIsACamelCaseString",
		score:   max(maxLeadingRunePenalty, 4*leadingRunePenalty) + 4*camelCaseBonus + sequentialBonus + 17*unmatchedRunePenalty,
		matches: []int{4, 6, 7, 12, 16},
		matched: true,
	}
	validateMatchesTestCase(t, "IaCCS", testCase)

	// testCase = matchesTestCase{
	// 	str:   "AbcdaAbcdaAbcda",
	// 	score: firstRuneBonus + 2*camelCaseBonus + sequentialBonus + 11*unmatchedRunePenalty,
	// 	// In this test case, there are two possible results with the exact same score.
	// 	// matches: []int{0, 5, 9, 10},
	// 	matches: []int{0, 4, 5, 10},
	// 	matched: true,
	// }
	// validateMatchesTestCase(t, "AAAA", testCase)
}

func TestPickEarliestMatchOnTie(t *testing.T) {
	testCase := matchesTestCase{
		str:   "xAxBxBxAxBxBxAxAxBxAxB",
		score: leadingRunePenalty + 2*camelCaseBonus + 20*unmatchedRunePenalty,
		// In this test case, there are multiple possible results with the exact same score.
		// In such cases, we prefer results that match as early as possible.
		matches: []int{1, 3},
		matched: true,
	}
	validateMatchesTestCase(t, "AB", testCase)
}

func TestSeparatorBonus(t *testing.T) {
	separators := []string{"/", "-", "_", " ", ".", ",", "\\", "\t"}

	for _, sep := range separators {
		testCase := matchesTestCase{
			str:     "Xaxbaxba" + sep + "baxb",
			score:   firstRuneBonus + sequentialBonus + separatorBonus + 10*unmatchedRunePenalty,
			matches: []int{0, 1, 9},
			matched: true,
		}
		validateMatchesTestCase(t, "Xab", testCase)

		testCase = matchesTestCase{
			str:     "Xaxbaxb" + sep + sep + "a" + sep + sep + "baxb",
			score:   firstRuneBonus + 2*separatorBonus + 13*unmatchedRunePenalty,
			matches: []int{0, 9, 12},
			matched: true,
		}
		validateMatchesTestCase(t, "Xab", testCase)
	}
}

func TestUnicodeLetterCasing(t *testing.T) {
	var cyrillicSmallLetterBe = '\u0431'
	var cyrillicCapitalLetterBe = '\u0411'
	validateLetterCasing(t, cyrillicSmallLetterBe, cyrillicCapitalLetterBe)

	var greekSmallLetterOmegaWithPsiliAndPerispomeniAndYpogegrammeni = '\u1FA6'
	var greekCapitalLetterOmegaWithPsiliAndPerispomeniAndProsgegrammeni = '\u1FAE'
	validateLetterCasing(t, greekSmallLetterOmegaWithPsiliAndPerispomeniAndYpogegrammeni, greekCapitalLetterOmegaWithPsiliAndPerispomeniAndProsgegrammeni)
}

func validateLetterCasing(t *testing.T, smallLetter, captialLetter rune) {
	smallString := string([]rune{'x', smallLetter, 'x', smallLetter, 'x'})
	capitalString := string([]rune{'X', captialLetter, 'X', captialLetter, 'X'})

	testCases := []matchesTestCase{
		matchesTestCase{
			str:     smallString,
			score:   leadingRunePenalty + 3*unmatchedRunePenalty,
			matches: []int{1, 3},
			matched: true,
		},
		matchesTestCase{
			str:     capitalString,
			score:   leadingRunePenalty + 3*unmatchedRunePenalty,
			matches: []int{1, 3},
			matched: true,
		},
	}

	patterns := []string{
		string([]rune{smallLetter, smallLetter}),
		string([]rune{smallLetter, captialLetter}),
		string([]rune{captialLetter, smallLetter}),
		string([]rune{captialLetter, captialLetter}),
	}
	for _, pattern := range patterns {
		for _, testCase := range testCases {
			validateMatchesTestCase(t, pattern, testCase)
		}
	}
}

func TestMatchEmptyStrings(t *testing.T) {
	testCase := matchesTestCase{
		str:     "",
		score:   0,
		matches: []int{},
		matched: true,
	}
	validateMatchesTestCase(t, "", testCase)
}

func TestMatchWithEmptyPattern(t *testing.T) {
	testCase := matchesTestCase{
		str:     "abc",
		score:   3 * unmatchedRunePenalty, // No leading-character penality
		matches: []int{},
		matched: true,
	}
	validateMatchesTestCase(t, "", testCase)
}

func TestMatchAgainstEmptyString(t *testing.T) {
	testCase := matchesTestCase{
		str:     "",
		matched: false,
	}
	validateMatchesTestCase(t, "abc", testCase)
}

func TestRankDiscardsMismatchs(t *testing.T) {
	haystack := []string{"b", "a", "b", "aa", "bb", "aaa", "a", "c"}
	expectedRanking := []string{"a", "a", "aa", "aaa"}

	matches := Rank("a", haystack)
	assert.Len(t, matches, len(expectedRanking))

	for i, expected := range expectedRanking {
		assert.Equal(t, expected, matches[i].Str)
	}
}

func TestRankIsStable(t *testing.T) {
	haystack := []string{"b", "a", "a", "b", "aaa", "aa", "bb", "aaa", "a", "c"}
	expectedIndices := []int{1, 2, 8, 5, 4, 7}

	matches := Rank("a", haystack)
	assert.Len(t, matches, len(expectedIndices))

	for i, expected := range expectedIndices {
		assert.Equal(t, expected, matches[i].Index)
	}
}

func TestRankReversedOrder(t *testing.T) {
	haystack := []string{"aaa", "aa", "a"}
	expectedRanking := []string{"a", "aa", "aaa"}

	matches := Rank("a", haystack)
	assert.Len(t, matches, len(expectedRanking))

	for i, expected := range expectedRanking {
		assert.Equal(t, expected, matches[i].Str)
	}
}

func TestRankPopulatesAllFields(t *testing.T) {
	haystack := []string{"hay", "stack"}

	matches := Rank("a", haystack)
	assert.Len(t, matches, 2)

	assert.Equal(t, "hay", matches[0].Str)
	assert.Equal(t, 0, matches[0].Index)
	assert.Equal(t, leadingRunePenalty+2*unmatchedRunePenalty, matches[0].Score)
	assert.Equal(t, []int{1}, matches[0].MatchedIndexes)

	assert.Equal(t, "stack", matches[1].Str)
	assert.Equal(t, 1, matches[1].Index)
	assert.Equal(t, 2*leadingRunePenalty+4*unmatchedRunePenalty, matches[1].Score)
	assert.Equal(t, []int{2}, matches[1].MatchedIndexes)
}

func TestRankMatchesNothing(t *testing.T) {
	haystack := []string{"x", "xx"}

	matches := Rank("a", haystack)
	assert.Len(t, matches, 0)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func validateMatchesTestCase(t *testing.T, pattern string, testCase matchesTestCase) {
	score, matches, matched := Matches(pattern, testCase.str)

	assert.Equal(t, testCase.matched, matched, "Pattern %q, string %q: 'matched' assertion failed", pattern, testCase.str)
	assert.EqualValues(t, testCase.matches, matches, "Pattern %q, string %q: 'matches' assertion failed", pattern, testCase.str)
	assert.Equal(t, testCase.score, score, "Pattern %q, string %q: 'score' assertion failed", pattern, testCase.str)
}
