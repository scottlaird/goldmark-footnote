goldmark-footnote
=================

[goldmark](https://github.com/yuin/goldmark)'s footnote extension, with
[Pandoc's inline footnotes](https://pandoc.org/MANUAL.html#footnotes) added:

```markdown
Here is an inline note.^[Inline notes are easier to write, since you don't
have to pick an identifier and move down to type the note.]
```

goldmark does not support this syntax and
[has declined to](https://github.com/yuin/goldmark/issues/47), on the grounds
that it can be done as an extension instead. This is that extension.

> **This replaces `extension.Footnote`. It does not add to it.**
> Registering both parses every footnote twice. Use one or the other.

Installation
------------

```
go get github.com/scottlaird/goldmark-footnote
```

Usage
-----

```go
md := goldmark.New(
	goldmark.WithExtensions(
		footnote.NewFootnote(footnote.WithInlineFootnotes()),
	),
)
```

Without `footnote.WithInlineFootnotes()` this behaves exactly like
`extension.Footnote`, because it is that extension: `^[` is ordinary text in
PHP Markdown Extra, so the syntax is opt-in.

Inline footnotes
----------------

An inline footnote is both the reference and the definition, written where it
is referenced rather than in a separate `[^ref]:` block. The two kinds mix
freely and share one numbering, in the order they are referenced:

```markdown
An inline note^[written here] and a reference note[^r] are numbered together.

[^r]: Written down here.
```

The contents are parsed as inline text, so a note cannot contain block elements
such as lists. Nesting, backslash escapes and code spans all work when finding
the closing bracket, and a note may span the lines of a paragraph:

```markdown
Brackets may be nested^[like [this] one], escaped^[like this \] one] or held
in a code span^[like this `]` one].
```

An unclosed `^[` is left as written, as is an escaped `\^[`.

Options
-------

The options are goldmark's, plus one:

| Option | Type | Description |
| ------ | ---- | ----------- |
| `footnote.WithInlineFootnotes` | `-` | Enables Pandoc style inline footnotes. Disabled by default. |
| `footnote.WithFootnoteIDPrefix` | `[]byte \| string` | A prefix for the id attributes. |
| `footnote.WithFootnoteIDPrefixFunction` | `func(gast.Node) []byte` | Determines the id attribute for a given node. |
| `footnote.WithFootnoteLinkTitle` | `[]byte \| string` | An optional title attribute for footnote links. |
| `footnote.WithFootnoteBacklinkTitle` | `[]byte \| string` | An optional title attribute for footnote backlinks. |
| `footnote.WithFootnoteLinkClass` | `[]byte \| string` | A class for footnote links. Defaults to `footnote-ref`. |
| `footnote.WithFootnoteBacklinkClass` | `[]byte \| string` | A class for footnote backlinks. Defaults to `footnote-backref`. |
| `footnote.WithFootnoteBacklinkHTML` | `[]byte \| string` | HTML content for footnote backlinks. Defaults to `&#x21a9;&#xfe0e;`. |

Sidenotes
---------

The nodes this produces are goldmark's own, from
`github.com/yuin/goldmark/extension/ast`, so anything written against those
keeps working. [goldmark-sidenote](https://github.com/scottlaird/goldmark-sidenote)
renders these notes in the margin, inline footnotes included:

```go
goldmark.WithExtensions(
	footnote.NewFootnote(footnote.WithInlineFootnotes()),
	sidenote.New(),
)
```

Relationship to goldmark
------------------------

`footnote.go` is a copy of goldmark's `extension/footnote.go`, taken at
**v1.8.5**, with inline footnote support added. The exported names are
goldmark's, unchanged and stuttering slightly as a result, for two reasons:
moving from `extension.Footnote` is then a change of import path and nothing
else, and a later goldmark release can be merged in without renaming anything
first.

That file has changed twice since 2022, so keeping up with it is not much work.

License
-------

MIT. `footnote.go` is derived from goldmark, copyright Yusuke Inuzuka; see
[LICENSE](LICENSE).
