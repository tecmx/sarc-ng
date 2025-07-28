package class

// Repository defines the data access operations for classes
// All methods are explicitly named with the Class entity
type Repository interface {
	ReadClassList() ([]Class, error)
	ReadClass(id uint) (*Class, error)
	CreateClass(class *Class) error
	UpdateClass(class *Class) error
	DeleteClass(id uint) error
}
