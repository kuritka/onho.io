package services


type IServiceRunner interface {
	UsingPostgre(connectionString string) (error, *Runner)
	UsingRabbit(connectionString string) *Runner
    AsCronJob(durationms int) *Run
    WithName() *Run
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

func (r *Runner) UsingPostgre(connectionString string) *Runner {
	return r
}

func (r *Runner) UsingRabbit(connectionString string) *Runner {
	return r
}

func (r *Runner) AsCronJob(durationms int) *Run {
	return &Run{ runner: r}
}

func (r *Runner) WithName() *Run {
	return &Run{ runner: r}
}

func (r *Run) Run() error {
	return nil
}







