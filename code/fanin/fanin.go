package main

import "fmt"
import "time"
import "math/rand"

func pipe(source <-chan string, sink chan<- string) {

	for {
		sink <- <-source
	}

}

func fanin(first, second <-chan string) <-chan string {

	sink := make(chan string)

	go pipe(first, sink)
	go pipe(second, sink)

	return sink

}

func boring(sink chan<- string, name string) {

	for i := 0; ; i++ {
		sink <- fmt.Sprintf("%s (%d)", name, i)
		time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
	}

}

type loop func(chan<- string, string)

func invoke(fn loop, name string) <-chan string {

	sink := make(chan string)

	go fn(sink, name)

	return sink

}

func main() {

	a := invoke(boring, "Alice")
	b := invoke(boring, "Bob")

	c := fanin(a, b)

	for i := 0; i < 10; i++ {

		fmt.Println(<-c)

	}

	fmt.Println("I'm bored now")

}
