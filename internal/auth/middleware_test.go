package auth

import (
	"net/http"
	"reflect"
	"testing"
)

func TestAddUser(t *testing.T) {
	user := User{ID: 10}

	req := &http.Request{}
	req = AddUser(req, user)

	reqUser := req.Context().Value(UserKey{})
	if !reflect.DeepEqual(user, reqUser) {
		t.Fatalf("Expected user %+v, but got %+v", user, reqUser)
	}
}

func TestGetUser(t *testing.T) {
	user := User{ID: 10}

	req := &http.Request{}
	req = AddUser(req, user)

	reqUser := GetUser(req)
	if !reflect.DeepEqual(user, *reqUser) {
		t.Fatalf("Expected user %+v, but got %+v", user, reqUser)
	}
}
