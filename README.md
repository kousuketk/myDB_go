# myDB_go

### docker build
```
$ docker build -t mydb:1.0 ./
$ docker run --name mydb_container -itv /Users/takahashi/Go2/myDB_go:/app -it mydb:1.0 /bin/bash
```

### run
```
# go run main.go 
```

### output
```
hello myDB
[{"tupleData":[{"type":1,"num":0,"str":"testInsert1"},{"type":1,"num":0,"str":"testInsert2"},{"type":1,"num":0,"str":"testInsert3"}]} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null} {"tupleData":null}]
```

### test
```
go test -v
```
