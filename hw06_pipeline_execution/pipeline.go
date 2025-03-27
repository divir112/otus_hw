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
		ch = stage(ch)
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
				// select {
				// case <-done:
				// 	return
				// case out <- e:
				// default:
				// 	if ok {
				// 		out <- e
				// 	}
				// }

			case <-done:
				return
			}
		}

		// for {
		// 	select {
		// 	case out <- <-in:
		// 	case <-done:
		// 		return
		// 	}
		// }

		// for elem := range in {
		// 	select {
		// 	case <-done:
		// 		return
		// 	default:
		// 		out <- elem
		// 	}
		// }
	}()
	return out
}
