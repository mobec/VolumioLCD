package display

import "sync"

//IView is the interface for all elements displayed on a screen
type IView interface {
	content() string
	update(deltaTime float64)
}

//==============================================================

//View is the base class for all views
type View struct {
	// mutex protecting the content of the view
	mutex sync.Mutex
}

//==============================================================

//TextView is a view containing only text
type TextView struct {
	View
	text string
}

//Content (text) of any length
func (v *TextView) content() string {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.text
}

//SetText modifies the TextViews content and marks the view as dirty
func (v *TextView) SetText(text string) {
	v.mutex.Lock()
	v.text = text
	v.mutex.Unlock()
}

//Update does nothing for static views
func (v *TextView) update(deltaTime float64) {}

//==============================================================

//ListView concatenates it's child views horizontally
type ListView struct {
	View
	list []IView
}

//content concateneates the ListView's child views
func (v *ListView) content() string {
	var str string
	v.mutex.Lock()
	for _, view := range v.list {
		str += view.content()
	}
	v.mutex.Unlock()
	return str
}

//SetList allows to modify the views elements after construction
func (v *ListView) SetList(list []IView) {
	v.mutex.Lock()
	v.list = list
	v.mutex.Unlock()
}

//update does nothing for static views
func (v *ListView) update(deltaTime float64) {}

//==============================================================

//ScrollView is a view with a fixed length in which child views
// of larger lengths can scroll through
type ScrollView struct {
	View
	child        IView
	length       int
	position     int
	reverse      bool
	childContent string
	time         float64
	speed        float64
}

//Content gathers the child's content and masks it to the current
// scroll window
func (v *ScrollView) content() string {
	// reset the position of the window if the child's content
	// has changed
	v.mutex.Lock()
	newContent := v.child.content()
	length := v.length
	v.mutex.Unlock()

	//reset the scrollview if child content has changed
	if v.childContent != newContent {
		v.position = 0
		v.childContent = newContent
	}

	//pad the content if it is too small
	for i := len(v.childContent); i < length; i++ {
		v.childContent = v.childContent + " "
	}

	// reverse if window reached the limits of the childs content
	if !v.reverse && v.position+v.length > len(v.childContent) {
		v.reverse = true
		v.position = len(v.childContent) - length
	} else if v.reverse && v.position < 0 {
		v.reverse = false
		v.position = 0
	}

	return v.childContent[v.position : v.position+v.length]
}

// Update the scroll view
func (v *ScrollView) update(deltaTime float64) {
	v.mutex.Lock()
	v.child.update(deltaTime)
	characterDuration := 1.0 / v.speed
	v.mutex.Unlock()

	v.time += deltaTime
	if v.time > 2*characterDuration {
		v.time = 2 * characterDuration
	}
	if v.time > characterDuration {
		v.time -= characterDuration
		if v.reverse {
			v.position--
		} else {
			v.position++
		}
	}
}

//SetChild of the view
func (v *ScrollView) SetChild(child IView) {
	v.mutex.Lock()
	v.child = child
	v.mutex.Unlock()
}

//SetSpeed of scrolling in characters per second
func (v *ScrollView) SetSpeed(speed float64) {
	v.mutex.Lock()
	v.speed = speed
	v.mutex.Unlock()
}

//SetLength of the scrollview window
func (v *ScrollView) SetLength(length int) {
	v.mutex.Lock()
	v.length = length
	v.mutex.Unlock()
}

//==============================================================

//RowView is a view spanning one row of a screen. It is limited
// to the physical length of each row on a screen
type RowView struct {
	View
	child  IView
	length int
}

//Content of the child view. Masked to the length of the row
func (v *RowView) content() string {
	v.mutex.Lock()
	content := v.child.content()
	v.mutex.Unlock()

	for i := len(content); i < v.length; i++ {
		content = content + " "
	}
	return content[:v.length]
}

//update the child view. deltaTime should be reasonably low for
// a smooth update of the views
func (v *RowView) update(deltaTime float64) {
	v.mutex.Lock()
	v.child.update(deltaTime)
	v.mutex.Unlock()
}

//SetChild of the row
func (v *RowView) SetChild(child IView) {
	v.mutex.Lock()
	v.child = child
	v.mutex.Unlock()
}
