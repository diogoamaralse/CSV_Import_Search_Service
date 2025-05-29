package service

import (
	"ImportAndSearchCsvFile/pkg/models"
	"strings"
	"testing"
)

func TestUserStore_ImportUsers(t *testing.T) {
	tests := []struct {
		name        string
		csvData     string
		wantErr     bool
		wantUserCnt int
	}{
		{
			name: "valid csv",
			csvData: `id	first_name	last_name	email_address	created_at	deleted_at	merged_at	parent_user_id
8	Hanah	Schmidt	Hanah_Schmidt1965@gmail.edu	1.36122E+12	-1	-1	-1`,
			wantErr:     false,
			wantUserCnt: 1,
		},
		{
			name: "malformed csv - skip bad line",
			csvData: `id	first_name	last_name	email_address	created_at	deleted_at	merged_at	parent_user_id
invalid	data	here
31	Emily	Tamm	EmilyTamm@gmail.edu	1.36137E+12	-1	-1	-1`,
			wantErr:     false,
			wantUserCnt: 1,
		},
		{
			name: "completely invalid csv",
			csvData: `invalid	header	here
invalid	data	here`,
			wantErr:     true,
			wantUserCnt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewUserStore()
			reader := strings.NewReader(tt.csvData)

			err := store.ImportUsers(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				userCount := 0
				store.mu.RLock()
				userCount = len(store.users)
				store.mu.RUnlock()

				if userCount != tt.wantUserCnt {
					t.Errorf("Expected %d users, got %d", tt.wantUserCnt, userCount)
				}
			}
		})
	}
}

func TestUserStore_GetUserByEmail(t *testing.T) {
	store := NewUserStore()
	testUser := models.User{
		ID:        1,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}

	csvData := `id	first_name	last_name	email_address	created_at	deleted_at	merged_at	parent_user_id
1	Test	User	test@example.com	123456789	-1	-1	-1`
	reader := strings.NewReader(csvData)
	if err := store.ImportUsers(reader); err != nil {
		t.Fatalf("Failed to import test data: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		wantUser models.User
		want     bool
	}{
		{"existing user", "test@example.com", testUser, true},
		{"nonexistent user", "nonexistent@example.com", models.User{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, got := store.GetUserByEmail(tt.email)
			if got != tt.want {
				t.Errorf("GetUserByEmail() got = %v, want %v", got, tt.want)
			}
			if got && gotUser != tt.wantUser {
				t.Errorf("GetUserByEmail() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
