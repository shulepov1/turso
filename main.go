package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type NewDatabaseArgs struct {
	Name  string
	Group string
}

func Create(username string) string {
	data := []byte(fmt.Sprintf(`{
		"name": "%s-db",
		"group": "default"
	}`, username))

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.turso.tech/v1/organizations/%s/databases", os.Getenv("ORG_NAME")), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Err when request", err)
	}

	token := os.Getenv("DB_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("err when reading", err)
	}

	log.Println(string(out))
	return ""
}

type FindDatabaseArgs struct {
	token        string
	Organization string
	DatabaseName string
}

func NewFindDatabase(token, org, name string) *FindDatabaseArgs {
	return &FindDatabaseArgs{
		token:        token,
		Organization: org,
		DatabaseName: name,
	}
}

func main() {
	username := "ilya" // уникальный username пользователя

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Err when loading env", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.turso.tech/v1/organizations/%s/databases/%s-db", os.Getenv("ORG_NAME"), username), nil)
	if err != nil {
		log.Fatal("Err when request", err)
	}

	token := os.Getenv("DB_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("err when reading", err)
	}
	if strings.Contains(string(out), "error") {
		Create(username)
	}

	log.Println(string(out))
}
