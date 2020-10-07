package p9

import (
	"image/color"

	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/f32color"
)

type _editor struct {
	th       *Theme
	font     text.Font
	textSize unit.Value
	// Color is the text color.
	color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	hint string
	// HintColor is the color of hint text.
	hintColor color.RGBA
	editor    *widget.Editor
	shaper    text.Shaper
}

func (th *Theme) Editor(editor *widget.Editor) *_editor {
	return &_editor{
		th:        th,
		editor:    editor,
		textSize:  th.textSize,
		color:     th.Colors.Get("Text"),
		shaper:    th.shaper,
		hint:      "hint",
		hintColor: th.Colors.Get("Hint"),
	}
}

func (e *_editor) Font(font string) *_editor {
	for i := range e.th.collection {
		if e.th.collection[i].Font.Typeface == text.Typeface(font) {
			e.font = e.th.collection[i].Font
			break
		}
	}
	return e
}
func (e *_editor) TextScale(scale float32) *_editor {
	e.textSize = e.th.textSize.Scale(scale)
	return e
}
func (e *_editor) Color(color string) *_editor {
	e.color = e.th.Colors.Get(color)
	return e
}
func (e *_editor) Hint(hint string) *_editor {
	e.hint = hint
	return e
}
func (e *_editor) HintColor(color string) *_editor {
	e.hintColor = e.th.Colors.Get(color)
	return e
}

func (e *_editor) Fn(c l.Context) l.Dimensions {
	defer op.Push(c.Ops).Pop()
	macro := op.Record(c.Ops)
	paint.ColorOp{Color: e.hintColor}.Add(c.Ops)
	tl := widget.Label{Alignment: e.editor.Alignment}
	dims := tl.Layout(c, e.shaper, e.font, e.textSize, e.hint)
	call := macro.Stop()
	if w := dims.Size.X; c.Constraints.Min.X < w {
		c.Constraints.Min.X = w
	}
	if h := dims.Size.Y; c.Constraints.Min.Y < h {
		c.Constraints.Min.Y = h
	}
	dims = e.editor.Layout(c, e.shaper, e.font, e.textSize)
	disabled := c.Queue == nil
	if e.editor.Len() > 0 {
		textColor := e.color
		if disabled {
			textColor = f32color.MulAlpha(textColor, 150)
		}
		paint.ColorOp{Color: textColor}.Add(c.Ops)
		e.editor.PaintText(c)
	} else {
		call.Add(c.Ops)
	}
	if !disabled {
		paint.ColorOp{Color: e.color}.Add(c.Ops)
		e.editor.PaintCaret(c)
	}
	return dims
}
