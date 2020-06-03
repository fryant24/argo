package apiclient

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"git.code.oa.com/henrylwang/argo/persist/sqldb"
	"git.code.oa.com/henrylwang/argo/pkg/apiclient/cronworkflow"
	workflowpkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflow"
	workflowarchivepkg "git.code.oa.com/henrylwang/argo/pkg/apiclient/workflowarchive"
	"git.code.oa.com/henrylwang/argo/pkg/apiclient/workflowtemplate"
	"git.code.oa.com/henrylwang/argo/pkg/client/clientset/versioned"
	"git.code.oa.com/henrylwang/argo/server/auth"
	cronworkflowserver "git.code.oa.com/henrylwang/argo/server/cronworkflow"
	workflowserver "git.code.oa.com/henrylwang/argo/server/workflow"
	workflowtemplateserver "git.code.oa.com/henrylwang/argo/server/workflowtemplate"
	"git.code.oa.com/henrylwang/argo/util/help"
)

type argoKubeClient struct {
}

func newArgoKubeClient(clientConfig clientcmd.ClientConfig) (context.Context, Client, error) {
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	wfClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	gatekeeper := auth.NewGatekeeper(auth.Server, wfClient, kubeClient, restConfig)
	ctx, err := gatekeeper.Context(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return ctx, &argoKubeClient{}, nil
}

func (a *argoKubeClient) NewWorkflowServiceClient() workflowpkg.WorkflowServiceClient {
	return &argoKubeWorkflowServiceClient{workflowserver.NewWorkflowServer(sqldb.ExplosiveOffloadNodeStatusRepo)}
}

func (a *argoKubeClient) NewCronWorkflowServiceClient() cronworkflow.CronWorkflowServiceClient {
	return &argoKubeCronWorkflowServiceClient{cronworkflowserver.NewCronWorkflowServer()}
}
func (a *argoKubeClient) NewWorkflowTemplateServiceClient() workflowtemplate.WorkflowTemplateServiceClient {
	return &argoKubeWorkflowTemplateServiceClient{workflowtemplateserver.NewWorkflowTemplateServer()}
}

func (a *argoKubeClient) NewArchivedWorkflowServiceClient() (workflowarchivepkg.ArchivedWorkflowServiceClient, error) {
	return nil, fmt.Errorf("it is impossible to interact with the workflow archive if you are not using the Argo Server, see " + help.CLI)
}
