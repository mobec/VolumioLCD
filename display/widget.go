package display

//IWidget is the interface expected by Screen for all widgets
type IWidget interface {
	Text() string
	Length() int
}

//TextField is the most basic widget containing text
type TextField struct {
	text string
}

//Text of the widget
func (w *TextField) Text() string { return w.text }

//SetText of the widget
func (w *TextField) SetText(text string) {
	w.text = text
}

//Length of the TextField in bytes
func (w *TextField) Length() int {
	return len(w.text)
}
