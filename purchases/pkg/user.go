package templates

type User struct {
	Id int `json:"-" db:"id"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

