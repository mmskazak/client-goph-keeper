package commands

import (
	"bytes"
	"client-goph-keerper/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

const Username = "username"
const Password = "password"
const Login = "login"
const Response = "Response: %v\n"

// LoginCommand инициализирует команду для входа пользователя.
func LoginCommand(s *storage.Storage) (*cobra.Command, error) {
	loginCmd := &cobra.Command{
		Use:   Login,
		Short: "Log in a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Проверяем, есть ли уже сохраненный токен
			isHaveToken, err := checkTokenExists(s)
			if err != nil {
				return fmt.Errorf("checkTokenExists: %w", err)
			}
			if isHaveToken {
				return errors.New("пользователь уже авторизован")
			}

			// Получаем логин и пароль из флагов
			username, err := cmd.Flags().GetString(Username)
			if err != nil {
				return fmt.Errorf("cmd.Flags().GetString: %w", err)
			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				return fmt.Errorf("cmd.Flags().GetString: %w", err)
			}

			// Формируем данные для JSON-запроса
			data := map[string]string{
				Username: username,
				Password: password,
			}

			body, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("ошибка кодирования JSON: %w", err)
			}

			// Выполняем HTTP-запрос для входа
			loginURL := path.Join(s.ServerURL, Login)
			req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(body))
			if err != nil {
				return fmt.Errorf("http.NewRequest: %w", err)
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("client.Do: %w", err)
			}
			defer resp.Body.Close() //nolint:errcheck //опустим тут проверку

			// Читаем ответ от сервера
			all, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("all.ReadAll: %w", err)
			}
			fmt.Println(string(all))

			if resp.StatusCode != http.StatusOK {
				fmt.Printf(Response, resp.Status)
				return nil
			}

			// Извлекаем токен из заголовка ответа
			token := resp.Header.Get("Authorization")
			if token == "" {
				return errors.New("токен не найден в заголовке")
			}

			// Сохраняем токен в базе данных
			_, err = s.DataBase.Exec(`
				INSERT OR REPLACE INTO app_params (key, value)
				VALUES (?, ?)`, "jwt_token", token)
			if err != nil {
				return fmt.Errorf("failed to save JWT token: %w", err)
			}
			fmt.Println("JWT токен успешно сохранен:", token)

			fmt.Printf(Response, resp.Status)
			return nil
		},
	}

	// Добавляем обязательные флаги login и password
	loginCmd.Flags().String(Login, "", "Login for the user")
	loginCmd.Flags().String(Password, "", "Password for the user")
	err := loginCmd.MarkFlagRequired(Login)
	if err != nil {
		return nil, fmt.Errorf("failed to mark `login` flag as required: %w", err)
	}
	err = loginCmd.MarkFlagRequired(Password)
	if err != nil {
		return nil, fmt.Errorf("failed to mark `password` flag as required: %w", err)
	}

	return loginCmd, nil
}
