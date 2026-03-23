package main

import (
	"fmt"
	"sync"
)

func main() {
	// start := time.Now()
	// defer func() {
	// 	fmt.Println(time.Since(start))
	// }()

	// smokeSignal := make(chan bool)
	// evilNinjas := []string{"Tommy", "Johnny", "Bobby", "Andy"}

	// for _, evilNinja := range evilNinjas {
	// 	go attack(evilNinja, smokeSignal)
	// 	fmt.Println(<-smokeSignal)
	// }

	var beeper sync.WaitGroup
	evilNinjas := []string{"Tommy", "Johnny", "Bobby", "Andy"}
	beeper.Add(len(evilNinjas))

	for _, evilNinja := range evilNinjas {
		go attack(evilNinja, &beeper)
	}
	beeper.Wait()
	fmt.Println("Mission completed")
}

// func attack(target string, attacked chan bool) {
// 	// time.Sleep(time.Second)
// 	// fmt.Println("Throwing ninja stars at", target)
// 	// attacked <- true
// }

func attack(evilNinja string, beeper *sync.WaitGroup) {
	fmt.Println("Attacked evil ninja:", evilNinja)
	beeper.Done()
}
