# asciiquarium.live

A hosted asciiquarium without any dependencies on any machine.

Try it yourself:
```
curl asciiquarium.live
```` 


You can specify the size too:
```
curl "http://asciiquarium.live?cols=100&rows=30"
```

Or use it fullscreen:
```
curl "http://asciiquarium.live?cols=$(tput cols)&rows=$(tput lines)"
```
