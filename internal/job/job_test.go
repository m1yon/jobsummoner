package job

// func TestSqliteJobService(t *testing.T) {
// 	t.Run("add job and immediately get added job", func(t *testing.T) {
// 		db := sqlitedb.NewTestDB()
// 		companyRepository := sqlitedb.NewInMemorySqliteCompanyRepository(db)
// 		jobRepository := sqlitedb.NewInMemorySqliteJobRepository(db, companyRepository)
// 		jobService := NewDefaultJobService(jobRepository)

// 	})

// }

// func TestCreateJobs(t *testing.T) {
// 	db := sqlitedb.NewTestDB()
// 	companyRepository := sqlitedb.NewInMemorySqliteCompanyRepository(db)
// 	jobRepository := sqlitedb.NewInMemorySqliteJobRepository(db, companyRepository)
// 	jobService := NewDefaultJobService(jobRepository)

// 	jobs := []jobsummoner.Job{
// 		{Position: "Software Engineer"},
// 		{Position: "Manager"},
// 	}

// 	jobService.AddJobs(jobs)

// 	// assert.Equal(t, []jobsummoner.Job{
// 	// 	{Position: "Software Engineer"},
// 	// 	{Position: "Manager"},
// 	// }, res)
// }
