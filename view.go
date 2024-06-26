package jobsummoner

type HomepageJobModel struct {
	Job
	LastPostedText string
}
type HomepageViewModel struct {
	Jobs            []HomepageJobModel
	Flash           string
	IsAuthenticated bool
}
