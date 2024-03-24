package api

import (
	"encoding/json"
	"testing"
)

// test Federation marshal
// should marshal type with correct field names
func TestFederationMarshal(t *testing.T) {
	// arrange
	sut := Federation{
		Id:    1,
		Owner: "owner",
	}
	want := `{"id":1,"owner":"owner"}`

	// act
	federationJson, _ := json.Marshal(sut)
	federationStr := string(federationJson)

	// assert
	if federationStr != want {
		t.Fatalf("Federation = %q want %q", federationStr, want)
	}
}
