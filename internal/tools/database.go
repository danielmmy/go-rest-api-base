package tools

import (
	"gorest/api"
)

type FederationRepository interface {
	Setup() error
	AddFederation(*api.Federation) (int, error)
	GetFederation(int) *api.Federation
	GetFederations() []*api.Federation
	UpdateFederation(*api.Federation) (int, error)
	DeleteFederation(int) (int, error)
}

func NewFederationRepository() (*FederationRepository, error) {
	var repo FederationRepository = new(mockDb)
	if err := repo.Setup(); err != nil {
		ErrorLogger.Println(err.Error())
		return nil, err
	}

	return &repo, nil
}
