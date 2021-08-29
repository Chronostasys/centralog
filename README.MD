# CENTRALOG  
a simple central log library buit on top of zap.  
- simple 
- fast 
- reliable
- no security guarentee

## Usage

First, start a server:
> Note: this lib should only be used in trusted environments, do not expose your log server to internet
```go
ls, err := logserver.CreateLogListener(&logserver.LogServerOptions{
    MongoUrl:           "mongodb://localhost:27018",
    Database:           "testlog",
    Collection:         "zaplog",
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
log.InitLogger(zap.NewProductionConfig(), "127.0.0.1:8001")
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


