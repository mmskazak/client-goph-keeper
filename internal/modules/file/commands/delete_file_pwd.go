package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

// SetDeleteFileCmd создает команду для удаления файла по ID.
func SetDeleteFileCmd(s *storage.Storage) (*cobra.Command, error) {
	deleteFileCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a file by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := cmd.Flags().GetString(FileID)
			if err != nil {
				return fmt.Errorf("file id is required: %w", err)
			}

			reqURL := path.Join(s.ServerURL, File, "delete", fileID)
			req, err := http.NewRequest(http.MethodGet, reqURL, http.NoBody)
			if err != nil {
				return fmt.Errorf(ErrSendRequest, err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set(Authorization, s.Token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf(ErrSendRequest, err)
			}
			defer resp.Body.Close() //nolint:errcheck //опустим здесь проверку

			fmt.Printf(Response, resp.Status)
			return nil
		},
	}

	deleteFileCmd.Flags().String(FileID, "", "File ID to delete")
	err := deleteFileCmd.MarkFlagRequired(FileID)
	if err != nil {
		return nil, fmt.Errorf("ошибка установки обязательного флага 'file_id': %w", err)
	}

	return deleteFileCmd, nil
}
