package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Result struct {
	message string
	err     error
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//candy

type iCandy interface {
	getFlavour() int
}

type candy struct {
	flavour int
}

func (c candy) getFlavour() int {
	return c.flavour
}

func randomCandy() candy {
	f := rand.Intn(4)
	return candy{flavour: f}
}

//candyEaters

type iCandyEater interface {
	eat(iCandy, chan Result)
	isFree() bool
	setFree(b bool)
	setRoot(r *candyServiceBase)
	getQueue() []iCandy
}

func (e *candyEater) getQueue() []iCandy {
	return e.root.queue
}

type candyEater struct {
	free bool
	root *candyServiceBase
}

func (e candyEater) isFree() bool {
	return e.free
}

func (e *candyEater) setFree(b bool) {
	e.free = b
}

func (e *candyEater) setRoot(r *candyServiceBase) {
	e.root = r
}

func (e *candyEater) eat(c iCandy, ch chan Result) {

	t := rand.Intn(4)
	fmt.Println("Eating a candy: random = ", t, " flavour = ", c.getFlavour())

	if t == 3 {
		r := Result{message: "I don't like this candy", err: errors.New("")}
		ch <- r
	} else {
		time.Sleep(time.Duration(t) * time.Second) // sleep imitates some time to eat a candy
		fmt.Println("Candy eaten, random time = ", t, " flavour = ", c.getFlavour())

		ch <- Result{}
	}

	// extract flavour from slice when candy is eaten
	for i, a := range e.root.flavours {
		if a == c.getFlavour() {
			e.root.flavours = append(e.root.flavours[:i], e.root.flavours[i+1:]...)
		}
	}

	e.free = true
	fmt.Println(e.root.queue) // logging queue

	e.eatNext()

}

func (e *candyEater) eatNext() {
	for i, a := range e.root.queue { // call next eat
		if !contains(e.root.flavours, a.getFlavour()) {
			err := make(chan Result)

			e.setFree(false)
			e.root.flavours = append(e.root.flavours, a.getFlavour())

			e.root.queue = append(e.root.queue[:i], e.root.queue[i+1:]...)
			go e.eat(a, err)
			res := <-err
			if res.err != nil {
				fmt.Println(res.message)
			}
			break
		}
	}
}

//candyBaseService

type candyServiceBase struct {
	candyEaters []iCandyEater
	flavours    []int
	queue       []iCandy
}

func newCandyServiceBase(e []iCandyEater) *candyServiceBase {
	r := candyServiceBase{candyEaters: e}
	for _, a := range r.candyEaters {
		a.setRoot(&r)
	}
	return &r
}

func (b *candyServiceBase) addCandy(c iCandy) {

	if contains(b.flavours, c.getFlavour()) { // put candy in queue in case the same flavour is already inside an eater
		b.queue = append(b.queue, c)
		return
	}

	errs := make(chan Result, len(b.candyEaters))

	for _, a := range b.candyEaters { // finding free candy eater

		if a.isFree() {

			go a.eat(c, errs)

			// error check was meant to be here, but reading from channel causes unexpected behaviour
			// any eat function calls view candy queue as empty

			a.setFree(false)
			b.flavours = append(b.flavours, c.getFlavour())
			return
		}

	}

	b.queue = append(b.queue, c) // if no eaters are available put candy in queue

	return

}

func main() {

	candies := []iCandy{candy{flavour: 2}, candy{flavour: 3}, candy{flavour: 2}, candy{flavour: 4}, candy{flavour: 1}, candy{flavour: 1}}

	var eaters []iCandyEater
	for i := 0; i < 2; i++ {
		eaters = append(eaters, &candyEater{free: true})
	}
	s := newCandyServiceBase(eaters)

	for _, a := range candies {
		s.addCandy(a)
	}

	var input string // make terminal wait
	fmt.Scanln(&input)
}
