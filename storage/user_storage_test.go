package storage

import (
	"github.com/MerlinDMC/dsapid"
	"os"
	"testing"
)

var (
	storageFilename string = os.TempDir() + "/users.json"
)

func TestNewUserStorage(t *testing.T) {
	storage := NewUserStorage(storageFilename)

	user1 := &dsapid.UserResource{
		Uuid:  "test1",
		Name:  "test1",
		Email: "test1@test.com",
		Token: "test1token",
	}

	user2 := &dsapid.UserResource{
		Uuid:  "test2",
		Name:  "test2",
		Email: "test2@test.com",
		Token: "test2token",
	}

	user3 := &dsapid.UserResource{
		Uuid:  "test3",
		Name:  "test3",
		Email: "test3@test.com",
		Token: "test3token",
	}

	storage.Add(user1.Uuid, user1)
	storage.Add(user2.Uuid, user2)
	storage.Add(user3.Uuid, user3)

	if len(storage.users) != 3 {
		t.Errorf("should have added 3 users but storage has %d", len(storage.users))
	}

	if v := storage.map_name_id["test1"]; v != "test1" {
		t.Errorf("should have mapped name:test1 to test1 but got %s", v)
	}

	if v := storage.map_email_id["test1@test.com"]; v != "test1" {
		t.Errorf("should have mapped email:test1@test.com to test1 but got %s", v)
	}

	if v := storage.map_token_id["test1token"]; v != "test1" {
		t.Errorf("should have mapped token:test1token to test1 but got %s", v)
	}
}

func TestLoadStorage(t *testing.T) {
	storage := NewUserStorage(storageFilename)

	storage.Get("") // dummy get to trigger load()

	if len(storage.users) != 3 {
		t.Errorf("should have added 3 users but storage has %d", len(storage.users))
	}

	if v := storage.map_name_id["test1"]; v != "test1" {
		t.Errorf("should have mapped name:test1 to test1 but got %s", v)
	}

	if v := storage.map_email_id["test1@test.com"]; v != "test1" {
		t.Errorf("should have mapped email:test1@test.com to test1 but got %s", v)
	}

	if v := storage.map_token_id["test1token"]; v != "test1" {
		t.Errorf("should have mapped token:test1token to test1 but got %s", v)
	}
}
