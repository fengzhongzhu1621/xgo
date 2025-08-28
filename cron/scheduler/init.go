package scheduler

var sched *Scheduler

func Init(s *Scheduler) {
	sched = NewScheduler()
}

func GetScheduler() *Scheduler {
	return sched
}
