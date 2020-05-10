# serror
Error package for golang services.

Error Package based on 
- [gRPC Error Handling](https://grpc.io/docs/guides/error/).
- [gRPC Error Codes](https://pkg.go.dev/google.golang.org/grpc/codes?tab=doc).
- [gRPC Error Status/Structure](https://pkg.go.dev/google.golang.org/grpc/status?tab=doc).
- [Google API's Error Details](https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto).

# Usage
```golang
err := InvalidArgs("Additional Message") //Create Invalid Args error with additional description
fmt.Println(err)                         //Print err info - code and message
//output: "rpc error: code = InvalidArgument desc = Additional Message"

h := ConvHTTP(err)                       //Convert to HTTP err format, can be used with json encoding
fmt.Printf("%+v\n", h.Rest)              //Print err in HTTP-Json Format
//output: "{Code:400 Desc:InvalidArgument Message:Additional Message Details:[]}"

//Creating Error with Details
//Create Field Violations and add it as Bad Request - Details in the error
v := Violation{
    Field:       "Field Name",                   //Pass field name that has violation
    Description: "mandatory field not provided", //Details of the violation
}
br := BadRequest{
    FieldViolations: []*Violation{&v}, //Pass the violations to BadRequest message
}

derr := InvalidArgs("Additional Message", &br) //Create Invalid Args error by passing addition description and violation details
fmt.Println(derr)                              //Print err info - code and message
//output: "rpc error: code = InvalidArgument desc = Additional Message"

fmt.Println(String(derr))                      //Print err with Details using String func
//output: "Code = InvalidArgument, Message = Additional Message, Details = [field_violations:<field:"Field Name" description:"mandatory field not provided" > ]"

h = ConvHTTP(derr)                             //Convert to HTTP err format, can be used with json encoding
fmt.Printf("%+v", h.Rest)                      //Print err in HTTP-Json Format
//output: "{Code:400 Desc:InvalidArgument Message:Additional Message Details:[field_violations:<field:"Field Name" description:"mandatory field not provided" > ]}"
```
Refer to package documentation for the list of errors and error details.
