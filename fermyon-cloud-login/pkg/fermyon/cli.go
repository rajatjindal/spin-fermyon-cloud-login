package fermyon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

type GetDeviceCodeOutput struct {
	DeviceCode      string `json:"deviceCode"`
	UserCode        string `json:"userCode"`
	VerificationURL string `json:"verificationUrl"`
	ExpiredIn       int    `json:"expiresIn"`
	Interval        int    `json:"interval"`
}

func GenerateDeviceCode(cloudLink string) (*GetDeviceCodeOutput, error) {
	cmd := exec.Command("spin", "login", "--url", cloudLink, "--get-device-code")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := runCmd(cmd)
	if err != nil {
		return nil, err
	}

	dc := &GetDeviceCodeOutput{}
	err = json.Unmarshal(stdout.Bytes(), dc)
	if err != nil {
		return nil, err
	}

	return dc, nil
}

func CheckDeviceCode(cloudLink, deviceCode string) error {
	return runCmd(exec.Command("spin", "login", "--url", cloudLink, "--check-device-code", deviceCode))
}

func runCmd(cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	stdoutWriters := []io.Writer{&stdout}
	stderrWriters := []io.Writer{&stderr}

	if cmd.Stdout != nil {
		stdoutWriters = append(stdoutWriters, cmd.Stdout)
	}

	if cmd.Stderr != nil {
		stderrWriters = append(stderrWriters, cmd.Stderr)
	}

	cmd.Stderr = io.MultiWriter(stderrWriters...)
	cmd.Stdout = io.MultiWriter(stdoutWriters...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running: %s\nstdout:%s\nstderr:%s\n: %w", cmd.String(), stdout.String(), stderr.String(), err)
	}

	return nil
}
