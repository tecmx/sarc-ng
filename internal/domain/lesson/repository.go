package lesson

// Repository defines the data access operations for lessons
// All methods are explicitly named with the Lesson entity
type Repository interface {
	ReadLessonList() ([]Lesson, error)
	ReadLesson(id uint) (*Lesson, error)
	CreateLesson(lesson *Lesson) error
	UpdateLesson(lesson *Lesson) error
	DeleteLesson(id uint) error
}
