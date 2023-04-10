package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func readInDone(in In, done In) Out {
	tmp := make(Bi)
	go func() {
		defer close(tmp)
		for {
			select {
			case <-done:
				return
			case i, ok := <-in:
				if !ok {
					return
				}
				tmp <- i
			}
		}
	}()
	return tmp
}

func runPipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, s := range stages {
		out = s(readInDone(out, done))
	}
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		out := make(Bi)
		close(out)
		return out
	}

	return runPipeline(in, done, stages...)
}
