package template

import (
	"context"
	"fmt"
	"log"

	"git.code.oa.com/henrylwang/pkg/errors"

	"github.com/spf13/cobra"

	"git.code.oa.com/henrylwang/argo/cmd/argo/commands/client"
	workflowtemplatepkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflowtemplate"
)

// NewDeleteCommand returns a new instance of an `argo delete` command
func NewDeleteCommand() *cobra.Command {
	var (
		all bool
	)

	var command = &cobra.Command{
		Use:   "delete WORKFLOW_TEMPLATE",
		Short: "delete a workflow template",
		Run: func(cmd *cobra.Command, args []string) {
			apiServerDeleteWorkflowTemplates(all, args)
		},
	}

	command.Flags().BoolVar(&all, "all", false, "Delete all workflow templates")
	return command
}

func apiServerDeleteWorkflowTemplates(allWFs bool, wfTmplNames []string) {
	ctx, apiClient := client.NewAPIClient()
	serviceClient := apiClient.NewWorkflowTemplateServiceClient()
	namespace := client.Namespace()
	var delWFTmplNames []string
	if allWFs {
		wftmplList, err := serviceClient.ListWorkflowTemplates(ctx, &workflowtemplatepkg.WorkflowTemplateListRequest{
			Namespace: namespace,
		})
		if err != nil {
			log.Fatal(err)
		}
		for _, wfTmpl := range wftmplList.Items {
			delWFTmplNames = append(delWFTmplNames, wfTmpl.Name)
		}

	} else {
		delWFTmplNames = wfTmplNames
	}
	for _, wfTmplNames := range delWFTmplNames {
		apiServerDeleteWorkflowTemplate(serviceClient, ctx, namespace, wfTmplNames)
	}
}

func apiServerDeleteWorkflowTemplate(client workflowtemplatepkg.WorkflowTemplateServiceClient, ctx context.Context, namespace, wftmplName string) {
	_, err := client.DeleteWorkflowTemplate(ctx, &workflowtemplatepkg.WorkflowTemplateDeleteRequest{
		Name:      wftmplName,
		Namespace: namespace,
	})
	if err != nil {
		errors.CheckError(err)
	}
	fmt.Printf("WorkflowTemplate '%s' deleted\n", wftmplName)
}
