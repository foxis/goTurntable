# goTurntable
Tiny-Go implementation of a turn-table controller

# 3D Printed design files

A simple turntable was used from Thingiverse: https://www.thingiverse.com/thing:4093452

ESP 8266 d1 mini was used for the brains. Motor was used the same as in the thing. 
Two touch sensors were added to control the speed of the rotation as well as LiPo batterry with USB-C charger.

ESP8266 was chosen in order to extend this turntable in the future to a fully autonomous 3D scanning rig.
The idea is that the ESP would be able to control a camera and report data via WiFi to a photogrammetry software.

Unfortunately TinyGo implementation for Espressif chips do not support neither WiFi nor BlueTooth at the moment.
