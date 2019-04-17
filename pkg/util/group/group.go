package group

import "sync"

type work func() error

type Group struct {
	wg    sync.WaitGroup
	works []work
}

func NewGroup(size int) *Group {
	return &Group{
		works: make([]work, 0, size),
	}
}

func (g *Group) Add(fn func() error) {
	g.works = append(g.works, fn)
}

func (g *Group) Run() []error {
	defer g.Reset()
	errs := make([]error, 0, len(g.works))
	for _, v := range g.works {
		w := v
		g.wg.Add(1)
		go func() {
			defer g.wg.Done()
			err := w()
			errs = append(errs, err)
		}()
	}

	g.wg.Wait()

	return errs
}

func (g *Group) Reset() {
	g.works = g.works[0:0]
}
