package util

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func ComparePassword(pwdHash, comPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(comPwd))
	return err == nil
}
