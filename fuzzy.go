package fuzzy

import (
	"sort"
	"unicode"
)

var (
	// Separators represent special characters in search strings.
	// If a rune is matched after such a separator, the resulting score will be better.
	//
	// The separators can be changed to yield optimal results depending on the actual use case.
	// However, they must not be changed while a concurrent match is in progress.
	Separators = []rune{' ', '_', '-', '.', ',', '/', '\\', '\t'}
)

const (
	maxRecursions = 10
)

const (
	sequentialBonus = 15 // Bonus for adjacent matches
	separatorBonus  = 20 // Bonus if a match occurs right after a separator
	camelCaseBonus  = 20 // Bonus if a matched rune is uppercase, while the preceeding rune is lower case
	firstRuneBonus  = 15 // Bonus if the first rune is matched

	leadingRunePenalty    = -5  // Penalty for every rune before the first match
	maxLeadingRunePenalty = -15 // Maximum penalty for leading runes
	unmatchedRunePenalty  = -1  // Penalty for rune that wasn't matched
)

// Match represents a matched string.
type Match struct {
	// The matched string.
	Str string

	// The index of the matched string within the supplied slice.
	Index int

	// The match score which can be used for comparisons with other calls of the same pattern.
	// Higher scores represent a better match. The value itself does not have an intrinsic value
	// and the value range depends on the used pattern.
	Score int

	// Indices of matched runes in Str. Useful for highlighting.
	MatchedIndexes []int
}

type matchSort []Match

func (f matchSort) Len() int           { return len(f) }
func (f matchSort) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f matchSort) Less(i, j int) bool { return f[i].Score > f[j].Score }

/*
Rank the input strings according to their match score in descending order.
If no string matched the pattern, an empty slice is returned.
*/
func Rank(pattern string, strings []string) []Match {
	patternRunes := []rune(pattern)
	res := make([]Match, 0)

	for i, str := range strings {
		strRunes := []rune(str)
		score, matches, matched := match(patternRunes, strRunes, strRunes, make([]int, 0), maxRecursions)
		if matched {
			res = append(res, Match{
				Str:            str,
				Index:          i,
				MatchedIndexes: matches,
				Score:          score,
			})
		}
	}
	sort.Sort(matchSort(res))
	return res
}

/*
Matches the pattern on the given input string.

Returns the score of the match which can be used for comparisons with other calls of the same pattern.
Higher scores represent a better match. The value itself does not have an intrinsic value
and the range depends on the provided pattern.

Returns a slice of indices that represent matched runes in str.
Useful for highlighting. If the pattern did not match, a nil-slice is returned.

The last return value shows if the string matched the given pattern.
*/
func Matches(pattern string, str string) (int, []int, bool) {
	patternRunes := []rune(pattern)
	strRunes := []rune(str)
	return match(patternRunes, strRunes, strRunes, make([]int, 0), maxRecursions)
}

func match(pattern []rune, str []rune, originalStr []rune, allMatches []int, remainingRecursions int) (int, []int, bool) {
	if remainingRecursions < 0 {
		return 0, nil, false
	}

	// The algorithm finds two potential matches:
	//	1) eager matching:
	//		The 'score' and 'matches' if every rune in the pattern is matched as soon as they appear in the search string
	var score int
	var matches = allMatches
	//  2) skip matching:
	//	    Everytime a rune in the pattern is matched, the algorithm also looks for the best match if we would have skipped that specific rune.
	//		This is recursevly done for every matched rune in the input pattern, but only the best "skip-match" is stored.
	var bestSkippingScore int
	var bestSkippingMatches []int
	// Later, the two results are compared and the better one is returned.

	for len(pattern) > 0 && len(str) > 0 {
		if unicode.ToLower(pattern[0]) != unicode.ToLower(str[0]) {
			str = str[1:]
			continue
		}

		skipScore, skipMatches, matched := match(pattern, str[1:], originalStr, matches, remainingRecursions-1)
		if matched { // The pattern matches multiple times
			if bestSkippingMatches == nil || skipScore > bestSkippingScore {
				bestSkippingScore = skipScore
				bestSkippingMatches = skipMatches
			}
		}

		strIdx := len(originalStr) - len(str) // The index within the original, non-sliced string
		matches = append(matches, strIdx)

		pattern = pattern[1:]
		str = str[1:]
	}

	if len(pattern) != 0 { // We couldn't match the whole pattern
		return 0, nil, false
	}

	// Calculate the score of the eager-match

	// Leading rune penality
	if len(matches) > 0 { // no penality if the pattern was an empty string
		penalty := leadingRunePenalty * matches[0]
		if penalty < maxLeadingRunePenalty {
			penalty = maxLeadingRunePenalty
		}
		score += penalty
	}

	// Unmatched rune penalty
	unmatched := len(originalStr) - len(matches)
	score += unmatchedRunePenalty * unmatched

	// Ordering bonuses
	for i, strIdx := range matches {
		if i > 0 {
			// Two runes in a row
			if matches[i-1]+1 == strIdx {
				score += sequentialBonus
			}
		}

		if strIdx > 0 {
			// Camel case
			leftNeighbor := originalStr[strIdx-1]
			current := originalStr[strIdx]

			if isLower(leftNeighbor) && isUpper(current) {
				score += camelCaseBonus
			}

			// Separator
			if isSeparator(leftNeighbor) {
				score += separatorBonus
			}
		} else {
			score += firstRuneBonus
		}
	}

	if bestSkippingMatches != nil && (bestSkippingScore > score) {
		// The skip-score is better than the eager score
		return bestSkippingScore, bestSkippingMatches, true
	}
	return score, matches, true
}

func isLower(r rune) bool {
	return r != unicode.ToUpper(r)
}

func isUpper(r rune) bool {
	return r != unicode.ToLower(r)
}

func isSeparator(r rune) bool {
	for _, s := range Separators {
		if s == r {
			return true
		}
	}
	return false
}
