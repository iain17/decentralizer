# logger
Quick and simple logger in GO.

## Use

```go

logger.AddOutput(logger.Stdout{
  MinLevel: logger.INFO, //logger.DEBUG,
  Colored:  true,
})
logger.Info("test")
```
