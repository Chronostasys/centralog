# CENTRALOG  
a simple central log library buit on top of zap.  
- simple 
- fast 
- reliable
- no security guarentee

## Usage

First, get it via:  
```bash
go get github.com/Chronostasys/centralog:v1.1.0
```

Then, start a server:
> Note: this lib should only be used in trusted environments, do not expose your log server to internet
```go
ls, err := logserver.CreateLogListener(&logserver.LogServerOptions{
    MongoUrl:           "mongodb://localhost:27018",
    Database:           "testlog",  // default log db
    Collection:         "zaplog",   // default log col
    ExpireAfterSeconds: 3600,
})
if err != nil {
    log.Fatal(err)
}
err = ls.Listen("0.0.0.0:8001")
if err != nil {
    log.Fatal(err)
}
```
Then, config the client:  
```go
log.InitLoggerWithOpt(zap.NewProductionConfig(), &log.LogOptions{
    Server:     "127.0.0.1:8001",
    Db:         "logtest",
    Collection: "logtest",
})
```
After that, log anything you want!  
```go
defer log.Sync() // flushes buffer, if any
log.Info("test").Any("testdata", map[string]interface{}{
    "hello": "world",
    "age":   18,
}).Log()
```
See [test](./test) directory for working examples


