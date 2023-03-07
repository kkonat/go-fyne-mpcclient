### REMOTE CONTROL CENTER ###

(a MPD (linux music player deamon) client with fyne GUI)

--------------------


This is a front-end for my custom network music player device. The device is built on an old Raspberry Pi single board computer running ubuntu and a python script controlling its additional function. The raspbbery communicates over the serial port with an additional Atmega controller, which drives LCD and communicates over MIDI with Behringer DEQ2496 studio equalizer and power relay. The additional hardware lets me switch audio inputs (by telling DEQ2496 over MIDI to switch input) and turn the power on/off of the whole audio rig (by telling raspbery pi to command the relay). 

Hence this application features regular MPD controls + additional buttons for communication with my additional server (the python script in RPi) listening on another port.

I wrote this to try my skills in golang and to test out the following things:

- fyne cross-platform GUI & multi-threading
- application design (split between GUI state and HW interface)
- channel communication between various parts of the app
- graceful error handling (network disconnection)
- various build techniques

