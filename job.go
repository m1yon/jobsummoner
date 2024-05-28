package jobsummoner

type JobService interface {
	GetJobs() []Job
}

type DefaultJobService struct{}

func NewDefaultJobService() *DefaultJobService {
	return &DefaultJobService{}
}

type Job struct {
	Name string
}

func (j *DefaultJobService) GetJobs() []Job {
	return []Job{
		{"Software Engineer"},
		{"Manager"},
	}
}
