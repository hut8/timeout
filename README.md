# timeout
Timeout utilities for Go

[This blog post](https://tilde.town/~hut8/post/go-http-timeouts/) explains the need for this library.

## Use case
When performing any long-running IO operation, a timeout may be necessary. In Go, when using the standard library's `http.Client`, there is no way to implement a timeout that will deal with a stalled connection without timing out the connection after a certain duration. In other words, if you want a timeout while establishing a connection, that is possible, and if you want a timeout while downloading a file after a certain duration, regardless of its status, then that is possible too. But if you only want the timeout to occur after not receiving data for a certain amount of time, you have to write some code and use a `context.Context`. This library does that in a reusable way.  

## Example Usage

```
func DownloadFile(path, url string) error {
    file, err := os.Create(path)
    if err != nil {
      return err
    }

    res, err := http.Get(url)
    if err != nil {
      return err
    }
    defer res.Body.Close()

    ctx := context.Background()
    request, err := http.NewRequestWithContext(
        ctx, http.MethodGet, url, nil)

    tw := NewTimeoutWriter(ctx, 20*time.Second)
    defer tw.Cancel()

    sink := io.MultiWriter(file, tw)
    readbytes, err := io.Copy(sink, res.Body)

    if readbytes != res.ContentLength {
      // ... seems bad
    }

    if err != nil {
       // definitely bad
       return err
    }

    return nil
}
```