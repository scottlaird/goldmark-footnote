package footnote

import (
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var inlineClosuresKey = parser.NewContextKey()

// inlineCloses reports whether FindClosure with inlineFootnoteFindClosureOptions,
// started just after the '[' at source offset opener, would find its closer.
// block must be positioned before that '[' and is left where it was.
//
// A failed search scans to the end of the block, so asking FindClosure at
// every opener is quadratic in a paragraph full of unclosed ones. Instead one
// scan answers for every later opener it saw outside a code span: from there
// on, a search started at that opener sees exactly the same code spans.
func inlineCloses(block text.Reader, pc parser.Context, opener int) bool {
	scan, _ := pc.Get(inlineClosuresKey).(map[int]bool)
	if closes, ok := scan[opener]; ok {
		return closes
	}
	scan = scanInlineClosures(block)
	pc.Set(inlineClosuresKey, scan)
	return scan[opener]
}

// scanInlineClosures follows FindClosure's rules for '[' and ']' from block's
// position to the end of the block, and returns, for each '[' outside a code
// span, whether it is closed.
func scanInlineClosures(block text.Reader) map[int]bool {
	line, pos := block.Position()
	defer block.SetPosition(line, pos)

	closes := map[int]bool{}
	var open []int
	codeSpan := 0
	for {
		bs, seg := block.PeekLine()
		if bs == nil {
			return closes
		}
		for i := 0; i < len(bs); i++ {
			c := bs[i]
			switch {
			case c == '`':
				n := 1
				for i+n < len(bs) && bs[i+n] == '`' {
					n++
				}
				if codeSpan == 0 {
					codeSpan = n
				} else if codeSpan == n {
					codeSpan = 0
				}
				i += n - 1
			case codeSpan != 0:
			case c == '\\' && i < len(bs)-1 && util.IsPunct(bs[i+1]):
				i++
			case c == '[':
				p := seg.Start + i - seg.Padding
				closes[p] = false
				open = append(open, p)
			case c == ']' && len(open) > 0:
				closes[open[len(open)-1]] = true
				open = open[:len(open)-1]
			}
		}
		block.AdvanceLine()
	}
}
