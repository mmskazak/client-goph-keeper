package file

import (
	"client-goph-keerper/internal/modules/file/commands"
	"client-goph-keerper/internal/storage"
	"fmt"

	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Manage files",
}

// InitFileCmd инициализирует команды управления файлами.
func InitFileCmd(s *storage.Storage) (*cobra.Command, error) {
	saveFileCmd, err := commands.SetSaveFileCmd(s)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации команды сохранения файла: %w", err)
	}
	deleteFileCmd, err := commands.SetDeleteFileCmd(s)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации команды удаления файла: %w", err)
	}
	getFileCmd, err := commands.SetGetFileCmd(s)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации команды получения файла: %w", err)
	}
	allFilesCmd, err := commands.SetAllFilesCmd(s)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации команды получения всех файлов: %w", err)
	}

	fileCmd.AddCommand(saveFileCmd)
	fileCmd.AddCommand(deleteFileCmd)
	fileCmd.AddCommand(getFileCmd)
	fileCmd.AddCommand(allFilesCmd)

	return fileCmd, nil
}
