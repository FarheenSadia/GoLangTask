package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

type User struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

const (
	userFile  = "user.data"
	auditFile = "audit.data"
)

func main() {

	if _, err := os.Stat(userFile); os.IsNotExist(err) {
		_, _ = os.Create(userFile)
	}

	for {
		fmt.Print("Choose action: (1) Insert  (2) Retrieve  (3) Exit : ")
		choice := readLine()
		switch choice {
		case "1":
			insertUser()
		case "2":
			retrieveUser()
		case "3":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func insertUser() {
	fmt.Print("Enter name: ")
	name := readLine()
	if strings.ToLower(name) == "exit" {
		return
	}

	fmt.Print("Enter email: ")
	email := readLine()

	user := User{
		Id:     uint(rand.IntN(1000)),
		Name:   name,
		Email:  email,
		Status: "Active",
	}

	data, _ := json.Marshal(user)
	data = append(data, '\n')

	f, _ := os.OpenFile(userFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	f.Write(data)

	writeAudit(user, "Insert")

	fmt.Println("User added successfully.")
}

func retrieveUser() {
	fmt.Print("Enter email to search: ")
	email := readLine()

	file, err := os.Open(userFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		var u User
		if err := json.Unmarshal(scanner.Bytes(), &u); err != nil {
			continue
		}
		if strings.EqualFold(u.Email, email) {
			fmt.Printf("Found User: %+v\n", u)
			writeAudit(u, "Retrieve")
			found = true
			break
		}
	}

	if !found {
		fmt.Println("No user with that email.")
	}
}

func writeAudit(u User, op string) {
	f, _ := os.OpenFile(auditFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	log := fmt.Sprintf("[%s] %s%%!(EXTRA string=%s, string=%s)\n",
		time.Now().Format("2006-01-02 15:04:05"), u.Name, u.Email, op)
	f.WriteString(log)
}

func readLine() string {
	in := bufio.NewReader(os.Stdin)
	text, _ := in.ReadString('\n')
	return strings.TrimSpace(text)
}
