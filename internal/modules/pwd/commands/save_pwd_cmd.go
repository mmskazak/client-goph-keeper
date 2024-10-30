package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var savePwdCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a password",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		login, _ := cmd.Flags().GetString("login")
		password, _ := cmd.Flags().GetString("password")
		token, _ := cmd.Flags().GetString("token")

		data := map[string]interface{}{
			"title":       title,
			"description": description,
			"credentials": map[string]string{
				"login":    login,
				"password": password,
			},
		}

		body, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("ошибка кодирования JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/pwd/save", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Printf("Response: %v\n", resp.Status)
		return nil
	},
}

func InitSavePwdCmd() *cobra.Command {
	savePwdCmd.Flags().String("title", "", "Title for the password entry")
	savePwdCmd.Flags().String("description", "", "Description for the password entry")
	savePwdCmd.Flags().String("login", "", "Login for the password entry")
	savePwdCmd.Flags().String("password", "", "Password for the password entry")
	savePwdCmd.Flags().String("token", "", "Bearer token for authentication")
	savePwdCmd.MarkFlagRequired("title")
	savePwdCmd.MarkFlagRequired("login")
	savePwdCmd.MarkFlagRequired("password")
	savePwdCmd.MarkFlagRequired("token")
	return savePwdCmd
}
