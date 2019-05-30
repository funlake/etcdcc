package timer

func NewTimer() *timer {
	return &timer{}
}

type timer struct {
	secondWheel *timeWheel
	minuteWheel *timeWheel
	hourWheel   *timeWheel
}

func (timer *timer) initialize() {
	timer.secondWheel = &timeWheel{
		interval: 1,
		maxpos:   60,
		curpos:   0,
		slot:     make(map[int]chan *SlotItem),
	}
	timer.minuteWheel = &timeWheel{
		interval: 60,
		maxpos:   60,
		curpos:   0,
		slot:     make(map[int]chan *SlotItem),
	}
	timer.hourWheel = &timeWheel{
		interval: 3600,
		maxpos:   24,
		curpos:   0,
		slot:     make(map[int]chan *SlotItem),
	}
}
func (timer *timer) SetInterval(timeout int, fun func()) *SlotItem {
	if timeout <= timer.minuteWheel.interval {
		return timer.secondWheel.SetInterval(timeout, fun)
	}
	if timeout > timer.minuteWheel.interval && timeout <= timer.hourWheel.interval {
		//用指针方能正确赋予sm.left值
		var sm *SlotItem
		sm = timer.minuteWheel.SetInterval(timeout, func() {
			//累计剩余left秒数，判断是否需要cross本次定时tick
			if sm.left > 0 && (sm.left+timer.minuteWheel.interval <= timeout) {
				sm.left = sm.left + timer.minuteWheel.interval
				//logs.Info(sm.left)
				//go to next round
				return
			}
			if timeout%timer.minuteWheel.interval > 0 {
				var ss *SlotItem
				//sm.left意思是定时任务离走完还剩余的秒数
				//mention : what happen if timeout - sm.left == 0 ?
				//需要修改Timewheel,先执行回调，再更新步数
				ss = timer.SetInterval((timeout-sm.left)%timer.minuteWheel.interval, func() {
					if sm != nil {
						//记录离此次分钟定时器结束还有多少秒
						//后续的定时间隔计算需考虑此sm.left
						if timer.secondWheel.curpos == 0 {
							//重置left,回到宇宙形成之初
							sm.left = 0
						} else {
							//得到调度后剩余步数，作为下次迭代累加基数
							sm.left = timer.minuteWheel.interval - ((timeout - sm.left) % timer.minuteWheel.interval)
						}
						//执行完后删除秒定时器，等待下次分钟级别的调度
						timer.secondWheel.StopInterval(ss)
					}
					fun()
				})
			} else {
				fun()
			}

		})
		return sm
	}
	if timeout > timer.hourWheel.interval {
		var sh *SlotItem
		sh = timer.hourWheel.SetInterval(timeout, func() {
			if sh.left > 0 && (sh.left+timer.hourWheel.interval <= timeout) {
				sh.left = sh.left + timer.hourWheel.interval
				//logs.Info(sh.left)
				//go to next round
				return
			}
			if timeout%timer.hourWheel.interval > 0 {
				var sm *SlotItem
				sm = timer.SetInterval((timeout-sh.left)%timer.hourWheel.interval, func() {
					if sm != nil {
						if timer.minuteWheel.curpos == 0 {
							//重置left,回到宇宙形成之初
							sm.left = 0
						} else {
							sh.left = timer.hourWheel.interval - ((timeout - sh.left) % timer.hourWheel.interval)
						}
						timer.minuteWheel.StopInterval(sm)
					}
					fun()
				})
			} else {
				fun()
			}
		})
		return sh
	}
	return nil
}
func (timer *timer) StopInterval(item *SlotItem) {
	if item.timeout <= timer.minuteWheel.interval {
		timer.secondWheel.StopInterval(item)
	}
	if item.timeout > timer.minuteWheel.interval && item.timeout <= timer.hourWheel.interval {
		timer.minuteWheel.StopInterval(item)
	}
	if item.timeout > timer.hourWheel.interval {
		timer.hourWheel.StopInterval(item)
	}
}
func (timer *timer) Ready() {
	timer.initialize()
	go timer.secondWheel.Start(func() {
		//every minutes pass
		timer.minuteWheel.UpdatePos(func() {
			//every hours pass
			timer.hourWheel.UpdatePos(func() {
				//nothing to do
			})
			timer.hourWheel.Invoke()
		})
		timer.minuteWheel.Invoke()
	})
}
