package tools

import (
	"math/rand"
	"net/http"
	"reflect"
	"testing"

	"gorest/api"
)

// test setup() with random error
// should return error
func TestSetupRandomError(t *testing.T) {
	// arrange
	r = rand.New(rand.NewSource(1))
	sut := new(mockDb)

	// act
	err := sut.Setup()

	// assert
	if err == nil {
		t.Fatal(`setup() = "nil" want "random error"`)
	}

	if err.Error() != "random error" {
		t.Fatalf(`setup() = %q want "random error"`, err)
	}
}

// test setup() without error
// should return nil
func TestSetupSuccess(t *testing.T) {
	// arrange
	r = rand.New(rand.NewSource(2))
	sut := new(mockDb)

	// act
	err := sut.Setup()

	// assert
	if err != nil {
		t.Fatalf(`setup() = %q want "nil"`, err)
	}
}

// test AddFederation(*api.Federation) (int, err) with duplicated federation
// should return error
func TestAddFederationDuplicated(t *testing.T) {
	// arrange
	federation := &api.Federation{
		Id: 1,
	}
	sut := new(mockDb)
	wantErr := "federation 1 already exists"
	wantCode := http.StatusBadRequest

	// act
	code, err := sut.AddFederation(federation)

	// assert
	if code != wantCode {
		t.Fatalf("AddFederation(federation) = %d want %d", code, wantCode)
	}

	if err.Error() != wantErr {
		t.Fatalf("AddFederation(federation) = %q want %q", err, wantErr)
	}
}

// test AddFederation(*api.Federation) (int, err) success
// should create  federation
func TestAddFederationSuccess(t *testing.T) {
	// arrange
	federation := &api.Federation{
		Id:    123,
		Owner: "Federation 123",
	}
	sut := new(mockDb)
	var wantErr error = nil
	wantCode := http.StatusCreated

	// act
	code, err := sut.AddFederation(federation)

	// assert
	if code != wantCode {
		t.Fatalf("AddFederation(federation) = %d want %d", code, wantCode)
	}

	if err != wantErr {
		t.Fatalf("AddFederation(federation) = %v want %v", err, wantErr)
	}

	var federation123 = sut.GetFederation(123)
	if federation123 != federation {
		t.Fatalf("AddFederation(federation) = %v want %v", federation123, federation)
	}
}

// test GetFederation with bad request
// should return nil
func TestGetFederationNotFound(t *testing.T) {
	// arrange
	id := -1
	sut := new(mockDb)

	// act
	fed := sut.GetFederation(id)

	// assert
	if fed != nil {
		t.Fatalf("GetFederation(id) = %v want %v", fed, nil)
	}
}

// test GetFederation without error
// should return federation
func TestGetFederationSuccess(t *testing.T) {
	// arrange
	id := 1
	sut := new(mockDb)

	// act
	fed := sut.GetFederation(id)

	// assert
	if fed.Id != id {
		t.Fatalf("GetFederation(id) = %v want federation: %q", fed, id)
	}
}

// test GetFederation success
// should return an ordered slice of federations
func TestGetFederationsSucess(t *testing.T) {
	// arrange
	federationData = map[int]*api.Federation{
		2: {Id: 2, Owner: "Owner 2"},
		1: {Id: 1, Owner: "Owner 1"},
	}
	sut := new(mockDb)
	dontWant := []*api.Federation{
		federationData[2],
		federationData[1],
	}
	want := []*api.Federation{
		federationData[1],
		federationData[2],
	}

	// act
	federations := sut.GetFederations()

	// assert
	if !reflect.DeepEqual(federations, want) {
		t.Fatalf("GetFederations() = %v want %v", federations, want)
	}

	if reflect.DeepEqual(federations, dontWant) {
		t.Fatalf("GetFederations() = %v want %v", federations, dontWant)
	}
}

// test UpdateFederation(*federation) not found
// should return error
func TestUpdateFederationNotFound(t *testing.T) {
	// arrange
	federation := &api.Federation{Id: -1}
	sut := new(mockDb)
	wantCode := http.StatusNotFound
	wantError := "federation -1 not found"

	// act
	code, err := sut.UpdateFederation(federation)

	// assert
	if code != wantCode {
		t.Fatalf("UpdateFederation(federation) = %d want %d", code, wantCode)
	}

	if err.Error() != wantError {
		t.Fatalf("UpdateFederation(federation) = %q want %q", err, wantError)
	}
}

// test UpdateFederation(*federation) success
// should update federation
func TestUpdateFederationSuccess(t *testing.T) {
	// arrange
	federation := &api.Federation{Id: 1, Owner: "new owner"}
	sut := new(mockDb)
	wantCode := http.StatusOK
	var wantError error = nil

	// act
	code, err := sut.UpdateFederation(federation)

	// assert
	if code != wantCode {
		t.Fatalf("UpdateFederation(federation) = %d want %d", code, wantCode)
	}

	if err != wantError {
		t.Fatalf("UpdateFederation(federation) = %v want %v", err, wantError)
	}

	if federationData[1].Owner != federation.Owner {
		t.Fatalf("UpdateFederation(federation) = %q want %q", federationData[1].Owner, federation.Owner)
	}

	if federationData[1].Id != federation.Id {
		t.Fatalf("UpdateFederation(federation) = %d want %d", federationData[1].Id, federation.Id)
	}
}

// test DeleteFederation(string)
// should deleteFederation
func TestDeleteFederationSucess(t *testing.T) {
	// arrange
	federationData = map[int]*api.Federation{
		1: {Id: 1, Owner: "Owner 1"},
		2: {Id: 20, Owner: "Owner 2"},
	}
	id := 1
	sut := new(mockDb)
	wantCode := http.StatusOK
	var wantError error = nil

	// act
	code, err := sut.DeleteFederation(id)
	_, ok := federationData[id]

	// assert
	if code != wantCode {
		t.Fatalf("DeleteFederation(id) = %d want %d", code, wantCode)
	}

	if err != wantError {
		t.Fatalf("DeleteFederation(id) = %v want %v", err, wantError)
	}

	if ok {
		t.Fatalf("DeleteFederation(id) = %v want %v", ok, false)
	}

}
