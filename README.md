## How to Use?

## Step 1:
- Install and set up the [collector](https://github.com/evenscribe/evenscribe-collector).

## Step 2:
- ```go install https://github.com/evenscribe/evenscribe-go```
- setup the logger
```go
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
log.InfoS("Hello world", "log-attributes", map[string]string{"hello": "world"}, "trace-id", 12333, "span-id", 11111)
```

