package emotions

import (
	"sort"
	"sync"
)

type aggregator struct {
	mu      sync.Mutex
	samples []EmotionSample
}

func (a *aggregator) add(s EmotionSample) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.samples = append(a.samples, s)
}

func (a *aggregator) closeWindow() (WindowSummary, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.samples) == 0 {
		return WindowSummary{}, false
	}

	dist := map[string]int{}
	begin, end := a.samples[0].TimestampMs, a.samples[0].TimestampMs
	for _, s := range a.samples {
		dist[s.Emotion]++
		if s.TimestampMs < begin {
			begin = s.TimestampMs
		}
		if s.TimestampMs > end {
			end = s.TimestampMs
		}
	}

	summary := WindowSummary{
		BeginMs:      begin,
		EndMs:        end,
		Samples:      len(a.samples),
		Distribution: dist,
		Predominant:  predominantEmotion(dist),
	}

	a.samples = nil
	return summary, true
}

func predominantEmotion(dist map[string]int) string {
	keys := make([]string, 0, len(dist))
	for k := range dist {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	predominant := ""
	maxCount := -1
	for _, k := range keys {
		if dist[k] > maxCount {
			maxCount = dist[k]
			predominant = k
		}
	}
	return predominant
}
