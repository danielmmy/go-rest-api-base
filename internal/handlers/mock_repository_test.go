package handlers

import (
	"sort"

	"gorest/api"
	"gorest/internal/tools"
)

type FederationRepositoryMock struct{}

var FederationRepositoryMockReturnCode = 0
var FederationRepositoryMockReturnError error = nil
var FederationRepositoryMockReturnReceivedFed *api.Federation = nil
var federationData = map[int]*api.Federation{
	1: {Id: 1, Owner: "Owner 1"},
	2: {Id: 2, Owner: "Owner 2"},
}

func ResetFederationRepositoryMock() {
	FederationRepositoryMockReturnCode = 0
	FederationRepositoryMockReturnError = nil
	FederationRepositoryMockReturnReceivedFed = nil
	federationData = map[int]*api.Federation{
		1: {Id: 1, Owner: "Owner 1"},
		2: {Id: 2, Owner: "Owner 2"},
	}
}

func NewFederationRepositoryMock() *tools.FederationRepository {
	var repo tools.FederationRepository = new(FederationRepositoryMock)
	return &repo
}

func (db *FederationRepositoryMock) Setup() error {
	return FederationRepositoryMockReturnError
}

func (db *FederationRepositoryMock) AddFederation(federation *api.Federation) (int, error) {
	FederationRepositoryMockReturnReceivedFed = federation
	return FederationRepositoryMockReturnCode, FederationRepositoryMockReturnError
}

func (db *FederationRepositoryMock) GetFederation(id int) *api.Federation {
	federation := federationData[id]
	return federation
}

func (db *FederationRepositoryMock) GetFederations() []*api.Federation {
	feds := make([]*api.Federation, 0, len(federationData))
	for _, fed := range federationData {
		feds = append(feds, fed)
	}

	sort.Slice(feds, func(i, j int) bool {
		return feds[i].Id < feds[j].Id
	})
	return feds
}

func (db *FederationRepositoryMock) UpdateFederation(federation *api.Federation) (int, error) {
	return FederationRepositoryMockReturnCode, FederationRepositoryMockReturnError
}

func (db *FederationRepositoryMock) DeleteFederation(id int) (int, error) {
	return FederationRepositoryMockReturnCode, FederationRepositoryMockReturnError
}
