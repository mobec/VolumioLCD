package display

//==============================================================

//Screen contains a collection of views for a physical display
type Screen struct {
	ID   uint32
	rows []RowView
}

//NewScreen for a display with displays physical parameters
func NewScreen(rows int, rowLength int) Screen {
	var s Screen
	s.rows = make([]RowView, rows)
	for i := 0; i < rows; i++ {
		s.rows[i].length = rowLength
	}
	return s
}

//update all the rows and therefore views of the screen
func (s *Screen) update(deltaTime float64) {
	for idx := range s.rows {
		s.rows[idx].update(deltaTime)
	}
}

//GetRow by index. Can be used to modify the contents of the screen
func (s *Screen) GetRow(idx int) *RowView {
	return &s.rows[idx]
}
