package zoo

// A=1 B=1 go test

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

type Monster struct {
	name string
	age  int32
	exit chan bool
}

func NewMonster(name string) *Monster {
	a := &Monster{name, 0, make(chan bool)}
	go func() {
		for {
			atomic.AddInt32(&a.age, 1)
			fmt.Println(a.ID(), a.age)

			select {
			case <-a.exit:
				return
			case <-time.After(time.Second):
			}
		}
	}()
	return a
}

func (a *Monster) ID() string {
	return a.name
}
func (a *Monster) Die() {
	for {
		time.Sleep(time.Second)
		if os.Getenv(a.ID()) == "" {
			a.exit <- true
			fmt.Println(a.ID(), "die")
			return
		}
		fmt.Println(a.ID(), "would not die")
	}
}

func TestZoo(t *testing.T) {
	z := NewZoo()

	for _, x := range "ABCDEFGHIJKLMN" {
		a := NewMonster(string(x))
		z.AddAnimal(a)
	}

	time.Sleep(time.Second * 3)

	fmt.Println(z.Dead)
	z.Destroy()
	fmt.Println(z.Dead)
}
