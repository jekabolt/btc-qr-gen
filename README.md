# btc-qr-gen
HTTP server for generation qr code with encoded invoice and btc address. 
Also btc incoming tx watcher and reports on payments.

## Build Project Locally

To build project run:

```shell script
$ make build 
```

## Run Project Locally

To build project run:

```shell script
$ make run 
```

## API

### generate address invoice 
### /qr/{amount}/{meta}

```
amount = 0.01 
meta = base64 encoded PaymentInfo = {"country":"us","city":"nyc","address":"15 west 107","zip":"10025","fullName":"Eugene Boltenko","phone":"+1234567","email":"jekabolt@yahoo.com","totalAmount":"0.1","success":true}
```

```shell script
curl --request GET \
  --url http://localhost:8080/v1/qr/0.01/eyJjb3VudHJ5IjoidXMiLCJjaXR5IjoibnljIiwiYWRkcmVzcyI6IjE1IHdlc3QgMTA3IiwiemlwIjoiMTAwMjUiLCJmdWxsTmFtZSI6IkV1Z2VuZSBCb2x0ZW5rbyIsInBob25lIjoiKzEyMzQ1NjciLCJlbWFpbCI6Impla2Fib2x0QHlhaG9vLmNvbSIsInRvdGFsQW1vdW50IjoiMC4xIiwic3VjY2VzcyI6dHJ1ZX0= \
  --header 'x-api-key: kek'
```

### response
png with generated address meta and amount 

----

# History

## 0.1.0
