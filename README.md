radix
=====

radix is a small and performant implementation of a radix tree in Golang.

Benchmarks
----------

    go test -run=XXX -bench=. -test.benchdir=~/Documents ./...
    BenchmarkRandomInsertSelf-4          1000000          3066 ns/op
    BenchmarkRandomInsertGoRadix-4        500000          3948 ns/op
    BenchmarkFilesSelf-4                     300       5143732 ns/op
    BenchmarkFilesGoRadix-4                  200       7923751 ns/op
    BenchmarkInsertFilesSelf-4               500       3251357 ns/op
    BenchmarkInsertFilesGoRadix-4            300       5282395 ns/op
