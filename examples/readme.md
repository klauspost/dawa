
Installing dependencies:

```go get github.com/klauspost/dawa/examples/...```


Running web example:
```
go run webapi.go
go run webapi-stream.go
go run bynavne-stream.go
go run query-adresse.go
go run query-list-reverse.go
go run query-list.go
go run query-adresse-geojson.go
```
Run file import examples:

1) Unpack files in examples.7z

2) Run the file you would like to test:
```
go run adgangsadresser-csv.go
go run adgangsadresser-json.go
go run adgangs-adresser-json.go
go run adresser-csv.go
go run adresser-json.go
go run postnumre-json.go
go run vejstykker-json.go
```