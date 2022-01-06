package command

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestCompletion(t *testing.T) {
	list := []string{"bash", "zsh", "fish", "powershell"}

	for _, name := range list {
		err := completionCmd.RunE(&cobra.Command{}, []string{name})

		assertNil(t, err, "Completion error for "+name)
		if fakeBuffer.Len() == 0 {
			t.Errorf("Expected completion result for %s, got empty", name)
		}
	}
}

func TestCompletion_WrongArgument(t *testing.T) {
	err := completionCmd.RunE(&cobra.Command{}, []string{"invalidArg"})

	assertErrorEquals(t, "invalid shell name provided: invalidArg", err)
}
