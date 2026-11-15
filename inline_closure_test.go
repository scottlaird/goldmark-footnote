package footnote

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// blockAt returns a block reader over src, split into lines, positioned at
// offset pos.
func blockAt(src []byte, pos int) text.Reader {
	segs := text.NewSegments()
	line, start := 0, 0
	for i := 0; i <= len(src); i++ {
		if i == len(src) || src[i] == '\n' {
			stop := min(i+1, len(src))
			if start < stop {
				if start <= pos {
					line = segs.Len()
				}
				segs.Append(text.NewSegment(start, stop))
			}
			start = i + 1
		}
	}
	r := text.NewBlockReader(src, segs)
	r.SetPosition(line, text.NewSegment(pos, segs.At(line).Stop))
	return r
}

func TestInlineClosesMatchesFindClosure(t *testing.T) {
	const alphabet = "^^[[]]``\\a \n"
	rng := rand.New(rand.NewSource(1))
	for n := 0; n < 3000; n++ {
		b := make([]byte, 1+rng.Intn(60))
		for i := range b {
			b[i] = alphabet[rng.Intn(len(alphabet))]
		}
		src := bytes.TrimLeft(b, "\n")
		var openers []int
		for i := 0; i+1 < len(src); i++ {
			if src[i] == '^' && src[i+1] == '[' {
				openers = append(openers, i)
			}
		}
		if len(openers) == 0 {
			continue
		}

		want := map[int]bool{}
		for _, q := range openers {
			r := blockAt(src, q)
			r.Advance(2)
			_, want[q] = r.FindClosure('[', ']', inlineFootnoteFindClosureOptions)
		}

		// The parser asks in source order; a shuffled order exercises replacing
		// the cached scan.
		shuffled := append([]int(nil), openers...)
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		for _, order := range [][]int{openers, shuffled} {
			pc := parser.NewContext()
			for _, q := range order {
				r := blockAt(src, q)
				l0, s0 := r.Position()
				if got := inlineCloses(r, pc, q+1); got != want[q] {
					t.Fatalf("src %q, opener at %d: closes = %v, want %v", src, q, got, want[q])
				}
				if l, s := r.Position(); l != l0 || s != s0 {
					t.Fatalf("src %q, opener at %d: reader moved", src, q)
				}
			}
		}
	}
}

func BenchmarkInlineFootnoteUnclosed(b *testing.B) {
	md := goldmark.New(goldmark.WithExtensions(NewFootnote(WithInlineFootnotes())))
	src := []byte(strings.Repeat("^[a ", 16000))
	var buf bytes.Buffer
	for range b.N {
		buf.Reset()
		if err := md.Convert(src, &buf); err != nil {
			b.Fatal(err)
		}
	}
}
