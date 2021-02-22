package gc

import (
	"fmt"
	"runtime/debug"
	"time"
)

func NewGcAwareTimer(name string) *GCAwareTimer {
	t := time.Now()
	return &GCAwareTimer{
		Name:  name,
		start: t,
	}
}

func (g *GCAwareTimer) End() {
	g.end = time.Now()

	// Pull out info from gc stats
	gcs := debug.GCStats{}
	debug.ReadGCStats(&gcs)
	for i, end := range gcs.PauseEnd {
		if end.After(g.end) {
			continue
		} else if end.After(g.start) {
			g.GcCount++
			g.GcTime += gcs.Pause[i]
		} else if end.Before(g.start) {
			break
		}
	}
	g.Elapsed = g.end.Sub(g.start)
	g.RunTime = g.Elapsed - g.GcTime
}

func (g GCAwareTimer) CSVHeader() string {
	return "duration,elapsed,name,gc_time,run_time,gc_count"
}

func (g GCAwareTimer) CSVString() string {
	// duration, elapsed, name, gc time, run time, gc count
	ret := fmt.Sprintf("%v,%v,%v,%v,%v,%v", g.Elapsed, g.Elapsed.Nanoseconds(), g.Name, g.GcTime.Nanoseconds(), g.RunTime.Nanoseconds(), g.GcCount)
	return ret
}
