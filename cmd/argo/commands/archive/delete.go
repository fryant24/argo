package archive

import (
	"fmt"

	"git.code.oa.com/henrylwang/pkg/errors"
	"github.com/spf13/cobra"

	client "git.code.oa.com/henrylwang/argo/cmd/argo/commands/client"
	workflowarchivepkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflowarchive"
)

func NewDeleteCommand() *cobra.Command {
	var command = &cobra.Command{
		Use: "delete UID...",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, apiClient := client.NewAPIClient()
			serviceClient, err := apiClient.NewArchivedWorkflowServiceClient()
			errors.CheckError(err)
			for _, uid := range args {
				_, err = serviceClient.DeleteArchivedWorkflow(ctx, &workflowarchivepkg.DeleteArchivedWorkflowRequest{Uid: uid})
				errors.CheckError(err)
				fmt.Printf("Archived workflow '%s' deleted\n", uid)
			}
		},
	}
	return command
}
