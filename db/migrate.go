package db

import "fmt"

func Migrate() {
	conn := Connect()
	fmt.Println("Migrating...")
	conn.AutoMigrate(&Article{})
	fmt.Println("Successfuly migrated...")
}
