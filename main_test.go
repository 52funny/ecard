package ecard

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

func TestRoutine(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	needJobs := []string{}
	lenOfJobs := len(needJobs)
	fmt.Println(lenOfJobs)
	time.Sleep(time.Second)
	job := make(chan string, lenOfJobs)
	result := make(chan struct{}, lenOfJobs)

	e := getEcard()
	e.MustLogin()
	for i := 0; i < 30; i++ {
		go func(id int, jobs chan string, results chan struct{}) {
			for item := range jobs {
				if ok, err := e.IsCookieOverDue(); err == nil {
					if ok {
						e.Login()
					}
				} else {
					fmt.Println(err)
					result <- struct{}{}
					continue
				}
				fee, err := e.ObtainDormitoryElectricity("0", "5", item)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(item, fee)
				results <- struct{}{}
			}
		}(i, job, result)
	}
	for i := 0; i < len(needJobs); i++ {
		job <- needJobs[i]
	}
	close(job)

	for j := 0; j < len(needJobs); j++ {
		<-result
	}
	fmt.Println(runtime.NumGoroutine())
}

func TestPool(t *testing.T) {
	jobs := []string{
		"237",
		"k237",
		"235",
		"237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235", "237",
		"k237",
		"235",
	}
	e := getEcard()
	e.MustLogin()
	defer ants.Release()
	var wait sync.WaitGroup
	p, err := ants.NewPoolWithFunc(50, func(item interface{}) {
		ans, _ := e.ObtainDormitoryElectricity("1", "5", item.(string))
		fmt.Println(item, ans)
		wait.Done()
	})
	if err != nil {
		panic(err)
	}
	for _, item := range jobs {
		wait.Add(1)
		p.Invoke(item)
	}
	wait.Wait()
}
