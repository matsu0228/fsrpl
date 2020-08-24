package main

import "time"

const fsTimeLayout = "2006-01-02T15:04:05Z"

// InterpretationEachValueForTime convert string to time.Time
func InterpretationEachValueForTime(mps map[string]interface{}) map[string]interface{} {
	for k, v := range mps {
		vs, ok := v.(string)
		if ok {
			if tm, err := time.Parse(fsTimeLayout, vs); err == nil {
				mps[k] = tm
				continue
			}
		}
		mps[k] = v
	}
	return mps
}
