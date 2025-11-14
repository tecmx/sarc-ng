package health

import (
	"fmt"
	"sarc-ng/pkg/rest/client"

	"github.com/spf13/cobra"
)

// NewCommand creates a command to check API health
func NewCommand(clientFactory func() *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check the health status of the SARC API",
		Long:  "Checks if the SARC API server is running and responding to requests.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := clientFactory()

			_, err := client.Health()
			if err != nil {
				fmt.Printf("❌ API is not healthy: %v\n", err)
				return err
			}

			fmt.Println("✅ API is healthy")
			return nil
		},
	}
}
