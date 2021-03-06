package commands

import (
	"git.code.oa.com/henrylwang/pkg/errors"
	"github.com/spf13/cobra"

	"git.code.oa.com/henrylwang/argo/cmd/argo/commands/client"
	workflowpkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflow"
)

func NewResubmitCommand() *cobra.Command {
	var (
		memoized      bool
		cliSubmitOpts cliSubmitOpts
	)
	var command = &cobra.Command{
		Use:   "resubmit [WORKFLOW...]",
		Short: "resubmit one or more workflows",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewWorkflowServiceClient()
			namespace := client.Namespace()

			for _, name := range args {
				created, err := serviceClient.ResubmitWorkflow(ctx, &workflowpkg.WorkflowResubmitRequest{
					Namespace: namespace,
					Name:      name,
					Memoized:  memoized,
				})
				errors.CheckError(err)
				printWorkflow(created, cliSubmitOpts.output, DefaultStatus)
				waitOrWatch([]string{created.Name}, cliSubmitOpts)
			}
		},
	}

	command.Flags().StringVarP(&cliSubmitOpts.output, "output", "o", "", "Output format. One of: name|json|yaml|wide")
	command.Flags().BoolVarP(&cliSubmitOpts.wait, "wait", "w", false, "wait for the workflow to complete")
	command.Flags().BoolVar(&cliSubmitOpts.watch, "watch", false, "watch the workflow until it completes")
	command.Flags().BoolVar(&memoized, "memoized", false, "re-use successful steps & outputs from the previous run (experimental)")
	return command
}
