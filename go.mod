module github.com/f-secure-foundry/tamago-example

go 1.15

require (
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/f-secure-foundry/tamago v0.0.0-20200916193256-f70a57fb311b
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/mkevac/debugcharts v0.0.0-20191222103121-ae1c48aa8615
	github.com/shirou/gopsutil v2.20.8+incompatible // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sys v0.0.0-20200917073148-efd3b9a0ff20 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gvisor.dev/gvisor v0.0.0-20200917080942-a11061d78a58
)

replace gvisor.dev/gvisor => github.com/f-secure-foundry/gvisor v0.0.0-20200812210008-801bb984d4b1
