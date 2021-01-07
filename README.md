1) Run `make build_all`
2) Go to /cmd/bin `cd cmd/bin`
3) Run app `./api`
4) Connect to `ws://localhost:8080/` via websocket
5) Send subscription `{"action": "subscribe", "symbols": ["XBTUSD", "ETHUSD"]}`
6) Receive data :)
