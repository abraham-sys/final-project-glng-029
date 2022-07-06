package helpers

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) string {
	salt, err := strconv.Atoi(os.Getenv("SALT"))
	fmt.Println(salt)
	if err != nil {
		panic(err)
	}
	password := []byte(p)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, salt)

	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

func CompareHashedPassword(h, p string) bool {
	hashedPassword, password := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	return err == nil
}
