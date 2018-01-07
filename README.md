# fuzzy

Fuzzy is a fast and simple go library to perform fuzzy string matching similar to Sublime Text.

Fuzzy can match against single strings, or rank a slice with thousands of strings based on a simple string pattern.
The result not only contains the quality of each individual the match, but also index information that can be used for highlighting matched characters.

The library is unicode-aware and treats input strings with multi-byte characters correctly.

It is also possible to configure the score calculations for different use cases.

## Performance

## Installation

## Credits

The algorithm closely assembles the functionality of Sublime Text's fuzzy search logic.
It is based on the findings of Forrest Smith, who wrote a [blog post](https://blog.forrestthewoods.com/reverse-engineering-sublime-text-s-fuzzy-match-4cffeed33fdb#.d05n81yjy) as well a reference implementation in [C++](https://github.com/forrestthewoods/lib_fts/blob/master/code/fts_fuzzy_match.h).

The library also took some ideas from [sahilm/fuzzy](https://github.com/sahilm/fuzzy), which is another fuzzy-search implementation in go.
