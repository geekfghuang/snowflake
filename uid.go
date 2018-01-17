/*
 * Copyright 2018 geekfghuang. All Rights Reserved.
 *
 * From the Twitter-Snowflake algorithm, the 64 bit self
 * increasing ID generation theory
 *
 * | timestamp(ms)42 | worker id(10) | sequence(12) |
 *******************************************************/
package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	Epoch         = 1516170660000 // from now on
	WorkerIdBits  = 10
	SequenceBits  = 12

	WorkerIdShift  = SequenceBits
	TimeStampShift = SequenceBits + WorkerIdBits

	SequenceMask = 0xfff // equal to getSequenceMask()
	MaxWorker    = 0x3ff // equal to getMaxWorkerId()
)

type Worker struct {
	workerId      int64
	lastTimeStamp int64
	sequence      int64
	maxWorkerId   int64
	lock          *sync.Mutex
}

func getMaxWorkerId() int64 {
	return -1 ^ -1 << WorkerIdBits
}

func getSequenceMask() int64 {
	return -1 ^ -1 << SequenceBits
}

func NewWorker(workerId int64) (worker *Worker, err error) {
	worker = new(Worker)
	worker.maxWorkerId = getMaxWorkerId()
	if workerId > worker.maxWorkerId || workerId < 0 {
		return nil, errors.New("worker not fit")
	}
	worker.workerId = workerId
	worker.lastTimeStamp = -1
	worker.sequence = 0
	worker.lock = new(sync.Mutex)
	return
}

// return in ms
func (worker *Worker) timeGen() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func (worker *Worker) timeReGen(last int64) int64 {
	ts := worker.timeGen()
	for ts <= last {
		ts = worker.timeGen()
	}
	return ts
}

func (worker *Worker) NextId() (id int64, err error) {
	worker.lock.Lock()
	defer worker.lock.Unlock()

	ts := worker.timeGen()
	if ts == worker.lastTimeStamp {
		worker.sequence = (worker.sequence + 1) & SequenceMask
		if worker.sequence == 0 {
			ts = worker.timeReGen(ts)
		}
	} else {
		worker.sequence = 0
	}

	if ts < worker.lastTimeStamp {
		return -1, errors.New("clock moved backwards, refuse gen id")
	}
	worker.lastTimeStamp = ts
	id = (ts - Epoch) << TimeStampShift | worker.workerId << WorkerIdShift | worker.sequence
	return
}

// reverse id to Time timestamp, workid, seq
func ParseId(id int64) (t time.Time, ts int64, workerId int64, seq int64) {
	seq = id & SequenceMask
	workerId = (id >> WorkerIdShift) & MaxWorker
	ts = (id >> TimeStampShift) + Epoch
	t = time.Unix(ts / 1000, (ts % 1000) * 1000000)
	return
}