# go fuzzy!
[![Build Status](https://travis-ci.org/maja42/fuzzy.svg?branch=master)](https://travis-ci.org/maja42/fuzzy)
[![GoDoc](https://godoc.org/github.com/maja42/fuzzy?status.svg)](https://godoc.org/github.com/maja42/fuzzy)

Fuzzy is a fast and simple go library to perform fuzzy string matching similar to Sublime Text.

Fuzzy can match against single strings, or rank a slice with thousands of strings based on a simple string pattern.
The result not only contains the score of each individual match, but also index information that can be used to highlight matched characters.

The library is unicode-aware and treats input strings with multi-byte characters correctly.

It is also possible to configure the score calculations for different use cases.

## Demo

A demo application that searches ~16k files from the Unreal Engine 4 codebase can be found in the _example folder.

```
cd _example/
go get github.com/fatih/color
go get github.com/jroimartin/gocui
go run main.go
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/maja42/fuzzy"
)

func main() {
	pattern := "re"
	strings := []string{
		"The Seven Samurai", "Bonnie and Clyde", "Reservoir Dogs", "Airplane!", "Pan's Labyrinth", "鋼の錬金術師",
	}

	matches := fuzzy.Rank(pattern, strings)

	for _, match := range matches {
		fmt.Println(match.Str)
		// Prints:
		// 	Reservoir Dogs
		// 	Airplane!
	}
}
```

## Performance

The algorithm is optimized for go. Matching patterns against ~61k file names from the Linux Kernel takes 13ms on an average Laptop.

## Installation

`go get github.com/maja42/fuzzy`

## Credits

The algorithm closely assembles the functionality of Sublime Text's fuzzy search logic.
It is based on the findings of Forrest Smith, who wrote a [blog post](https://blog.forrestthewoods.com/reverse-engineering-sublime-text-s-fuzzy-match-4cffeed33fdb#.d05n81yjy) as well a reference implementation in [C++](https://github.com/forrestthewoods/lib_fts/blob/master/code/fts_fuzzy_match.h).

The library also took some ideas from [sahilm/fuzzy](https://github.com/sahilm/fuzzy), which is another fuzzy-search implementation in go.
