package apiclient

import (
	"context"

	"k8s.io/client-go/tools/clientcmd"

	cronworkflowpkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/cronworkflow"
	workflowpkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflow"
	workflowarchivepkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflowarchive"
	workflowtemplatepkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflowtemplate"
)

type Client interface {
	NewArchivedWorkflowServiceClient() (workflowarchivepkg.ArchivedWorkflowServiceClient, error)
	NewWorkflowServiceClient() workflowpkg.WorkflowServiceClient
	NewCronWorkflowServiceClient() cronworkflowpkg.CronWorkflowServiceClient
	NewWorkflowTemplateServiceClient() workflowtemplatepkg.WorkflowTemplateServiceClient
}

func NewClient(argoServer string, authSupplier func() string, clientConfig clientcmd.ClientConfig) (context.Context, Client, error) {
	if argoServer != "" {
		return newArgoServerClient(argoServer, authSupplier())
	} else {
		return newArgoKubeClient(clientConfig)
	}
}
