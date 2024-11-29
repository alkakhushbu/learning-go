package stores

type Database interface {
	Create(string) error
	Update(string) error
	Delete(string) error
}

// func CreateUser(db Database) {
// 	db.Create("Alka")
// }

// func UpdateUser(db Database) {
// 	db.Update("Alka")
// }


// func DeleteUser(db Database) {
// 	db.Delete("Alka")
// }