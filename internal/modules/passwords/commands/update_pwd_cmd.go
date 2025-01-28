package commands

import (
	"bytes"
	"client-goph-keerper/internal/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

// SetUpdatePasswordCmd создает команду обновления пароля.
func SetUpdatePasswordCmd(s *storage.Storage) (*cobra.Command, error) {
	updatePwdCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a password by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Получаем значения флагов
			pwdID, err := cmd.Flags().GetString(PwdID)
			if err != nil {
				return fmt.Errorf("error getting pwd_id from flags: %w", err)
			}
			title, err := cmd.Flags().GetString(Title)
			if err != nil {
				return fmt.Errorf("error getting title from flags: %w", err)
			}
			login, err := cmd.Flags().GetString(Login)
			if err != nil {
				return fmt.Errorf("error getting login from flags: %w", err)
			}
			password, err := cmd.Flags().GetString(Password)
			if err != nil {
				return fmt.Errorf("error getting password from flags: %w", err)
			}

			// Формируем JSON-данные для отправки
			data := map[string]interface{}{
				PwdID: pwdID,
				Title: title,
				"credentials": map[string]string{
					Login:    login,
					Password: password,
				},
			}

			body, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("ошибка кодирования JSON: %w", err)
			}

			// Создаем и отправляем запрос
			reqURL := path.Join(s.ServerURL, Pwd, "update")
			req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(body))
			if err != nil {
				return fmt.Errorf(ErrCreateRequest, err)
			}

			req.Header.Set(ContentType, applicationJSON)
			req.Header.Set(Authorization, s.Token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf(ErrSendRequest, err)
			}
			defer resp.Body.Close() //nolint:errcheck //опустим проверку

			fmt.Printf(Response, resp.Status)
			return nil
		},
	}

	updatePwdCmd.Flags().String(PwdID, "", "Password entry ID")
	updatePwdCmd.Flags().String(Title, "", "Title for the password entry")
	updatePwdCmd.Flags().String(Login, "", "Login for the password entry")
	updatePwdCmd.Flags().String(Password, "", "Password for the password entry")

	err := updatePwdCmd.MarkFlagRequired(PwdID)
	if err != nil {
		return nil, fmt.Errorf("error setting required flag 'pwd_id': %w", err)
	}
	err = updatePwdCmd.MarkFlagRequired(Title)
	if err != nil {
		return nil, fmt.Errorf("error setting required flag 'title': %w", err)
	}
	err = updatePwdCmd.MarkFlagRequired(Login)
	if err != nil {
		return nil, fmt.Errorf("error setting required flag 'login': %w", err)
	}
	err = updatePwdCmd.MarkFlagRequired(Password)
	if err != nil {
		return nil, fmt.Errorf("error setting required flag 'password': %w", err)
	}

	return updatePwdCmd, nil
}
