package cron

import (
	"fmt"
	"runtime"
	"sort"
	"time"
)

//
type Cron struct {
	entries []*Entry
	add     chan *Entry
	stop    chan struct{}
	running bool
	output  Outputer
}

// entry信息 输出接口
type Outputer interface {
	Output(*Entry) error
}

type DefaultOutput struct{}

func (w *DefaultOutput) Output(e *Entry) error {
	log.Info(e.Next.String(), e.Name)
	return nil
}

// Job is an interface for submitted cron jobs.
type Job interface {
	Run()
}

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
	Name string

	// shedule the job
	Schedule *Schedule

	// the next time the job will run
	Next time.Time

	// the last time this job was run
	Prev time.Time

	// the job to run
	Job Job
}

// byTime is a wrapper for sorting the entry array by time
// (with zero time at the end).
type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}

func New() *Cron {
	return NewWithLocation(time.Now().Location())
}

// NewWithLocation returns a new Cron job runner.
func NewWithLocation(location *time.Location) *Cron {
	return &Cron{
		entries: nil,
		add:     make(chan *Entry),
		stop:    make(chan struct{}),
		running: false,
		output:  &DefaultOutput{},
	}
}

// Set the custom output. 设置自定义输出
func (c *Cron) SetOutput(output Outputer) {
	c.output = output
}

func (c *Cron) Output(e *Entry) {
	c.output.Output(e)
}

// A wrapper that turns a func() into a cron.Job
type FuncJob func()

func (f FuncJob) Run() { f() }

// AddFunc adds a func to the Cron to be run on the schedule
func (c *Cron) AddFunc(name, cronExpr string, cmd func()) error {
	return c.AddJob(name, cronExpr, FuncJob(cmd))
}

func (c *Cron) AddJob(name, cronExpr string, cmd Job) error {
	schedule, err := Parse(cronExpr)
	if err != nil {
		return err
	}
	c.Schedule(name, schedule, cmd)
	return nil
}

// Schedule adds a Job to the Cron to be run on the given shedule.
func (c *Cron) Schedule(name string, schedule *Schedule, cmd Job) {
	entry := &Entry{
		Name:     name,
		Schedule: schedule,
		Job:      cmd,
	}
	if !c.running {
		c.entries = append(c.entries, entry)
		return
	}
	c.add <- entry
}

func (c *Cron) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

func (c *Cron) runWithRecovery(j Job) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err := fmt.Sprintf("cron: panic running job: %v\n%s", r, buf)
			log.Error(err)
		}
	}()
	j.Run()
}

func (c *Cron) run() {
	now := time.Now()
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		sort.Sort(byTime(c.entries))

		var timer *time.Timer
		if len(c.entries) == 0 || c.entries[0].Next.IsZero() {
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(c.entries[0].Next.Sub(now))
		}

		select {
		case now = <-timer.C:
			for _, e := range c.entries {
				if e.Next.After(now) || e.Next.IsZero() {
					break
				}
				c.Output(e)
				go c.runWithRecovery(e.Job)
				e.Prev = e.Next
				e.Next = e.Schedule.Next(now)
			}
		case newEntry := <-c.add:
			timer.Stop()
			now = time.Now()
			newEntry.Next = newEntry.Schedule.Next(now)
			c.entries = append(c.entries, newEntry)

		case <-c.stop:
			timer.Stop()
			return
		}
	}
}
