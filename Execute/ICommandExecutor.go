package Execute

type ICommandExecutor interface {
	Execute(command string) (string, error)
}
