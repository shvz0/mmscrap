package db

import "fmt"

func Migrate() {
	fmt.Println("Migrating...")
	Db.AutoMigrate(&Article{})
	fmt.Println("Successfuly migrated...")
}
