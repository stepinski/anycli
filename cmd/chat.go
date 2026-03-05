package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stepinski/anycli/internal/api"
	"github.com/stepinski/anycli/internal/config"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with your AnythingLLM workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1.loading config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		// 2. creating api client for AnythingLLM
		client := api.New(cfg)

		// 3. get the message to send to RAG
		if len(args) == 0 {
			return fmt.Errorf("message required")
		}
		message := args[0]

		// 4. sending message to RAG
		resp, err := client.Chat(cmd.Context(), message)
		if err != nil {
			return fmt.Errorf("chatting: %w", err)
		}

		// 5. printing the response
		fmt.Println(resp.TextResponse)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
