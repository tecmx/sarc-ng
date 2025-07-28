package lesson

// Usecase defines the business logic operations for lesson management
type Usecase interface {
	GetAllLessons() ([]Lesson, error)
	GetLesson(id uint) (*Lesson, error)
	CreateLesson(lesson *Lesson) error
	UpdateLesson(lesson *Lesson) error
	DeleteLesson(id uint) error
}
