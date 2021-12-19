package system

type executableCommand interface {
	Run() error
	Output() ([]byte, error)
}

type CommandExecutionerInterface interface {
	RunCommand(command string, args []string) error
	RunCommandForResult(command string, args []string) ([]byte, error)
}

type defaultCommandExecutioner struct {
	builder commandBuilderInterface
}

func InitCommandExecutioner(builder commandBuilderInterface) CommandExecutionerInterface {
	return &defaultCommandExecutioner{
		builder: builder,
	}
}

func (c defaultCommandExecutioner) RunCommand(command string, args []string) error {
	cmd := c.builder.buildInteractiveCommand(command, args)

	return cmd.Run()
}

func (c defaultCommandExecutioner) RunCommandForResult(command string, args []string) ([]byte, error) {
	cmd := c.builder.buildCommand(command, args)

	return cmd.Output()
}
