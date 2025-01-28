package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"

	"github.com/spf13/cobra"
)

// LogoutCommand инициализирует команду для выхода пользователя.
func LogoutCommand(s *storage.Storage) (*cobra.Command, error) {
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out the current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Вызываем функцию для удаления токена из базы данных
			err := s.RemoveTokenFromDB()
			if err != nil {
				return fmt.Errorf("failed to log out: %w", err)
			}

			fmt.Println("You have logged out")
			return nil
		},
	}

	return logoutCmd, nil
}
