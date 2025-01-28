package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

// SetDeletePasswordCmd создает команду удаления пароля по ID.
func SetDeletePasswordCmd(s *storage.Storage) (*cobra.Command, error) {
	deletePwdCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a password by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Получаем значение флага
			pwdID, _ := cmd.Flags().GetString(PwdID)

			// Создаем URL для запроса
			url := path.Join(s.ServerURL, Pwd, "delete", pwdID)
			req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
			if err != nil {
				return fmt.Errorf(ErrCreateRequest, err)
			}

			// Устанавливаем заголовки
			req.Header.Set(ContentType, applicationJSON)
			req.Header.Set(Authorization, s.Token)

			// Отправляем запрос
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

	// Определяем флаг
	deletePwdCmd.Flags().String(PwdID, "", "Password entry ID")

	// Устанавливаем обязательный флаг
	err := deletePwdCmd.MarkFlagRequired(PwdID)
	if err != nil {
		return nil, fmt.Errorf("error setting required flag 'pwd_id': %w", err)
	}

	return deletePwdCmd, nil
}
