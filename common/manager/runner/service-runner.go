package runner

import (
	"context"
	"github.com/kuritka/onho.io/common/manager/depresolver"
	"sync"

	"github.com/kuritka/onho.io/common/log"
	"github.com/kuritka/onho.io/common/utils"
)

var zerologger = log.Logger()

// Runner doc
type Runner struct {
//	grpcRunner     *GrpcRunner
	restRunner     *RestRunner
	//jobRunner      *JobRunner
	//cronRunner     *CronRunner
}

// NewServiceRunner doc
func NewServiceRunner() *Runner {
	return &Runner{}
}


func (r *Runner) WithRest(rr *RestRunner) *Runner{
	utils.FailOnNil(rr,"Rest Runner")
	r.restRunner = rr
	return r
}

func (r *Runner) WithWebSocketsServer() *Runner{
	utils.FailFast("Not implemented")
	return r
}

func (r *Runner) WithGrpc() *Runner{
	utils.FailFast("Not implemented")
	return r
}


func (r *Runner) MustRun(ctx context.Context, res depresolver.Dependencies) {
	var wg sync.WaitGroup

	if r.restRunner != nil {
		r.runOne(ctx, &wg, r.restRunner)
	}



	wg.Wait()
}



func (r *Runner) runOne(ctx context.Context, wg *sync.WaitGroup, s singleRunner) {
	wg.Add(1)
	go func() {
		if err := s.run(ctx); err != nil {
			zerologger.Panic().Err(err).Str("runner", s.String()).Msg("cannot run")
		}
		wg.Done()
	}()
}