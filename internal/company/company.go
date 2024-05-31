package company

import (
	"context"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type DefaultCompanyService struct {
	companyRepository jobsummoner.CompanyRepository
}

func NewDefaultCompanyService(repository jobsummoner.CompanyRepository) *DefaultCompanyService {
	return &DefaultCompanyService{repository}
}

func (c *DefaultCompanyService) CreateCompany(ctx context.Context, company jobsummoner.Company) (string, error) {
	id, err := c.companyRepository.CreateCompany(ctx, company)

	if err != nil {
		return "", errors.Wrap(err, "error creating company in company service")
	}

	return id, nil
}

func (c *DefaultCompanyService) DoesCompanyExist(ctx context.Context, id string) (bool, error) {
	doesCompanyExist, err := c.companyRepository.DoesCompanyExist(ctx, id)

	if err != nil {
		return false, errors.Wrap(err, "error finding company in company service")
	}

	return doesCompanyExist, nil
}
