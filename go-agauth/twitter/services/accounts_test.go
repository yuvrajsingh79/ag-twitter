package services

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-agauth/twitter/users"
	"github.com/stretchr/testify/assert"
)

func TestAccountService_VerifyCredentials(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/1.1/account/verify_credentials.json", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		assertQuery(t, map[string]string{"include_entities": "false", "include_email": "true"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"name": "AdGreetz AdChef", "id": 1273568267781767168}`)
	})

	client := NewClient(httpClient)
	user, _, err := client.Accounts.VerifyCredentials(&AccountVerifyParams{IncludeEntities: Bool(false), IncludeEmail: Bool(true)})
	expected := &users.User{Name: "AdGreetz AdChef", ID: 1273568267781767168}
	assert.Nil(t, err)
	assert.Equal(t, expected, user)
}
