package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

const Pwd = "pwd"
const Response = "Response: %s\n"
const ErrCreateRequest = "ошибка создания запроса: %w"
const ErrSendRequest = "ошибка отправки запроса: %w"
const ContentType = "Content-Type"
const Authorization = "Authorization"
const applicationJSON = "application/json"
const PwdID = "pwd_id"
const Title = "title"
const Login = "login"
const Password = "password"

// SetAllPasswordsCmd создает команду для получения всех паролей пользователя.
func SetAllPasswordsCmd(s *storage.Storage) (*cobra.Command, error) {
	allPwdCmd := &cobra.Command{
		Use:   "all",
		Short: "List all passwords for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Формируем URL для запроса
			url := path.Join(s.ServerURL, Pwd, "all")

			// Создаем запрос
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

			// Чтение тела ответа
			responseData, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("ошибка чтения ответа: %w", err)
			}

			fmt.Printf("Status: %v\n", resp.Status)
			fmt.Printf(Response, responseData)
			return nil
		},
	}

	return allPwdCmd, nil
}
