package fundconnext

const (
	// FundConnextDemoPath is a Path for Demo
	FundConnextDemoPath string = "https://stage.fundconnext.com"
	// FundConnextStagingPath is a Path for Staging
	FundConnextStagingPath string = "https://demo.fundconnext.com"
	// FundConnextProductionPath is a Path for Production
	FundConnextProductionPath string = "https://www.fundconnext.com"
)

// Endpoint is
func Endpoint() string {
	return FundConnextDemoPath
}
