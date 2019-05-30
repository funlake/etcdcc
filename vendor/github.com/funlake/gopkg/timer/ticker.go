package timer

import "sync"

var (
	cron   *timer
	once   = sync.Once{}
	ticker *Ticker
)

func NewTicker() *Ticker {
	once.Do(func() {
		cron = NewTimer()
		cron.Ready()
		ticker = &Ticker{slotItems: make(map[string]*SlotItem)}
	})
	return ticker
}

type Ticker struct {
	slotItems map[string]*SlotItem
}

func (this *Ticker) Set(second int, key string, fun func()) {
	slotkey := string(second) + "_" + key
	if _, ok := this.slotItems[slotkey]; !ok {
		this.slotItems[slotkey] = cron.SetInterval(second, fun)
	}
}
func (this *Ticker) Stop(second int, key string) {
	slotkey := string(second) + "_" + key
	if _, ok := this.slotItems[slotkey]; ok {
		cron.StopInterval(this.slotItems[slotkey])
		delete(this.slotItems, slotkey)
	}
}
