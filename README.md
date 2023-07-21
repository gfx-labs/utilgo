# go


go utils v2


what is in this package?



### bufpool

make your own, or use, a global thread safe pool of *bytes.Buffers

use Get or Put when using a size annotated version

alternatively, use GetStd and PutStd to get the sync.Pool from std


### frand

this is a fork of frand with some random related utilities added to it


### graceful


graceful provides a way to run a application in main.


the first arg is the total timeout

you should put the goroutine loop that should be cancelled on the donech

the context that is available in that function will be cancelled after the stop signal is sent after timeout seconds

the cancel function will be run as well on signal end


### gun

this is gun, the opinionated configuration format.

tl;dr put `gun.Load(&SomeStruct)` somewhere

### ifx


some interfaces, idk. not really that useful


### lambda

some generic functions that perform basic functional operations on slices


### rdx

a global preconfigured version of miniredis, to be used as a temporary embedded store


### nkv

use any kv file to provide an easy access db