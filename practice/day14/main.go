package main

import (
	"fmt"
	"sync"
)

/*
q1.  Make a struct Theater with the following fields: Seats(int=1), RWMutex, userName chan string.

     Create two methods over a struct

     The first method book a seat in the theater. If the seat is equal to zero, no one can book it.
     ( In the booking method, put simple print statements that show booking has been made if seats are available)

     Once the seat is booked in Theater, add the name of the user in the userName channel.
     Create a second Method, printInvoice(),  It fetches the userName from the channel and print it.

    Note:-
     Call the bookSeats & printInvoice() method as a goroutine in the main function.
     For example:-

     for i:=1; i<=3; i++ {
          go t.bookSeats()
     }
     go t.printInvoice()

     The program should quit gracefully without deadlock.

*/

type Theater struct {
	seats    int
	mutex    *sync.RWMutex
	userName chan string
}

func (t *Theater) bookSeats(userName string, wg *sync.WaitGroup) {
	defer wg.Done()
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.seats == 0 {
		fmt.Println("Seat not available for:", userName)
		return
	}
	// don't do locking here
	// t.mutex.Lock()
	// defer t.mutex.Unlock()
	t.seats--
	t.userName <- userName
	fmt.Println("bookSeats:", userName)
}

func (t *Theater) printInvoice(wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range t.userName {
		fmt.Println("printInvoice:", v)
	}
}

func main() {
	t := Theater{seats: 2, mutex: new(sync.RWMutex), userName: make(chan string)}
	users := []string{"Alka", "Komal", "Sandra", "Diwakar", "Kanchan", "Gobiya"}

	// equivalent to using new
	wg := &sync.WaitGroup{}
	workerWG := &sync.WaitGroup{}

	for _, v := range users {
		workerWG.Add(1)
		go t.bookSeats(v, workerWG)
	}

	wg.Add(1)
	go t.printInvoice(wg)

	workerWG.Wait()
	close(t.userName)

	wg.Wait()
}
