package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var saveTextCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a text entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		textContent, _ := cmd.Flags().GetString("text_content")

		data := map[string]string{
			"title":        title,
			"description":  description,
			"text_content": textContent,
		}

		body, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("ошибка кодирования JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/text/save", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		// Замените на действительный токен
		req.Header.Set("Authorization", "Bearer YOUR_TOKEN_HERE")

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

func InitSaveTextCmd() *cobra.Command {
	saveTextCmd.Flags().String("title", "", "Title of the text entry")
	saveTextCmd.Flags().String("description", "", "Description of the text entry")
	saveTextCmd.Flags().String("text_content", "", "Content of the text entry")
	saveTextCmd.MarkFlagRequired("title")
	saveTextCmd.MarkFlagRequired("description")
	saveTextCmd.MarkFlagRequired("text_content")
	return saveTextCmd
}
