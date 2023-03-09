for /f "skip=1" %%x in ('wmic os get localdatetime') do if not defined buildDate set buildDate=%x

go build -o .\cmd\bin\cc.exe -ldflags "-X main.buildDate=%buildDate%" .\cmd\
fyne package --os windows --exe .\bin\cc.exe --id test.kkonat.com --src .\cmd --release -icon Icon.png
cp .\cmd\bin\cc.exe \Users\Mieczu\go\bin\remotecc.exe
cp .\cmd\remotecc-config.yml \Users\Mieczu\go\bin