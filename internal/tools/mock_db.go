package tools

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"gorest/api"
)

type mockDb struct{}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var federationData = map[int]*api.Federation{
	1: {Id: 1, Owner: "Owner 1"},
	2: {Id: 2, Owner: "Owner 2"},
}

func (db *mockDb) Setup() error {

	if n := r.Intn(10); n == 1 {
		return errors.New("random error")
	}
	return nil
}

func (db *mockDb) AddFederation(federation *api.Federation) (int, error) {
	if _, ok := federationData[federation.Id]; ok {
		return http.StatusBadRequest, fmt.Errorf("federation %d already exists", federation.Id)
	}

	federationData[federation.Id] = federation
	return http.StatusCreated, nil
}

func (db *mockDb) GetFederation(id int) *api.Federation {
	fed := federationData[id]
	return fed
}

func (db *mockDb) GetFederations() []*api.Federation {
	// simulate delay
	ticker := time.NewTicker(1 * time.Second)
	<-ticker.C

	federations := make([]*api.Federation, 0, len(federationData))
	for _, fed := range federationData {
		federations = append(federations, fed)
	}

	sort.Slice(federations, func(i, j int) bool {
		return federations[i].Id < federations[j].Id
	})
	return federations
}

func (db *mockDb) UpdateFederation(federation *api.Federation) (int, error) {
	dbFederation, ok := federationData[federation.Id]
	if !ok {
		return http.StatusNotFound, fmt.Errorf("federation %d not found", federation.Id)
	}

	dbFederation.Owner = federation.Owner
	return http.StatusOK, nil
}

func (db *mockDb) DeleteFederation(id int) (int, error) {
	delete(federationData, id)
	return http.StatusOK, nil
}
