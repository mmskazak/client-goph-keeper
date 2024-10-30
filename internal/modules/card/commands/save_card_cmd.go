package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var saveCardCmd = &cobra.Command{
	Use:   "save",
	Short: "Save a card entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		number, _ := cmd.Flags().GetString("number")
		pincode, _ := cmd.Flags().GetString("pincode")
		cvv, _ := cmd.Flags().GetString("cvv")
		expire, _ := cmd.Flags().GetString("expire")

		data := map[string]string{
			"title":       title,
			"description": description,
			"number":      number,
			"pincode":     pincode,
			"cvv":         cvv,
			"expire":      expire,
		}

		body, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("ошибка кодирования JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/card/save", bytes.NewBuffer(body))
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

func InitSaveCardCmd() *cobra.Command {
	saveCardCmd.Flags().String("title", "", "Title of the card entry")
	saveCardCmd.Flags().String("description", "", "Description of the card entry")
	saveCardCmd.Flags().String("number", "", "Card number")
	saveCardCmd.Flags().String("pincode", "", "PIN code")
	saveCardCmd.Flags().String("cvv", "", "CVV")
	saveCardCmd.Flags().String("expire", "", "Expiration date")
	saveCardCmd.MarkFlagRequired("title")
	saveCardCmd.MarkFlagRequired("description")
	saveCardCmd.MarkFlagRequired("number")
	saveCardCmd.MarkFlagRequired("pincode")
	saveCardCmd.MarkFlagRequired("cvv")
	saveCardCmd.MarkFlagRequired("expire")
	return saveCardCmd
}
