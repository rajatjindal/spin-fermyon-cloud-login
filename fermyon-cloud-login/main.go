package main

import (
	"fmt"
	"os"

	"github.com/rajatjindal/spin-fermyon-cloud-login/fermyon-cloud-login/pkg/fermyon"
)

func main() {
	err := login()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func login() error {
	cloudLink := fermyon.GetCloudLink(os.Getenv("environment"))

	code, err := fermyon.GenerateDeviceCode(cloudLink)
	if err != nil {
		return fmt.Errorf("generating device code %w", err)
	}

	apiToken, err := fermyon.LoginWithGithub(cloudLink, os.Getenv("E2E_GH_USERNAME"), os.Getenv("E2E_GH_PASSWORD"))
	if err != nil {
		return fmt.Errorf("login with Github to Fermyon cloud: %w", err)
	}

	err = fermyon.ActivateDeviceCode(cloudLink, apiToken, code.UserCode)
	if err != nil {
		return fmt.Errorf("activating device code: %w", err)
	}

	err = fermyon.CheckDeviceCode(cloudLink, code.DeviceCode)
	if err != nil {
		return fmt.Errorf("checking device code: %w", err)
	}

	return nil
}
