package service

import (
	"ImportAndSearchCsvFile/pkg/models"

	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Service interface {
	ImportUsers(reader io.Reader) error
	GetUserByEmail(email string) (models.User, bool)
}

type UserStore struct {
	mu sync.RWMutex
	// I decided to use email as key for fast lookups
	users map[string]models.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]models.User),
	}
}

func (s *UserStore) ImportUsers(reader io.Reader) error {
	r := csv.NewReader(reader)
	r.Comma = '\t'
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	for i, h := range headers {
		headers[i] = strings.ToLower(strings.ReplaceAll(h, " ", "_"))
	}

	var newUsers = make(map[string]models.User)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Skipping malformed line: %v", err)
			continue
		}

		user, err := parseUser(headers, record)
		if err != nil {
			log.Printf("Skipping invalid user record: %v", err)
			continue
		}

		newUsers[user.Email] = user
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = newUsers

	return nil
}

func (s *UserStore) GetUserByEmail(email string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[email]
	return user, exists
}

func parseUser(headers []string, record []string) (models.User, error) {
	var user models.User
	var err error

	if len(record) != len(headers) {
		return user, errors.New("record length doesn't match headers")
	}

	fieldMap := make(map[string]string)
	for i, header := range headers {
		fieldMap[header] = record[i]
	}

	user.ID, err = parseInt(fieldMap["id"])
	if err != nil {
		return user, fmt.Errorf("invalid id: %v", err)
	}

	user.FirstName = fieldMap["first_name"]
	user.LastName = fieldMap["last_name"]
	user.Email = fieldMap["email_address"]

	user.CreatedAt, err = parseTimestamp(fieldMap["created_at"])
	if err != nil {
		return user, fmt.Errorf("invalid created_at: %v", err)
	}

	user.DeletedAt, err = parseTimestamp(fieldMap["deleted_at"])
	if err != nil {
		return user, fmt.Errorf("invalid deleted_at: %v", err)
	}

	user.MergedAt, err = parseTimestamp(fieldMap["merged_at"])
	if err != nil {
		return user, fmt.Errorf("invalid merged_at: %v", err)
	}

	user.ParentUserID, err = parseInt(fieldMap["parent_user_id"])
	if err != nil {
		return user, fmt.Errorf("invalid parent_user_id: %v", err)
	}

	if err := user.Validate(); err != nil {
		return user, fmt.Errorf("Validation error:", err)
	}

	return user, nil
}

func parseInt(s string) (int, error) {
	if s == "-1" {
		return -1, nil
	}
	return strconv.Atoi(s)
}

func parseTimestamp(s string) (int64, error) {
	if s == "-1" {
		return -1, nil
	}

	if strings.Contains(s, "E") || strings.Contains(s, "e") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return int64(f), nil
	}

	return strconv.ParseInt(s, 10, 64)
}
