package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		out := make(Bi)
		close(out)
		return out
	}
	if len(stages) == 0 {
		return in
	}

	ch := in

	if done == nil {
		for _, s := range stages {
			ch = s(ch)
		}
		return ch
	}

	for _, s := range stages {
		out := make(Bi)

		go func(in In) {
			defer close(out)

			for {
				select {
				case <-done:
					// for range in {
					// 	// draining 'in' channel to allow existing goroutines to finish
					// }
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- v
				}
			}
		}(ch)

		ch = s(out)
	}

	return ch
}
