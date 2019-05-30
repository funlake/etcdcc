package timer

import (
	"sync"
	"time"
)

type timeWheel struct {
	mu       sync.Mutex
	interval int
	maxpos   int
	curpos   int
	slot     map[int]chan *SlotItem
}

type SlotItem struct {
	do      func()
	timeout int
	dead    bool
	left    int
	pos     int
}

func (tw *timeWheel) GetSlotPos(timeout int) int {
	findSlot := (timeout/tw.interval + tw.curpos) % tw.maxpos
	return findSlot
}

func (tw *timeWheel) GetSlot(timeout int) (chan *SlotItem, int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	findSlot := tw.GetSlotPos(timeout)
	if tw.slot[findSlot] == nil {
		//tw.slot[findSlot] = list.New()
		tw.slot[findSlot] = make(chan *SlotItem, 5000)
	}
	return tw.slot[findSlot], findSlot
}
func (tw *timeWheel) SetInterval(timeout int, callback func()) *SlotItem {
	//find next slot
	slot, pos := tw.GetSlot(timeout)
	si := &SlotItem{do: callback, timeout: timeout, dead: false, left: 0, pos: pos}
	//slot.PushFront(si)
	go func(si *SlotItem) {
		select {
		case slot <- si:
			//logs.Info("into slot")
		default:
			//full
		}
	}(si)
	return si
}

func (tw *timeWheel) StopInterval(si *SlotItem) {
	si.dead = true
}

//无锁情况下不能直接用tw.curpos,得传进来
func (tw *timeWheel) ReWheel(si *SlotItem, curpos int) {

	findSlot := (si.timeout/tw.interval + curpos) % tw.maxpos
	//slot := tw.GetSlot(si.timeout)
	tw.mu.Lock()
	if tw.slot[findSlot] == nil {
		tw.slot[findSlot] = make(chan *SlotItem, 1000)
	}
	tw.mu.Unlock()
	si.pos = findSlot
	go func(si *SlotItem) {
		//logs.Warn(cap(tw.slot[si.pos]))
		tw.mu.Lock()
		tw.slot[si.pos] <- si
		tw.mu.Unlock()
	}(si)
}
func (tw *timeWheel) Invoke() {
	//invoke callback
	//赋值很关键，处理过程与curpos更新分离开
	curpos := tw.curpos
	//tw.mu.Lock()
	//defer tw.mu.Unlock()
	if tw.slot[curpos] != nil {
		go func(curpos int) {
		finish:
			for {
				select {
				case si := <-tw.slot[curpos]:
					if si.dead == false {
						//tw.SetInterval(si.timeout, si.do)
						tw.ReWheel(si, curpos)
						go si.do()
					}
				default:
					break finish
				}
			}
		}(curpos)
	}
}

func (tw *timeWheel) UpdatePos(roundTrigger func()) {
	if tw.curpos == tw.maxpos-1 {
		tw.curpos = 0
		if roundTrigger != nil {
			roundTrigger()
		}
	} else {
		tw.curpos = tw.curpos + 1
	}
}
func (tw *timeWheel) Start(roundTrigger func()) {
	t := time.NewTicker(time.Second * time.Duration(tw.interval))
	//if 0 slot got something
	tw.Trigger(nil)
	for {
		<-t.C
		tw.Trigger(roundTrigger)
	}
}

func (tw *timeWheel) Trigger(roundTrigger func()) {
	//先执行回调，再更新步数
	tw.Invoke()
	tw.UpdatePos(roundTrigger)
}
