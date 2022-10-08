module github.com/techyshishy/nirn-revenue-service

go 1.16

require (
	github.com/spf13/cobra v1.5.0
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261
	google.golang.org/genproto v0.0.0-20220908141613-51c1cc9bc6d0
	google.golang.org/protobuf v1.28.1
)

require github.com/inconshreveable/mousetrap v1.0.1 // indirect

replace github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 => github.com/techyshishy/gopher-lua v0.0.0-20220830034647-7a18b4e379e9
