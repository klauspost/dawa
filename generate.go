package dawa

//go:generate codecgen -o values.generated.go postnummer.go awstime.go adresse.go adgangsadressse.go vejstykke.go

// Using this will speed up decoding about 20%
// To use, run this in the current directory:
//
// go get "github.com/ugorji/go/codec/codecgen"
// go generate
