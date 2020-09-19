package zoo

import "sync"

type Living interface {
	Die() // sync call to die
	ID() string
}

type Zoo struct {
	sync.RWMutex
	livings map[string]Living
	Dead    *SafeMap
}

// New `Zoo` to manage `Animal`s
func NewZoo(xs ...Living) *Zoo {
	rv := &Zoo{livings: make(map[string]Living), Dead: NewSafeMap()}
	for _, x := range xs {
		rv.AddAnimal(x)
	}
	return rv
}

// Add `Animal` to `Zoo`
func (z *Zoo) AddAnimal(x Living) *Zoo {
	z.Lock()
	defer z.Unlock()

	z.livings[x.ID()] = x
	z.Dead.Set(x.ID(), false)
	return z
}

// Destroy the zoo
func (z *Zoo) Destroy() {
	z.Lock()
	defer z.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(z.livings))

	for _, x := range z.livings {
		go func(x Living) {
			defer func() {
				z.Dead.Set(x.ID(), true)
				wg.Done()
			}()
			x.Die()
		}(x)
	}

	wg.Wait()
}

func (z *Zoo) IsDead(x Living) bool {
	if v := z.Dead.Get(x.ID()); v != nil {
		return v.(bool)
	}
	return true
}
