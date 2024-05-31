package jobsummoner

import "context"

type Company struct {
	ID       string
	Name     string
	Url      string
	Avatar   string
	SourceID string
}

type CompanyService interface {
	CreateCompany(ctx context.Context, company Company) (string, error)
}

type CompanyRepository interface {
	AddCompany(ctx context.Context, company Company) (string, error)
	DoesCompanyExist(ctx context.Context, id string) (bool, error)
}
