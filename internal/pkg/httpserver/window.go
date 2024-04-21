package httpserver

import (
	"strconv"
	"time"
)

type window struct {
	SLIDING_WINDOW_SIZE      int
	start                    int
	end                      int
	running                  bool
	hitTimeRecords           []int64
	hitTimeRecordsSerialised string //For storage in file. Format hitTime1,hitTime2,hitTime3..
	hitCountInWindow         int
}

func (w *window) appendNewHit() int64 {
	timeInt := time.Now().Unix()
	w.hitTimeRecords = append(w.hitTimeRecords, timeInt)
	w.appendNewHitSerialised(timeInt)
	return timeInt
}

func (w *window) appendNewHitSerialised(currHitTime int64) {
	if len(w.hitTimeRecords) == 1 {
		w.hitTimeRecordsSerialised = strconv.FormatInt(currHitTime, 10)
	} else {
		w.hitTimeRecordsSerialised = w.hitTimeRecordsSerialised + "," + strconv.FormatInt(currHitTime, 10)
	}
}

func (w *window) moveStart() int64 {
	w.start++
	w.hitCountInWindow--
	return w.hitTimeRecords[w.start]
}

func (w *window) moveEnd() int64 {
	w.end++
	w.hitCountInWindow++
	return w.hitTimeRecords[w.end]
}
