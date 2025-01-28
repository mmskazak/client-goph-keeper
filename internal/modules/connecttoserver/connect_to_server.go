package connecttoserver

import (
	"client-goph-keerper/internal/modules/connecttoserver/commands"
	"client-goph-keerper/internal/storage"
	"fmt"

	"github.com/spf13/cobra"
)

var initAppCmd = &cobra.Command{
	Use:   "setServer",
	Short: "Server for job application",
}

// StartsCmd инициализация команд - настроек клиента, для возможности начала работы.
func StartsCmd(s *storage.Storage) (*cobra.Command, error) {
	setServerCmd, err := commands.SetPathRemoteServerCommand(s)
	if err != nil {
		return nil, fmt.Errorf("error setting set server command: %w", err)
	}

	initAppCmd.AddCommand(setServerCmd)

	return initAppCmd, nil
}
