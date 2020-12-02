package myutil

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// https://github.com/joeshaw/envdecode <- maybe try to use this
func GetEnvOrDefault(envName string, fallback string) string {
	envValue, found := os.LookupEnv(envName)
	if !found {
		return fallback
	}
	return envValue
}

func MissingEnvVariableMsg(envName string) string {
	return fmt.Sprintf("Missing enviroment variable: %v", envName)
}

func ValidateRowsAffected(res sql.Result, w http.ResponseWriter, log *log.Logger) error {
	num, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error with database: %v", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}

	if num == 0 {
		http.Error(w, "No such item", http.StatusBadRequest)
		return err
	}

	return nil
}

func Empty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
