package models

type User struct {
	Name  string
	Email string
	Pay   int
	Bonus int
}

func (u *User) TotalSalary() int {
	if u.Pay == 0 || u.Bonus == 0 {
		return 0
	}
	return u.Pay + u.Bonus
}
