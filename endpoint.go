package fundconnext

import (
	"errors"
)

const (
	// FundConnextDemoPath is a Path for Demo
	FundConnextDemoPath string = "https://stage.fundconnext.com"
	// FundConnextStagingPath is a Path for Staging
	FundConnextStagingPath string = "https://demo.fundconnext.com"
	// FundConnextProductionPath is a Path for Production
	FundConnextProductionPath string = "https://www.fundconnext.com"
)

// Endpoint is
func Endpoint(env string, p string) (string, error) {
	switch env {
	case "production":
		return FundConnextProductionPath + "/" + p, nil
	case "stage":
		return FundConnextStagingPath + "/" + p, nil
	case "demo":
		return FundConnextDemoPath + "/" + p, nil
	default:
		return "", errors.New("Environment is not defined")
	}
}
