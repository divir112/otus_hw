package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := gen(in, done)

	for _, stage := range stages {
		in = stage(ch)
		ch = gen(in, done)
	}

	return ch
}

func gen(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case e, ok := <-in:
				if !ok {
					return
				}
				out <- e

			case <-done:
				go func() {
					for {
						_, ok := <-in
						if !ok {
							return
						}
					}
				}()
				return
			}
		}
	}()
	return out
}
