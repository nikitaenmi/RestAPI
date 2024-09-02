package cu

import (
	"bufio"
	"log"
	"os"
)

func CheckUser(path string, name string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == name {
			return true
		}
	}
	return false

}
