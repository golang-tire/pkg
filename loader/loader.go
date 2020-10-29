package loader

import (
	"context"
	"sort"
	"sync"
)

var (
	loads  items
	lock   sync.RWMutex
	lastId int
)

func (l items) Len() int {
	return len(l)
}

func (l items) Less(i, j int) bool {
	if loads[i].priority == 0 && loads[j].priority == 0 {
		return loads[i].id < loads[j].id
	}
	return loads[i].priority < loads[j].priority
}

func (l items) Swap(i, j int) {
	loads[i], loads[j] = loads[j], loads[i]
}

type items []loadItem

type loadItem struct {
	id       int
	function func(ctx context.Context, args interface{}) (interface{}, error)
	priority int
}

// Add new loader function to loader stack, functions will sort and run base on priority or
// their add order if priority was zero
func Add(function func(ctx context.Context, args interface{}) (interface{}, error), priority int) {
	lock.Lock()
	defer lock.Unlock()

	loads = append(loads, loadItem{
		id:       lastId,
		function: function,
		priority: priority,
	})

	lastId++
}

// Flush remove all loader functions
func Flush() {
	loads = nil
}

// Run loader with given args that will pass to each loader function
// args can be modify by each loader function and modified version will pass to next one
func Run(ctx context.Context, args interface{}) error {
	// sort items
	sort.Sort(loads)

	var err error
	for _, item := range loads {
		args, err = item.function(ctx, args)
		if err != nil {
			return err
		}
	}
	return nil
}
