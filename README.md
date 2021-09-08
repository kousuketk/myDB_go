# myDB_go

### docker build, run
```
$ docker build -t mydb:1.0 ./
$ docker run --name mydb_container -itv /Users/takahashi/Go2/myDB_go:/app -it mydb:1.0 /bin/bash
or(control + P + Q)
$ docker start mydb_container
$ docker exec -it mydb_container bash
```


### test
```
go test -v
```