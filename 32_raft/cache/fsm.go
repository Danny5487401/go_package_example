package cache

import (
	"encoding/json"
	"io"
	"log"

	"github.com/hashicorp/raft"
)

// FSM ： finite state machine，有限状态机
type FSM struct {
	Ctx *StCachedContext
	Log *log.Logger
}

type LogEntryData struct {
	Key   string
	Value string
}

// Apply applies a Raft log entry to the key-value store.
func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	e := LogEntryData{}
	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		panic("Failed unmarshaling Raft log entry. This is a bug.")
	}
	ret := f.Ctx.St.Cm.Set(e.Key, e.Value)
	f.Log.Printf("fms.Apply(), logEntry:%s, ret:%v\n", logEntry.Data, ret)
	return ret
}

// Snapshot returns a latest snapshot
// 生成一个快照结构
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &Snapshot{Cm: f.Ctx.St.Cm}, nil
}

// Restore stores the key-value store to a previous state.
// 服务重启的时候，会先读取本地的快照来恢复数据，在FSM里面定义的Restore函数会被调用
func (f *FSM) Restore(serialized io.ReadCloser) error {
	// 这里我们就简单的对数据解析json反序列化然后写入内存即可
	return f.Ctx.St.Cm.UnMarshal(serialized)
}

/*
type FSM interface {
	Apply log is invoked once a log entry is committed.
	It returns a value which will be made available in the
	ApplyFuture returned by Raft.Apply method if that
	method was called on the same Raft node as the FSM
	Apply(*Log) interface{}

	// Snapshot is used to support log compaction. This call should
	// return an FSMSnapshot which can be used to save a point-in-time
	// snapshot of the FSM. Apply and Snapshot are not called in multiple
	// threads, but Apply will be called concurrently with Persist. This means
	// the FSM should be implemented in a fashion that allows for concurrent
	// updates while a snapshot is happening.
	Snapshot() (FSMSnapshot, error)

	// Restore is used to restore an FSM from a snapshot. It is not called
	// concurrently with any other command. The FSM must discard all previous
	// state.
	Restore(io.ReadCloser) error
}

*/
