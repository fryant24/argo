package api

import (
	"context"

	"git.code.oa.com/henrylwang/argo/pkg/apiclient/workflow"
	wfv1 "git.code.oa.com/henrylwang/argo/pkg/apis/workflow/v1alpha1"
)

func SubmitWorkflowToAPIServer(apiGRPCClient workflow.WorkflowServiceClient, ctx context.Context, wf *wfv1.Workflow, dryRun bool) (*wfv1.Workflow, error) {

	wfReq := workflow.WorkflowCreateRequest{
		Namespace:    wf.Namespace,
		Workflow:     wf,
		ServerDryRun: dryRun,
	}
	return apiGRPCClient.CreateWorkflow(ctx, &wfReq)
}
