package main

// Based on https://github.com/sahilm/fuzzy/blob/master/_example/main.go

import (
	"fmt"
	"io"

	"io/ioutil"
	"strings"

	"time"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/maja42/fuzzy"
)

var filenamesBytes []byte
var err error

var filenames []string

var g *gocui.Gui

func main() {
	filenamesBytes, err = ioutil.ReadFile("ue4_filenames.txt")
	// filenamesBytes, err = ioutil.ReadFile("linux_filenames.txt")
	if err != nil {
		panic(err)
	}

	filenames = strings.Split(string(filenamesBytes), "\n")

	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("inputView", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "< Search Pattern >"
		if _, err := g.SetCurrentView("inputView"); err != nil {
			return err
		}
		v.Editor = gocui.EditorFunc(editor)
	}

	if v, err := g.SetView("results", 0, 3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = true
		v.Title = "< Search Results >"
		updateResults(v, "")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		fallthrough
	case key == gocui.KeySpace:
		v.EditWrite(ch)
		g.Update(func(gui *gocui.Gui) error {
			results, _ := g.View("results")
			updateResults(results, strings.TrimSpace(v.ViewBuffer()))
			return nil
		})

	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		g.Update(func(gui *gocui.Gui) error {
			results, _ := g.View("results")
			updateResults(results, strings.TrimSpace(v.ViewBuffer()))
			return nil
		})
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		g.Update(func(gui *gocui.Gui) error {
			results, _ := g.View("results")
			updateResults(results, strings.TrimSpace(v.ViewBuffer()))
			return nil
		})
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	}
}

func updateResults(v *gocui.View, pattern string) {
	v.Clear()

	t := time.Now()
	matches := fuzzy.Rank(pattern, filenames)
	elapsed := time.Since(t)

	fmt.Fprintf(v, "Searching %d strings in total.\n", len(filenames))
	fmt.Fprintf(v, "Found %d matches in %v\n", len(matches), elapsed)
	for _, match := range matches {
		printMatchedString(v, match)
	}
}

var matchColor = color.New(color.FgRed)

func printMatchedString(w io.Writer, match fuzzy.Match) {
	matchIdx := 0
	for i, r := range match.Str {
		if matchIdx < len(match.MatchedIndexes) && i == match.MatchedIndexes[matchIdx] {
			matchColor.Fprint(w, string(r))
			matchIdx++
		} else {
			fmt.Fprint(w, string(r))
		}
	}
	fmt.Fprintln(w, "")
}
