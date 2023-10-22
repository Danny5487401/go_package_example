<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [confluent-kafka源码分析](#confluent-kafka%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [C库参考链接](#c%E5%BA%93%E5%8F%82%E8%80%83%E9%93%BE%E6%8E%A5)
  - [生产者](#%E7%94%9F%E4%BA%A7%E8%80%85)
  - [消费者](#%E6%B6%88%E8%B4%B9%E8%80%85)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# confluent-kafka源码分析
底层使用c库

## C库参考链接
配置链接： https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
源码链接： https://github.com/edenhill/librdkafka

## 生产者
结构体
```go
type Producer struct {
	events         chan Event
	produceChannel chan *Message   //需要发送的消息
	handle         handle

	// Terminates the poller() goroutine
	// 是否关闭poller
	pollerTermChan chan bool
}
```
生产者和消费者的公共结构体
```go

type handle struct {
	rk  *C.rd_kafka_t
	rkq *C.rd_kafka_queue_t

	// Forward logs from librdkafka log queue to logs channel.
	logs          chan LogEvent
	logq          *C.rd_kafka_queue_t
	closeLogsChan bool

	// Topic <-> rkt caches
	rktCacheLock sync.Mutex
	// topic name -> rkt cache
	rktCache map[string]*C.rd_kafka_topic_t
	// rkt -> topic name cache
	rktNameCache map[*C.rd_kafka_topic_t]string

	// Cached instance name to avoid CGo call in String()
	name string

	//
	// cgo map
	// Maps C callbacks based on cgoid back to its Go object
	cgoLock   sync.Mutex
	cgoidNext uintptr
	cgomap    map[int]cgoif

	//
	// producer
	//
	p *Producer

	// Forward delivery reports on Producer.Events channel
	fwdDr bool

	// Enabled message fields for delivery reports and consumed messages.
	msgFields *messageFields

	//
	// consumer
	//
	c *Consumer

	// WaitGroup to wait for spawned go-routines to finish.
	waitGroup sync.WaitGroup
}
```

初始化流程
```go
func NewProducer(conf *ConfigMap) (*Producer, error) {
    
    // 1. librdkafka版本确认
	err := versionCheck()
	if err != nil {
		return nil, err
	}

    // 2. 初始化空结构体
	p := &Producer{}

	// before we do anything with the configuration, create a copy such that
	// the original is not mutated.
	confCopy := conf.clone()

	// 3. 配置获取转换
    // ....
	// Convert ConfigMap to librdkafka conf_t
	cConf, err := confCopy.convert()
	if err != nil {
		return nil, err
	}

	cErrstr := (*C.char)(C.malloc(C.size_t(256)))
	defer C.free(unsafe.Pointer(cErrstr))

    // 4. 注册生产者关心的一系列事件
	C.rd_kafka_conf_set_events(cConf, C.RD_KAFKA_EVENT_DR|C.RD_KAFKA_EVENT_STATS|C.RD_KAFKA_EVENT_ERROR|C.RD_KAFKA_EVENT_OAUTHBEARER_TOKEN_REFRESH)

    // 5. cgo初始化创建生产者实例
	// Create librdkafka producer instance
	p.handle.rk = C.rd_kafka_new(C.RD_KAFKA_PRODUCER, cConf, cErrstr, 256)
	if p.handle.rk == nil {
		return nil, newErrorFromCString(C.RD_KAFKA_RESP_ERR__INVALID_ARG, cErrstr)
	}

	p.handle.p = p
	p.handle.setup()
	// 6. 获取生产队列的主replication
	p.handle.rkq = C.rd_kafka_queue_get_main(p.handle.rk)
	p.events = make(chan Event, eventsChanSize)
	p.produceChannel = make(chan *Message, produceChannelSize)
	p.pollerTermChan = make(chan bool)

	if logsChanEnable {
		p.handle.setupLogQueue(logsChan, p.pollerTermChan)
	}

	p.handle.waitGroup.Add(1)
	// 7. 起协程监听生产者事件
	go func() {
		poller(p, p.pollerTermChan)
		p.handle.waitGroup.Done()
	}()

	// non-batch or batch producer, only one must be used
	var producer func(*Producer)
	// 根据配置go.batch.producer选择不同生产者
	if batchProducer {
		producer = channelBatchProducer
	} else {
		producer = channelProducer
	}

	p.handle.waitGroup.Add(1)
	// 8.起协程生产消息
	go func() {
		producer(p)
		p.handle.waitGroup.Done()
	}()

	return p, nil
}
```

生产者事件
```go
// poller polls the rd_kafka_t handle for events until signalled for termination
func poller(p *Producer, termChan chan bool) {
	for {
		select {
		case _ = <-termChan:
			return

		default:
			_, term := p.handle.eventPoll(p.events, 100, 1000, termChan)
			if term {
				return
			}
			break
		}
	}
}

```
事件分发函数
```go
//  kafka/event.go

// eventPoll polls an event from the handler's C rd_kafka_queue_t,
// translates it into an Event type and then sends on `channel` if non-nil, else returns the Event.
// term_chan is an optional channel to monitor along with producing to channel
// to indicate that `channel` is being terminated.
// returns (event Event, terminate Bool) tuple, where Terminate indicates
// if termChan received a termination event.
func (h *handle) eventPoll(channel chan Event, timeoutMs int, maxEvents int, termChan chan bool) (Event, bool) {

	var prevRkev *C.rd_kafka_event_t
	term := false

	var retval Event

	if channel == nil {
		maxEvents = 1
	}
out:
	for evcnt := 0; evcnt < maxEvents; evcnt++ {
		var evtype C.rd_kafka_event_type_t
		var gMsg C.glue_msg_t
		gMsg.want_hdrs = C.int8_t(bool2cint(h.msgFields.Headers))
		// cgo的事件进行了封装
		rkev := C._rk_queue_poll(h.rkq, C.int(timeoutMs), &evtype, &gMsg, prevRkev)
		prevRkev = rkev
		timeoutMs = 0

		retval = nil

		switch evtype {
		// 事件fetch
		case C.RD_KAFKA_EVENT_FETCH:
			// Consumer fetch event, new message.
			// Extracted into temporary gMsg for optimization
			retval = h.newMessageFromGlueMsg(&gMsg) 
        // 事件 reblance
		case C.RD_KAFKA_EVENT_REBALANCE:
			// Consumer rebalance event
			retval = h.c.handleRebalanceEvent(channel, rkev)

		case C.RD_KAFKA_EVENT_ERROR:
			// Error event
			cErr := C.rd_kafka_event_error(rkev)
			if cErr == C.RD_KAFKA_RESP_ERR__PARTITION_EOF {
				crktpar := C.rd_kafka_event_topic_partition(rkev)
				if crktpar == nil {
					break
				}

				defer C.rd_kafka_topic_partition_destroy(crktpar)
				var peof PartitionEOF
				setupTopicPartitionFromCrktpar((*TopicPartition)(&peof), crktpar)

				retval = peof

			} else if int(C.rd_kafka_event_error_is_fatal(rkev)) != 0 {
				// A fatal error has been raised.
				// Extract the actual error from the client
				// instance and return a new Error with
				// fatal set to true.
				cFatalErrstrSize := C.size_t(512)
				cFatalErrstr := (*C.char)(C.malloc(cFatalErrstrSize))
				defer C.free(unsafe.Pointer(cFatalErrstr))
				cFatalErr := C.rd_kafka_fatal_error(h.rk, cFatalErrstr, cFatalErrstrSize)
				fatalErr := newErrorFromCString(cFatalErr, cFatalErrstr)
				fatalErr.fatal = true
				retval = fatalErr

			} else {
				retval = newErrorFromCString(cErr, C.rd_kafka_event_error_string(rkev))
			}

		case C.RD_KAFKA_EVENT_STATS:
			retval = &Stats{C.GoString(C.rd_kafka_event_stats(rkev))}

		case C.RD_KAFKA_EVENT_DR:
			// Producer Delivery Report event
			// Each such event contains delivery reports for all
			// messages in the produced batch.
			// Forward delivery reports to per-message's response channel
			// or to the global Producer.Events channel, or none.
			rkmessages := make([]*C.rd_kafka_message_t, int(C.rd_kafka_event_message_count(rkev)))

			cnt := int(C.rd_kafka_event_message_array(rkev, (**C.rd_kafka_message_t)(unsafe.Pointer(&rkmessages[0])), C.size_t(len(rkmessages))))

			for _, rkmessage := range rkmessages[:cnt] {
				msg := h.newMessageFromC(rkmessage)
				var ch *chan Event

				if rkmessage._private != nil {
					// Find cgoif by id
					cg, found := h.cgoGet((int)((uintptr)(rkmessage._private)))
					if found {
						cdr := cg.(cgoDr)

						if cdr.deliveryChan != nil {
							ch = &cdr.deliveryChan
						}
						msg.Opaque = cdr.opaque
					}
				}

				if ch == nil && h.fwdDr {
					ch = &channel
				}

				if ch != nil {
					select {
					case *ch <- msg:
					case <-termChan:
						break out
					}

				} else {
					retval = msg
					break out
				}
			}

		case C.RD_KAFKA_EVENT_OFFSET_COMMIT:
			// Offsets committed
			cErr := C.rd_kafka_event_error(rkev)
			coffsets := C.rd_kafka_event_topic_partition_list(rkev)
			var offsets []TopicPartition
			if coffsets != nil {
				offsets = newTopicPartitionsFromCparts(coffsets)
			}

			if cErr != C.RD_KAFKA_RESP_ERR_NO_ERROR {
				retval = OffsetsCommitted{newErrorFromCString(cErr, C.rd_kafka_event_error_string(rkev)), offsets}
			} else {
				retval = OffsetsCommitted{nil, offsets}
			}

		case C.RD_KAFKA_EVENT_OAUTHBEARER_TOKEN_REFRESH:
			ev := OAuthBearerTokenRefresh{C.GoString(C.rd_kafka_event_config_string(rkev))}
			retval = ev

		case C.RD_KAFKA_EVENT_NONE:
			// poll timed out: no events available
			break out

		default:
			if rkev != nil {
				fmt.Fprintf(os.Stderr, "Ignored event %s\n",
					C.GoString(C.rd_kafka_event_name(rkev)))
			}

		}

		if retval != nil {
			if channel != nil {
				select {
				case channel <- retval:
				case <-termChan:
					retval = nil
					term = true
					break out
				}
			} else {
				break out
			}
		}
	}

	if prevRkev != nil {
		C.rd_kafka_event_destroy(prevRkev)
	}

	return retval, term
}
```


生产者生产消息:channelProducer为案例
```go
// channel_producer serves the ProduceChannel channel
func channelProducer(p *Producer) {
	for m := range p.produceChannel {
		err := p.produce(m, C.RD_KAFKA_MSG_F_BLOCK, nil)
		if err != nil {
			m.TopicPartition.Error = err
			p.events <- m
		}
	}
}
```
```go
func (p *Producer) produce(msg *Message, msgFlags int, deliveryChan chan Event) error {
	

	// Three problems:
	//  1) There's a difference between an empty Value or Key (length 0, proper pointer) and
	//     a null Value or Key (length 0, null pointer).
	//  2) we need to be able to send a null Value or Key, but the unsafe.Pointer(&slice[0])
	//     dereference can't be performed on a nil slice.
	//  3) cgo's pointer checking requires the unsafe.Pointer(slice..) call to be made
	//     in the call to the C function.
	//
	// Solution:
	//  Keep track of whether the Value or Key were nil (1), but let the valp and keyp pointers
	//  point to a 1-byte slice (but the length to send is still 0) so that the dereference (2)
	//  works.
	//  Then perform the unsafe.Pointer() on the valp and keyp pointers (which now either point
	//  to the original msg.Value and msg.Key or to the 1-byte slices) in the call to C (3).
	//
	
    // .一些预处理...

	
	cErr := C.do_produce(p.handle.rk, crkt,
		C.int32_t(msg.TopicPartition.Partition),
		C.int(msgFlags)|C.RD_KAFKA_MSG_F_COPY,
		valIsNull, unsafe.Pointer(&valp[0]), C.size_t(valLen),
		keyIsNull, unsafe.Pointer(&keyp[0]), C.size_t(keyLen),
		C.int64_t(timestamp),
		(*C.tmphdr_t)(unsafe.Pointer(&tmphdrs[0])), C.size_t(tmphdrsCnt),
		(C.uintptr_t)(cgoid))
	if cErr != C.RD_KAFKA_RESP_ERR_NO_ERROR {
		if cgoid != 0 {
			p.handle.cgoGet(cgoid)
		}
		return newError(cErr)
	}

	return nil
}
```
注意：produce函数和初始化的时候注册的函数底层调用的是同一个，不同的是，初始化的时候需要等待事件的到来


## 消费者

消息内容
```go
// Message represents a Kafka message
type Message struct {
	TopicPartition TopicPartition
	Value          []byte
	Key            []byte
	Timestamp      time.Time
	TimestampType  TimestampType
	Opaque         interface{}
	Headers        []Header
}
```

同理
