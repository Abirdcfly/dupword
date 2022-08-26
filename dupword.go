// MIT License
//
// Copyright (c) 2022 Abirdcfly
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package dupword defines an Analyzer that checks that duplicate words
// int the source code.
package dupword

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	Doc = `checks for duplicate words in the source code (usually miswritten)

This analyzer checks miswritten duplicate words in comments or package doc or string declaration`
	Message = "Duplicate words found"
)

var (
	noKeyWordErr = errors.New("should have at lease one keyword")
	defaultWord  = []string{"the", "and", "a"}
)

type analyzer struct {
	KeyWord []string
}

func NewAnalyzer(useDefaultWord bool, keyWord ...string) (*analysis.Analyzer, error) {
	var word []string
	if !useDefaultWord {
		if len(keyWord) == 0 {
			return nil, noKeyWordErr
		}
		word = keyWord
	} else {
		word = append(defaultWord, keyWord...)
	}
	a := analyzer{KeyWord: word}
	return &analysis.Analyzer{
		Name:             "dupword",
		Doc:              Doc,
		Requires:         []*analysis.Analyzer{inspect.Analyzer},
		Run:              a.run,
		RunDespiteErrors: true,
	}, nil
}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder(nil, func(n ast.Node) {
		if f, ok := n.(*ast.File); ok {
			a.fixDuplicateWordInComment(pass, f)
		}
		if lit, ok := n.(*ast.BasicLit); ok {
			a.fixDuplicateWordInString(pass, lit)
		}
	})
	return nil, nil
}

func (a *analyzer) fixDuplicateWordInComment(pass *analysis.Pass, f *ast.File) {
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			update, find := a.Check(c.Text)
			if find {
				pass.Report(analysis.Diagnostic{Pos: c.Slash, End: c.End(), Message: Message, SuggestedFixes: []analysis.SuggestedFix{{
					Message: "Update",
					TextEdits: []analysis.TextEdit{{
						Pos:     c.Slash,
						End:     c.End(),
						NewText: []byte(update),
					}},
				}}})
			}
		}
	}
}

func (a *analyzer) fixDuplicateWordInString(pass *analysis.Pass, lit *ast.BasicLit) {
	if lit.Kind != token.STRING {
		return
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		fmt.Printf("lit.Value:%v, err: %v\n", lit.Value, err)
		// fall back to default
		value = lit.Value
	}
	quote := value != lit.Value
	update, find := a.Check(value)
	if quote {
		update = strconv.Quote(update)
	}
	if find {
		pass.Report(analysis.Diagnostic{Pos: lit.Pos(), End: lit.End(), Message: Message, SuggestedFixes: []analysis.SuggestedFix{{
			Message: "Update",
			TextEdits: []analysis.TextEdit{{
				Pos:     lit.Pos(),
				End:     lit.End(),
				NewText: []byte(update),
			}},
		}}})
	}
}

func CheckOneKey(raw, key string) (new string, find bool) {
	if x := strings.Split(raw, key); len(x) < 2 {
		return
	}
	newLine := strings.Builder{}
	wordStart, symbolStart := 0, 0
	curWord, preWord := "", ""
	lastSymbol := ""
	var lastRune int32
	for i, w := range raw {
		if !unicode.IsSpace(w) && unicode.IsSpace(lastRune) {
			symbol := raw[symbolStart:i]
			if !(curWord == key && curWord == preWord) {
				newLine.WriteString(lastSymbol)
				newLine.WriteString(curWord)
			} else {
				find = true
				newLine.WriteString(lastSymbol)
				symbol = ""
			}
			lastSymbol = symbol
			preWord = curWord
			wordStart = i
		} else if unicode.IsSpace(w) && !unicode.IsSpace(lastRune) {
			symbolStart = i
			curWord = raw[wordStart:i]
		} else if i == len(raw)-1 {
			word := raw[wordStart:]
			if !(word == key && word == preWord) {
				newLine.WriteString(lastSymbol)
				newLine.WriteString(word)
			}
		}
		lastRune = w
	}
	if find {
		new = newLine.String()
	}
	return
}

func (a *analyzer) Check(raw string) (update string, find bool) {
	for _, key := range a.KeyWord {
		updateOne, findOne := CheckOneKey(raw, key)
		if findOne {
			raw = updateOne
			find = findOne
			update = updateOne
		}
	}
	return
}
