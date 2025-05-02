package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)

	result := make(DomainStat)
	substr := "." + domain
	for scanner.Scan() {
		var user User
		line := scanner.Bytes()

		if err := user.UnmarshalJSON(line); err != nil {
			return nil, err
		}

		if strings.Contains(user.Email, substr) {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[key]++
		}
	}

	return result, nil
}
