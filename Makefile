all:
	# BUILD=`date +%FT%T%z`
	BUILD=$(date +%Y-%m-%d)
	go build -o ./cmd/bin/cc.exe -ldflags "-X main.buildDate=${BUILD}" ./cmd/
	fyne package --os windows --exe .\bin\cc.exe --id test.kkonat.com --release -icon Icon.png
	cp .\cmd\bin\cc.exe \Users\Mieczu\go\bin\remotecc.exe