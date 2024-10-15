package Execute

import "os/exec"

type Shell struct {
	IsAvailable bool
	Name        string
}

func NewShell() *Shell {
	return &Shell{
		Name:        "Shell",
		IsAvailable: true,
	}
}

func (c *Shell) Execute(command string) (string, error) {
	// setting
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
