package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type LogEntry struct {
	Time   string `json:"time"`
	Action string `json:"action"`
}

var logChan = make(chan LogEntry, 5)

func main() {
	go logWorker()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n1. Insert")
		fmt.Println("2. Get")
		fmt.Println("3. Exit")
		fmt.Print("Enter choice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			name := readLine(reader, "Enter Name: ")
			email := readLine(reader, "Enter Email: ")

			if findUser(email) != nil {
				fmt.Println("User with this email already exists.")
				continue
			}

			user := User{
				ID:     rand.Intn(9000) + 1000,
				Name:   name,
				Email:  email,
				Status: "active",
			}
			saveUser(user)

			logChan <- LogEntry{
				Time:   time.Now().Format("2006-01-02 15:04:05"),
				Action: "insert",
			}

			fmt.Println("Successfully saved.")

		case 2:
			email := readLine(reader, "Enter Email: ")
			user := findUser(email)

			if user != nil {
				fmt.Printf("Found: %+v\n", *user)
				logChan <- LogEntry{
					Time:   time.Now().Format("2006-01-02 15:04:05"),
					Action: "get",
				}
			} else {
				fmt.Println("No user found with that email.")
			}

		case 3:
			close(logChan)
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func saveUser(user User) {
	f, _ := os.OpenFile("users.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	data, _ := json.Marshal(user)
	f.Write(append(data, '\n'))
}

func findUser(email string) *User {
	data, err := os.ReadFile("users.txt")
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var u User
		json.Unmarshal([]byte(line), &u)
		if u.Email == email {
			return &u
		}
	}
	return nil
}

func logWorker() {
	for entry := range logChan {
		f, _ := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		data, _ := json.Marshal(entry)
		f.Write(append(data, '\n'))
		f.Close()
	}
}
