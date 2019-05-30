##### Intro
Simple timewheel pkg implement by golang
###### Install
go get github.com/funlake/gopkg/timer
###### Import
import github.com/funlake/gopkg/timer
```
cron := timer.NewTimer()
cron.Ready()
st := cron.SetInterval(3,func(){
    //do anything in each 3 seconds
})
//stop interval
cron.StopInterval(st)
```


