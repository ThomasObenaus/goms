package main

const (
	// PathHealth is the url path for health end-point
	PathHealth = "/health"

	// PathMetrics path for the metrics end-point
	PathMetrics = "/metrics"

	// PathBuildInfo specifies the path where to get the build information about sokar
	PathBuildInfo = "/api/build"

	// PathConfig specifies the path where to get the config information used by sokar
	PathConfig = "/api/config"

	PathCompany      = "/api/company/:id"
	PathCompaniesAll = "/api/companies/list"
	PathCompanies    = "/api/companies"

	PathUser = "/api/user"
)
