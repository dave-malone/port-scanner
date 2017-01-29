
```bash
port-scanner -timeout 1000 -hostname localhost -start 1 -end 49151
```


## Troubleshooting

`socket: too many open files`

Try increasing the max number of file descriptors: `ulimit -n 512`. Then, re-run this program with a larger `max-conns` setting.

On OSX, you may need to try something like [this](http://blog.dekstroza.io/ulimit-shenanigans-on-osx-el-capitan/)
