// Package footnote implements PHP Markdown Extra footnotes for goldmark,
// together with Pandoc's inline footnotes:
//
//	Here is an inline note.^[Inline notes are easier to write, since you
//	don't have to pick an identifier and move down to type the note.]
//
// It is goldmark's own footnote extension with that syntax added, and it
// replaces that extension rather than adding to it:
//
//	goldmark.New(
//		goldmark.WithExtensions(footnote.NewFootnote(footnote.WithInlineFootnotes())),
//	)
//
// Registering this alongside extension.Footnote would parse every footnote
// twice. Use one or the other.
//
// The exported names are goldmark's, so moving from extension.Footnote is a
// change of import path and nothing else, and so that footnote.go stays a
// copy of goldmark's extension/footnote.go that a future release can be
// merged into. That file was copied at goldmark v1.8.5.
//
// The nodes are goldmark's own, from github.com/yuin/goldmark/extension/ast,
// so anything written against those keeps working.
package footnote
