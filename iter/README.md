<div align="center">
    <h1>iter</h1>
    <a href="https://goreportcard.com/report/github.com/ameghdadian/x/iter">
     <img src="https://goreportcard.com/badge/github.com/ameghdadian/x/iter" height="20" alt="Go Report Card">
    </a>
    <a href="https://pkg.go.dev/github.com/ameghdadian/x/iter">
        <img src="https://pkg.go.dev/badge/github.com/ameghdadian/x/iter.svg" height="20" alt="Go Reference">
    </a>
    <a href="https://github.com/ameghdadian/x/actions/workflows/github-actions-iter.yaml/badge.svg">
        <img src="https://github.com/ameghdadian/x/actions/workflows/github-actions-iter.yaml/badge.svg" height="20" alt="CI">
    </a>
    <a href="#">
     <img src="https://img.shields.io/coverallsCoverage/github/ameghdadian/x" height="20" alt="Code Test Coverage">
    </a>

  <h3><em>iter provides tools to simplify working with Iterators(officially introduced in Go 1.23 release).</em></h3>
</div>


## Example

```go
    package main

    import (
        "fmt"
        "slices"
        "github.com/ameghdadian/iter"
    )

    func main() {
        a := []int{1,2,3}
        b := []int{4,5,6}

        // Combine arbitrary number of slices into an iterator
        it := iter.Concat(a,b)
        for v := range it {
            fmt.Println(v)
        }

        // Combine arbitrary number of iterators into a single iterator
        it1 := slices.Values(a)
        it2 := slices.Values(b)

        combined := iter.ConcatIter(it1, it2)
        for v := range combined {
            fmt.Println(v)
        }

        // ...
    }
```

## License
This is project is licensed under the Apache Version 2.0 License. See [LICENSE](LICENSE) for details.
