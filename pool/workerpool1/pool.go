package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNoIdleWorkerInPool = errors.New("no idle worker in pool") // workerpool中任务已满，没有空闲goroutine用于处理新任务
	ErrWorkerPoolFreed    = errors.New("workerpool freed")       // workerpool已终止运行
)

type Task func()

type Pool struct {
	capacity int           // work大小
	active   chan struct{} // work的计数channel
	tasks    chan Task     // 用户任务channel

	wg   sync.WaitGroup //
	quit chan struct{}  // exit的信号量
}

const (
	defaultCapacity = 100
	maxCapacity     = 10000
)

/*可以理解为一个构造函数*/
func New(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		active:   make(chan struct{}, capacity),
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
	}

	fmt.Printf("workerPool start\n")

	go p.run()

	return p
}

/*主要提供线程池的初始化操作，以及检测退出信号量，主要通过select实现*/
func (p *Pool) run() {
	idx := 0

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			idx++
			p.newWorker(idx)
		}
	}
}

/*主要是实现一个对tasks检测的routine，同时也需要对退出信号量进行检测*/
func (p *Pool) newWorker(i int) {
	p.wg.Add(1)

	go func() {
		defer func() {
			// 对panic进行处理
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", i)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", i)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", i)
				t()
			}
		}
	}()
}

func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	}
}

/*释放线程池，等待当前work执行完*/
func (p *Pool) Free() {
	close(p.quit) // make sure all worker and p.run exit and schedule return error
	p.wg.Wait()
	fmt.Printf("workerpool freed\n")
}
