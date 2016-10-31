radix
=====

radix is a small and performant implementation of a radix tree in Golang.

Benchmarks
----------

    BenchmarkRandomInsertSelf-4        2000000         572 ns/op
    BenchmarkRandomInsertGoRadix-4     2000000         695 ns/op
    BenchmarkFilesSelf-4                 50000       27972 ns/op
    BenchmarkFilesGoRadix-4              30000       40529 ns/op
    BenchmarkInsertFilesSelf-4          100000       17351 ns/op
    BenchmarkInsertFilesGoRadix-4        50000       29221 ns/op
