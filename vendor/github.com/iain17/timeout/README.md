# Timeout
This simple package allows you to call a function and specify a timeout. Once called it will hang until the method is either finished executing or the expire time has exceeded.

# Why
I tried numerous retry packages that included a timeout. However all either didn't actually inform the function that the callee had stopped waiting or actually stopped hanging.

# Usage
Here we have an example where we have a closure (basically an in-line function with no name that’s passed as a parameter to another function) that will expire in 10 seconds. Because the actual process takes 12 seconds.
```go
    print("Starting")
    timeout.Do(func(ctx context.Context) {
        time.Sleep(12 * time.Second)
    }, 10 * time.Second)
    print("Done")
```

Another example where we call **callMe** and it completes its task within 3 seconds. So it will print done after 3 seconds.
```go
    print("Starting")
    timeout.Do(callMe, 10 * time.Second)
    print("Done")
    
    func callMe(ctx context.Context) {
        time.Sleep(3 * time.Second)
    }
```

## How do I extract a result?
Let us say expensiveCallThatMightTakeTooLong that might do an external call. Just use the shared memory with a closure call.
```go
    print("Starting")
    result := 0
    timeout.Do(func(ctx context.Context) {
        result = expensiveCallThatMightTakeTooLong()
    }, 10 * time.Second)
    print("Done")
    print(result)
```