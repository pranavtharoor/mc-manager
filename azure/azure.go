package azure

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

func azCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("az")
	cmd.Env = os.Environ()
	cmd.Args = append(cmd.Args, args...)
	cmd.Args = append(cmd.Args, "-o", "json")
	return cmd
}

func azStart(c chan string, args ...string) error {
	defer close(c)
	cmd := azCmd(args...)
	stdout, oErr := cmd.StdoutPipe()
	stderr, eErr := cmd.StderrPipe()
	if oErr != nil || eErr != nil {
		return errors.New("error piping azure command with args: " + strings.Join(args, " "))
	}
	merged := io.MultiReader(stderr, stdout)
	scanner := bufio.NewScanner(merged)
	err := cmd.Start()
	if err != nil {
		return errors.New("error starting azure command with args: " + strings.Join(args, " "))
	}
	for scanner.Scan() {
		out := scanner.Text()
		c <- out
	}
	cmd.Wait()
	return nil
}

func az(result interface{}, args ...string) error {
	cmd := azCmd(args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	if result == nil {
		return nil
	}

	if _, ok := result.(string); ok {
		result = string(out)
		return nil
	}

	return json.Unmarshal(out, result)
}
