## linux_tools_in_go

This is a simple pseudo-docker clone based on the talk by Liz Rice - Containers From Scratch: https://www.youtube.com/watch?v=8fi7uSYlOdc

It works similar to docker run: go run main.go run image <cmd> <params>

There were a few things to fix from the original presentation. It must be run by root to work, although Liz listed a presentation on how to modify it to work in user space. https://speakerdeck.com/lizrice/rootless-containers-from-scratch
