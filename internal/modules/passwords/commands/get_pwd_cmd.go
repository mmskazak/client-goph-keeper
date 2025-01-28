package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

// SetGetPasswordCmd создает команду получения пароля по ID.
func SetGetPasswordCmd(s *storage.Storage) (*cobra.Command, error) {
	getPwdCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a password by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Получаем значение флага
			pwdID, _ := cmd.Flags().GetString(PwdID)

			// Формируем URL запроса
			url := path.Join(s.ServerURL, Pwd, "get", pwdID)

			// Создаем запрос
			req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
			if err != nil {
				return fmt.Errorf(ErrCreateRequest, err)
			}

			req.Header.Set(ContentType, applicationJSON)
			req.Header.Set(Authorization, s.Token)

			// Отправляем запрос
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf(ErrSendRequest, err)
			}
			defer resp.Body.Close() //nolint:errcheck //опустим здесь проверку

			// Чтение ответа
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("ошибка чтения ответа: %w", err)
			}

			fmt.Printf(Response, resp.Status)
			fmt.Printf("Data: %s\n", data)
			return nil
		},
	}

	// Определяем флаги
	getPwdCmd.Flags().String(PwdID, "", "Password entry ID")
	// Устанавливаем обязательные флаги
	err := getPwdCmd.MarkFlagRequired(PwdID)
	if err != nil {
		return nil, fmt.Errorf("error setting required flag 'pwd_id': %w", err)
	}

	return getPwdCmd, nil
}
