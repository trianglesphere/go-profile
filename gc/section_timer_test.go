package gc

import (
	"fmt"
	"runtime/debug"
	"testing"
	"time"
)

func newTestTimer(unixStart int64) *GCAwareSectionTimer {
	t := time.Unix(unixStart, 0)
	return &GCAwareSectionTimer{
		Name:  "",
		start: t,
		prev:  t,
	}
}

func (g *GCAwareSectionTimer) testMark(label string, unixMark int64) {
	t := time.Unix(unixMark, 0)
	mark := mark{name: label, start: g.prev, end: t}
	g.prev = t
	g.marks = append(g.marks, mark)
}

func (g *GCAwareSectionTimer) testEndMark(label string, unixMark int64) {
	g.testMark(label, unixMark)
	g.end = time.Unix(unixMark, 0)
}

// runSequence runs a test squence with the simulated marks and gc info.
// marks and pauseEnds are in unix seconds and then converted to times
// The first mark is the start and the last mark is the end.
func runSequence(marks, pauseEnds []int64, pauseTimes []time.Duration) *GCAwareSectionTimer {
	t := newTestTimer(marks[0])
	for i, m := range marks[1 : len(marks)-1] {
		t.testMark(fmt.Sprint(i), m)
	}
	t.testEndMark("end", marks[len(marks)-1])
	var ends []time.Time
	for _, end := range pauseEnds {
		ends = append(ends, time.Unix(end, 0))
	}
	// Reverse arrays to match GCStats formats
	for i, j := 0, len(ends)-1; i < j; i, j = i+1, j-1 {
		ends[i], ends[j] = ends[j], ends[i]
	}
	for i, j := 0, len(pauseTimes)-1; i < j; i, j = i+1, j-1 {
		pauseTimes[i], pauseTimes[j] = pauseTimes[j], pauseTimes[i]
	}
	gcs := debug.GCStats{
		Pause:    pauseTimes,
		PauseEnd: ends,
	}
	t.calculateGCTimes(gcs)
	return t
}

func TestGcAwareTimer(t *testing.T) {
	testCases := []struct {
		marks           []int64
		pause           []time.Duration
		pauseEnd        []int64
		expectedGcCount []int
		expectedGcTimes []time.Duration
	}{
		{
			marks:           []int64{10, 20},
			pause:           []time.Duration{100 * time.Nanosecond},
			pauseEnd:        []int64{15},
			expectedGcCount: []int{1},
			expectedGcTimes: []time.Duration{100 * time.Nanosecond},
		},
		{
			marks:           []int64{10, 12, 20},
			pause:           []time.Duration{100 * time.Nanosecond},
			pauseEnd:        []int64{11},
			expectedGcCount: []int{1, 0},
			expectedGcTimes: []time.Duration{100 * time.Nanosecond, 0 * time.Nanosecond},
		},
		{
			marks:           []int64{10, 12, 20},
			pause:           []time.Duration{100 * time.Nanosecond},
			pauseEnd:        []int64{15},
			expectedGcCount: []int{0, 1},
			expectedGcTimes: []time.Duration{0 * time.Nanosecond, 100 * time.Nanosecond},
		},
		{
			marks:           []int64{10, 12, 20},
			pause:           []time.Duration{100 * time.Nanosecond, 200 * time.Nanosecond},
			pauseEnd:        []int64{15, 16},
			expectedGcCount: []int{0, 2},
			expectedGcTimes: []time.Duration{0 * time.Nanosecond, 300 * time.Nanosecond},
		},
		{
			marks:           []int64{10, 12, 20},
			pause:           []time.Duration{100 * time.Nanosecond, 100 * time.Nanosecond, 200 * time.Nanosecond},
			pauseEnd:        []int64{9, 15, 16},
			expectedGcCount: []int{0, 2},
			expectedGcTimes: []time.Duration{0 * time.Nanosecond, 300 * time.Nanosecond},
		},
		{
			marks:           []int64{10, 12, 20},
			pause:           []time.Duration{100 * time.Nanosecond, 100 * time.Nanosecond, 200 * time.Nanosecond, 50 * time.Nanosecond},
			pauseEnd:        []int64{9, 15, 16, 22},
			expectedGcCount: []int{0, 2},
			expectedGcTimes: []time.Duration{0 * time.Nanosecond, 300 * time.Nanosecond},
		},
	}
	for i, testCase := range testCases {
		timer := runSequence(testCase.marks, testCase.pauseEnd, testCase.pause)
		timer.verifyComputedFields(t, i)
		timer.verifyExpected(t, i, testCase.expectedGcCount, testCase.expectedGcTimes)
		fmt.Println(timer)
	}
}

func (g *GCAwareSectionTimer) verifyComputedFields(t *testing.T, caseNumber int) {
	var gcCount int
	var runTime, gcTime, elapsed time.Duration
	for _, section := range g.Sections {
		// Sum each section
		gcCount += section.GcCount
		runTime += section.RunTime
		gcTime += section.GcTime
		elapsed += section.Elapsed
		// Sanity check each section
		if section.Elapsed != section.RunTime+section.GcTime {
			t.Errorf("case %v. Elapsed != runtime + gc time: %v = %v + %v", caseNumber, elapsed, section.RunTime, section.GcTime)
		}
		if section.Elapsed < 0 || section.RunTime < 0 || section.GcTime < 0 {
			t.Errorf("case %v. Negative duration. Elapsed: %v; runtime: %v; gc time: %v;", caseNumber, section.Elapsed, section.RunTime, section.GcTime)
		}
	}
	// Verify sums
	if g.Elapsed != elapsed {
		t.Errorf("case %v. Elapsed is not the same as the sum of the sections. Expected %v, got %v", caseNumber, elapsed, g.Elapsed)
	}
	if g.GcTime != gcTime {
		t.Errorf("case %v. GcTime is not the same as the sum of the sections. Expected %v, got %v", caseNumber, gcTime, g.GcTime)
	}
	if g.RunTime != runTime {
		t.Errorf("case %v. RunTime is not the same as the sum of the sections. Expected %v, got %v", caseNumber, runTime, g.RunTime)
	}
	if g.GcCount != gcCount {
		t.Errorf("case %v. GcCount is not the same as the sum of the sections. Expected %v, got %v", caseNumber, gcCount, g.GcCount)
	}
}

func (g *GCAwareSectionTimer) verifyExpected(t *testing.T, caseNumber int, gcCounts []int, gcTimes []time.Duration) {
	for i, section := range g.Sections {
		if section.GcCount != gcCounts[i] {
			t.Errorf("case %v. Mismatch in %v-th section gc count. Expected %v, Got %v", caseNumber, i, gcCounts[i], section.GcCount)
		}
		if section.GcTime != gcTimes[i] {
			t.Errorf("case %v. Mismatch in %v-th section gc time. Expected %v, Got %v", caseNumber, i, gcTimes[i], section.GcTime)
		}
	}
}
