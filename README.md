## How to Use?

## Step 1:
- Install and set up the [collector](https://github.com/evenscribe/evenscribe-collector).

## Step 2:
- ```go install https://github.com/evenscribe/evenscribe-go```
- setup the logger
```go
    import (
        log "github.com/evenscribe/evenscribe-go"
    )

	err := log.NewLogger(
			log.NewOptions().
				WithDefaults().
				WithResourceAttributes(map[string]string{"evnviroment": "testing"}).
				WithServiceName("logger_test").
				WithPrintToStdout(false),
		)
```

## Step 3:
- Just get logging.

```go
log.Info("i got clicked...")
log.Fatal("oops... i broke")
```


```go
log.InfoS("Hello world", "log-attributes", map[string]string{"hello": "world"}, "trace-id", 12333, "span-id", 11111)
```

