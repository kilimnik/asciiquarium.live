# asciiquarium.live

A hosted asciiquarium without any dependencies on any machine.

[![asciicast](https://asciinema.org/a/539009.svg)](https://asciinema.org/a/539009)

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

## Related Projects
* [parrot.live](https://github.com/hugomd/parrot.live)
* [ascii-live](https://github.com/hugomd/ascii-live)
* [asciiquarium original](https://robobunny.com/projects/asciiquarium/html/)