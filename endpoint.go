package fundconnext

import (
	"errors"
)

const (
	// FundConnextDemoPath is a Path for Demo
	fundConnextDemoPath string = "https://demo.fundconnext.com"
	// FundConnextStagingPath is a Path for Staging
	fundConnextStagingPath string = "https://stage.fundconnext.com"
	// FundConnextProductionPath is a Path for Production
	fundConnextProductionPath string = "https://www.fundconnext.com"
)

// Endpoint is
func endpoint(env string, p string) (string, error) {
	switch env {
	case "production":
		return fundConnextProductionPath + "/" + p, nil
	case "stage":
		return fundConnextStagingPath + "/" + p, nil
	case "demo":
		return fundConnextDemoPath + "/" + p, nil
	default:
		return "", errors.New("Environment is not defined")
	}
}
