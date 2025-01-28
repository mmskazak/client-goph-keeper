package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

const File = "file"
const Response = "Response: %v\n"
const Authorization = "Authorization"
const ErrSendRequest = "ошибка отправки запроса: %w"
const FileID = "file_id"

// SetAllFilesCmd создает команду для получения списка всех файлов для пользователя.
func SetAllFilesCmd(s *storage.Storage) (*cobra.Command, error) {
	allFilesCmd := &cobra.Command{
		Use:   "all",
		Short: "List all files for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			reqURL := path.Join(s.ServerURL, File, "all")
			req, err := http.NewRequest(http.MethodGet, reqURL, http.NoBody)
			if err != nil {
				return fmt.Errorf(ErrSendRequest, err)
			}

			req.Header.Set(Authorization, s.Token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf(ErrSendRequest, err)
			}
			defer resp.Body.Close() //nolint:errcheck //опустим здесь ошибку

			fmt.Printf(Response, resp.Status)
			return nil
		},
	}

	return allFilesCmd, nil
}
