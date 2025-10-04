package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()

	assert.Equal(t, "gokeenapi", cmd.Use)
	assert.NotNil(t, cmd.PersistentPreRunE)

	// Check that config flag exists
	configFlag := cmd.PersistentFlags().Lookup("config")
	assert.NotNil(t, configFlag)
	assert.Equal(t, "string", configFlag.Value.Type())

	// Check that subcommands are added
	subcommands := cmd.Commands()
	assert.True(t, len(subcommands) > 0)

	// Verify some expected subcommands exist
	commandNames := make(map[string]bool)
	for _, subcmd := range subcommands {
		commandNames[subcmd.Use] = true
	}

	expectedCommands := []string{
		CmdShowInterfaces,
		CmdAddRoutes,
		CmdDeleteRoutes,
		CmdAddDnsRecords,
		CmdDeleteDnsRecords,
		CmdAddAwg,
		CmdUpdateAwg,
		CmdDeleteKnownHosts,
		CmdExec,
	}

	for _, expectedCmd := range expectedCommands {
		assert.True(t, commandNames[expectedCmd], "Expected command %s not found", expectedCmd)
	}
}
