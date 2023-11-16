package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Item struct {
	value    string
	priority time.Time
	index    int
}

type PriorityQueue []*Item

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Swap swaps two items in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority.Before(pq[j].priority)
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the highest priority from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// This is the most basic example: it prints events to the terminal as we
// receive them.
func watch(paths ...string) {
	if len(paths) < 1 {
		exit("must specify at least one path to watch")
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		exit("creating a new watcher: %s", err)
	}
	defer w.Close()

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("tick")
				evalThePurge(&pq)
			}
		}
	}()

	// Start listening for events.
	go watchLoop(w, &pq)

	// Add all paths from the commandline.
	for _, p := range paths {
		err = w.Add(p)
		if err != nil {
			exit("%q: %s", p, err)
		}
	}

	printTime("ready; press ^C to exit")
	<-make(chan struct{}) // Block forever
}

func evalThePurge(pq *PriorityQueue) {
	fmt.Printf("the purge %v\n", len(*pq))
	for len(*pq) >= 1 && (*pq)[0].priority.Before(time.Now().Add(-10*time.Second)) {
		item := heap.Pop(pq).(*Item)
		fmt.Printf(item.priority.Format("15:04:05.0000") + " " + item.value + "\n")
		fileInfo, err := os.Stat(item.value)
		// Checks for the error
		if err != nil {
			// ignore not exists errors.  Multiple writes can put multiple entries on the heap
			if !os.IsNotExist(err) {
				log.Fatal(err)
			} else {
				continue
			}
		}

		// Gives the modification time
		modificationTime := fileInfo.ModTime()
		if modificationTime.After(item.priority) {
			heap.Push(pq, &Item{value: item.value, priority: time.Now()})
			fmt.Println("file modified.  Putting it back. %v", item.value)
		} else {
			fmt.Println("removing %v", item.value)
			e := os.Remove(item.value)
			if e != nil {
				log.Fatal(e)
			}
		}
	}
}

func watchLoop(w *fsnotify.Watcher, pq *PriorityQueue) {
	i := 0

	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			printTime("ERROR: %s", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			// Just print the event nicely aligned, and keep track how many
			// events we've seen.
			i++
			printTime("%3d %s", i, e)
			heap.Push(pq, &Item{value: e.Name, priority: time.Now()})
		}
		evalThePurge(pq)
	}
}
