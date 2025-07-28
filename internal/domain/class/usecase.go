package class

// Usecase defines the business logic operations for class management
type Usecase interface {
	GetAllClasses() ([]Class, error)
	GetClass(id uint) (*Class, error)
	CreateClass(class *Class) error
	UpdateClass(class *Class) error
	DeleteClass(id uint) error
}
