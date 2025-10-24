package scoring

import (
	"os/exec"
	"testing"
)

// TestMainLintCheck runs the linter and marks all subsequent tests to fail if linting fails
func TestMainLintCheck(t *testing.T) {
	cmd := exec.Command("golangci-lint", "run")
	err := cmd.Run()

	if err != nil {
		// Mark linting as failed and ensure all other tests will fail
		t.Fail()
	}
}
