package commands

import (
	"fmt"

	"git.code.oa.com/henrylwang/pkg/cli"
	"github.com/spf13/cobra"

	"git.code.oa.com/henrylwang/argo/cmd/argo/commands/auth"
	"git.code.oa.com/henrylwang/argo/cmd/argo/commands/cron"
	"git.code.oa.com/henrylwang/argo/util/help"

	"git.code.oa.com/henrylwang/argo/cmd/argo/commands/archive"
	"git.code.oa.com/henrylwang/argo/cmd/argo/commands/client"
	"git.code.oa.com/henrylwang/argo/cmd/argo/commands/template"
	"git.code.oa.com/henrylwang/argo/util/cmd"
)

const (
	// CLIName is the name of the CLI
	CLIName = "argo"
)

// NewCommand returns a new instance of an argo command
func NewCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CLIName,
		Short: "argo is the command line interface to Argo",
		Example: fmt.Sprintf(`
If you're using the Argo Server (e.g. because you need large workflow support or workflow archive), please read %s.`, help.CLI),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewCompletionCommand())
	command.AddCommand(NewDeleteCommand())
	command.AddCommand(NewGetCommand())
	command.AddCommand(NewLintCommand())
	command.AddCommand(NewListCommand())
	command.AddCommand(NewLogsCommand())
	command.AddCommand(NewResubmitCommand())
	command.AddCommand(NewResumeCommand())
	command.AddCommand(NewRetryCommand())
	command.AddCommand(NewServerCommand())
	command.AddCommand(NewSubmitCommand())
	command.AddCommand(NewSuspendCommand())
	command.AddCommand(auth.NewAuthCommand())
	command.AddCommand(NewWaitCommand())
	command.AddCommand(NewWatchCommand())
	command.AddCommand(NewTerminateCommand())
	command.AddCommand(archive.NewArchiveCommand())
	command.AddCommand(cmd.NewVersionCmd(CLIName))
	command.AddCommand(template.NewTemplateCommand())
	command.AddCommand(cron.NewCronWorkflowCommand())
	client.AddKubectlFlagsToCmd(command)
	client.AddArgoServerFlagsToCmd(command)

	// global log level
	var logLevel string
	command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cli.SetLogLevel(logLevel)
	}
	command.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "Set the logging level. One of: debug|info|warn|error")

	return command
}
