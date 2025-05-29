package service

import (
	"ImportAndSearchCsvFile/pkg/models"
	"encoding/csv"
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
	// I decided to use email as key for fast lookups
	mu    sync.RWMutex
	users map[string]models.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]models.User),
	}
}

func (s *UserStore) ImportUsers(reader io.Reader) error {
	r := csv.NewReader(reader)
	r.TrimLeadingSpace = true
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	for i, h := range headers {
		headers[i] = strings.ToLower(strings.ReplaceAll(h, " ", "_"))
	}

	newUsers := make(map[string]models.User)
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

	if len(headers) != len(record) {
		return user, fmt.Errorf("header/record length mismatch: got %d headers and %d fields", len(headers), len(record))
	}

	fieldMap := make(map[string]string)
	for i, header := range headers {
		fieldMap[header] = strings.TrimSpace(record[i])
	}

	var err error

	user.ID, err = strconv.Atoi(fieldMap["id"])
	if err != nil {
		return user, fmt.Errorf("invalid id '%s': %v", fieldMap["id"], err)
	}

	user.FirstName = fieldMap["first_name"]
	user.LastName = fieldMap["last_name"]
	user.Email = fieldMap["email_address"]

	if user.CreatedAt, err = parseTimestamp(fieldMap["created_at"]); err != nil {
		return user, fmt.Errorf("invalid created_at: %v", err)
	}
	if user.DeletedAt, err = parseTimestamp(fieldMap["deleted_at"]); err != nil {
		return user, fmt.Errorf("invalid deleted_at: %v", err)
	}
	if user.MergedAt, err = parseTimestamp(fieldMap["merged_at"]); err != nil {
		return user, fmt.Errorf("invalid merged_at: %v", err)
	}
	if user.ParentUserID, err = parseInt(fieldMap["parent_user_id"]); err != nil {
		return user, fmt.Errorf("invalid parent_user_id: %v", err)
	}

	if err := user.Validate(); err != nil {
		return user, fmt.Errorf("validation error: %w", err)
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

	if strings.ContainsAny(s, "Ee") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return int64(f), nil
	}

	return strconv.ParseInt(s, 10, 64)
}
