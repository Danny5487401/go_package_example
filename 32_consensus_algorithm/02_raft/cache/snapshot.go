package cache

import (
	"github.com/hashicorp/raft"
)

type Snapshot struct {
	Cm *cacheManager
}

// Persist saves the FSM Snapshot out to the given sink.
// 在Persist里面，自己把缓存里面的数据用json格式化的方式来生成快照，sink.Write就是把快照写入snapStore，我们刚才定义的是FileSnapshotStore，所以会把数据写入文件。
func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	snapshotBytes, err := s.Cm.Marshal()
	if err != nil {
		sink.Cancel()
		return err
	}

	if _, err := sink.Write(snapshotBytes); err != nil {
		sink.Cancel()
		return err
	}

	if err := sink.Close(); err != nil {
		sink.Cancel()
		return err
	}
	return nil
}

func (f *Snapshot) Release() {}
