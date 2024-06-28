package main

import (
	"fmt"
	"sync"
	"time"
)

type cowboy struct {
	health int
	damage int
}

type mapOfCowboys struct {
	cowboys map[string]cowboy
	m       sync.RWMutex
}

func main() {

	mapCowboys := mapOfCowboys{
		cowboys: make(map[string]cowboy),
	}

	initCowboys(mapCowboys.cowboys)

	syncCh := make(chan int)
	var wg sync.WaitGroup

	fmt.Println("Fight!!!")

	for name := range mapCowboys.cowboys {
		wg.Add(1)
		go mapCowboys.startShoot(name, syncCh, &wg)
	}

	close(syncCh)
	wg.Wait()

	for winner := range mapCowboys.cowboys {
		fmt.Println(winner, "wins!")
	}

}

func (cowboys *mapOfCowboys) startShoot(name string, syncCh chan int, wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		_, ok := <-syncCh
		if !ok {
			break
		}
	}

	for {
		if !cowboys.shooting(name) {
			return
		}
		time.Sleep(time.Second)
	}

}

func (cowboys *mapOfCowboys) shooting(name string) bool {
	cowboys.m.Lock()
	defer cowboys.m.Unlock()

	if len(cowboys.cowboys) == 1 {
		return false
	}

	if cowboys.deleteMyself(name) {
		return false
	}

	nameToShoot := cowboys.getRandomCowboyName(name)

	fmt.Println(name, "health:", cowboys.cowboys[name].health, "hit", nameToShoot, "on", cowboys.cowboys[name].damage, "with health", cowboys.cowboys[nameToShoot].health)

	cowboys.shootTheName(name, nameToShoot)
	return true
}

func (cowboys *mapOfCowboys) deleteMyself(name string) bool {
	if cowboys.cowboys[name].health <= 0 {
		delete(cowboys.cowboys, name)
		return true
	}
	return false
}

func (cowboys *mapOfCowboys) shootTheName(nameOfShooter, nameToShoot string) {
	cb := cowboys.cowboys[nameToShoot]
	cb.health -= cowboys.cowboys[nameOfShooter].damage

	cowboys.cowboys[nameToShoot] = cb
	if cb.health <= 0 {
		delete(cowboys.cowboys, nameToShoot)

	}
}

func (cowboys *mapOfCowboys) getRandomCowboyName(nameOfShooter string) string {
	var cowboyName string
	for name := range cowboys.cowboys {
		cowboyName = name
		if nameOfShooter != cowboyName {
			break
		}
	}

	return cowboyName
}

func initCowboys(cowboys map[string]cowboy) {

	cowboys["John"] = cowboy{
		health: 2,
		damage: 1,
	}

	cowboys["Bill"] = cowboy{
		health: 8,
		damage: 2,
	}

	cowboys["Sam"] = cowboy{
		health: 10,
		damage: 1,
	}

	cowboys["Peter"] = cowboy{
		health: 5,
		damage: 3,
	}

	cowboys["Philip"] = cowboy{
		health: 50,
		damage: 1,
	}
}
