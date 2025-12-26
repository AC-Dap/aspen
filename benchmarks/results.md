12/25/2025
```
> go test .\benchmarks\ -bench BenchmarkRouter
goos: windows
goarch: amd64
pkg: aspen/benchmarks
cpu: Intel(R) Core(TM) Ultra 9 285H
BenchmarkRouter-16                  	    2312	    518764 ns/op	       0 B/op	       0 allocs/op
BenchmarkRouterWithMiddleware-16    	    1245	    956618 ns/op	  320012 B/op	   10000 allocs/op
```
