package fingerprinting

import (
	"encoding/json"
	"os/exec"
)

type FingerprintResult struct {
	Duration    float64 `json:"duration"`
	FingerPrint string  `json:"fingerprint"`
}

func Fingerprint(filepath string) (FingerprintResult, error) {
	cmd := exec.Command("fpcalc", "-json", filepath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return FingerprintResult{}, err
	}

	if err := cmd.Start(); err != nil {
		return FingerprintResult{}, err
	}

	fp := FingerprintResult{}

	if err := json.NewDecoder(stdout).Decode(&fp); err != nil {
		return FingerprintResult{}, err
	}
	if err := cmd.Wait(); err != nil {
		return FingerprintResult{}, err
	}
	return fp, nil
}
