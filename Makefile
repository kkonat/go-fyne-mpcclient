all:
	fyne package --os windows --src .\cmd --exe ..\bin\cc.exe --id test.kkonat.com --release -icon Icon.png
	cp .\cmd\bin\cc.exe \Users\Mieczu\go\bin\remotecc.exe