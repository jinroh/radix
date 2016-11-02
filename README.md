radix
=====

radix is a small and performant implementation of a radix tree in Golang.

Benchmarks
----------

    go test -run=XXX -bench=. -test.benchdir=~/Documents ./...
    BenchmarkRandomInsertSelf-4           500000          3156 ns/op
    BenchmarkRandomInsertGoRadix-4        500000          4082 ns/op
    BenchmarkFilesSelf-4                     200       7612412 ns/op
    BenchmarkFilesGoRadix-4                  200       8288134 ns/op
    BenchmarkInsertFilesSelf-4               300       4546510 ns/op
    BenchmarkInsertFilesGoRadix-4            300       5510951 ns/op
