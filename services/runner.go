package services


type IServiceRunner interface {
	UsingPostgre(connectionString string) (error, *Runner)
	UsingRabbit(connectionString string)
	AsWebApp(port int) (error, *Runner)
	AsConsole() (error, *Run)
	Run() error
}

type Runner struct {
}

type Run struct {
	runner *Runner
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) UsingDatabase(connectionString string) *Runner {
	return r
}

func (r *Runner) UsingRabbit(connectionString string) *Runner {
	return r
}

func (r *Runner) AsWebApp(port int) *Run {
	return &Run{ runner: r}
}

func (r *Runner) AsConsole() *Run {
	return &Run{ runner: r}
}

func (r *Run) Run() error {
	return nil
}







