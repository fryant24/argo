package commands

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"git.code.oa.com/henrylwang/pkg/cli"
	kubecli "git.code.oa.com/henrylwang/pkg/kube/cli"

	"git.code.oa.com/henrylwang/argo"
	"git.code.oa.com/henrylwang/argo/util"
	"git.code.oa.com/henrylwang/argo/util/cmd"
	"git.code.oa.com/henrylwang/argo/workflow/common"
	"git.code.oa.com/henrylwang/argo/workflow/executor"
	"git.code.oa.com/henrylwang/argo/workflow/executor/docker"
	"git.code.oa.com/henrylwang/argo/workflow/executor/k8sapi"
	"git.code.oa.com/henrylwang/argo/workflow/executor/kubelet"
	"git.code.oa.com/henrylwang/argo/workflow/executor/pns"
)

const (
	// CLIName is the name of the CLI
	CLIName = "argoexec"
)

var (
	clientConfig       clientcmd.ClientConfig
	logLevel           string // --loglevel
	glogLevel          int    // --gloglevel
	podAnnotationsPath string // --pod-annotations
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cli.SetLogLevel(logLevel)
	cli.SetGLogLevel(glogLevel)
}

func NewRootCommand() *cobra.Command {
	var command = cobra.Command{
		Use:   CLIName,
		Short: "argoexec is the executor sidecar to workflow containers",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewInitCommand())
	command.AddCommand(NewResourceCommand())
	command.AddCommand(NewWaitCommand())
	command.AddCommand(cmd.NewVersionCmd(CLIName))

	clientConfig = kubecli.AddKubectlFlagsToCmd(&command)
	command.PersistentFlags().StringVar(&podAnnotationsPath, "pod-annotations", common.PodMetadataAnnotationsPath, "Pod annotations file from k8s downward API")
	command.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "Set the logging level. One of: debug|info|warn|error")
	command.PersistentFlags().IntVar(&glogLevel, "gloglevel", 0, "Set the glog logging level")

	return &command
}

func initExecutor() *executor.WorkflowExecutor {
	config, err := clientConfig.ClientConfig()
	checkErr(err)

	namespace, _, err := clientConfig.Namespace()
	checkErr(err)

	clientset, err := kubernetes.NewForConfig(config)
	checkErr(err)

	podName, ok := os.LookupEnv(common.EnvVarPodName)
	if !ok {
		log.Fatalf("Unable to determine pod name from environment variable %s", common.EnvVarPodName)
	}

	tmpl, err := executor.LoadTemplate(podAnnotationsPath)
	checkErr(err)

	var cre executor.ContainerRuntimeExecutor
	switch os.Getenv(common.EnvVarContainerRuntimeExecutor) {
	case common.ContainerRuntimeExecutorK8sAPI:
		cre, err = k8sapi.NewK8sAPIExecutor(clientset, config, podName, namespace)
	case common.ContainerRuntimeExecutorKubelet:
		cre, err = kubelet.NewKubeletExecutor()
	case common.ContainerRuntimeExecutorPNS:
		cre, err = pns.NewPNSExecutor(clientset, podName, namespace, tmpl.Outputs.HasOutputs())
	default:
		cre, err = docker.NewDockerExecutor()
	}
	checkErr(err)

	wfExecutor := executor.NewExecutor(clientset, podName, namespace, podAnnotationsPath, cre, *tmpl)
	yamlBytes, _ := json.Marshal(&wfExecutor.Template)
	vers := argo.GetVersion()
	log.Infof("Executor (version: %s, build_date: %s) initialized (pod: %s/%s) with template:\n%s", vers, vers.BuildDate, namespace, podName, string(yamlBytes))
	return &wfExecutor
}

// checkErr is a convenience function to panic upon error
func checkErr(err error) {
	if err != nil {
		util.WriteTeriminateMessage(err.Error())
		panic(err.Error())
	}
}
