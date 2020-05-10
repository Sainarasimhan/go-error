package svcerr

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWrap(t *testing.T) {
	type args struct {
		text string
		err  error
	}
	tests := []struct {
		name       string
		args       args
		wantErrMsg string
	}{
		{
			name: "Wrap Error",
			args: args{
				err:  NotFound("Not Found"),
				text: "Outer Text",
			},
			wantErrMsg: "Outer Text" + ":" + "(rpc error: code = NotFound desc = Not Found)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wrap(tt.args.text, tt.args.err); err != nil {
				if err.Error() != tt.wantErrMsg {
					t.Errorf("Wrap() have error = %s, wantErr %s", err, tt.wantErrMsg)
				}
			} else {
				t.Errorf("Wrap() error = %v, wantErr", err)
			}
		})
	}
}

func TestErrorFuncs(t *testing.T) {
	type args struct {
		msg     string
		details []proto.Message
	}

	var (
		badReq = BadRequest{FieldViolations: []*Violation{
			{
				Field:       "Test",
				Description: "Test",
			},
		}}
	)
	tests := []struct {
		fn          func(string, ...proto.Message) error
		name        string
		args        args
		wantErrcode codes.Code
		wantErrmsg  string
		wantDetails []proto.Message
	}{
		{
			fn:   InvalidArgs,
			name: "Invalid Args",
			args: args{
				msg:     "Invalid msg",
				details: []proto.Message{&badReq},
			},
			wantErrcode: codes.InvalidArgument,
			wantErrmsg:  "Invalid msg",
			wantDetails: []proto.Message{&badReq},
		},
		{
			fn:          InternalErr,
			name:        "Internal Error",
			args:        args{msg: "Internal Error"},
			wantErrcode: codes.Internal,
			wantErrmsg:  "Internal Error",
		},
		{
			fn:          NotFound,
			name:        "Not Found Error",
			args:        args{msg: "Not Found"},
			wantErrcode: codes.NotFound,
			wantErrmsg:  "Not Found",
		},
		{
			fn:          Unauthenticated,
			name:        "Unauthenticated Error",
			args:        args{msg: "Unauthenticated"},
			wantErrcode: codes.Unauthenticated,
			wantErrmsg:  "Unauthenticated",
		},
		{
			fn:          Unimplemented,
			name:        "Unimplemented Error",
			args:        args{msg: "Unimplemented"},
			wantErrcode: codes.Unimplemented,
			wantErrmsg:  "Unimplemented",
		},
		{
			fn:          DeadlineExceeded,
			name:        "DealLine Error",
			args:        args{msg: "DeadLineExceeded"},
			wantErrcode: codes.DeadlineExceeded,
			wantErrmsg:  "DeadLineExceeded",
		},
		{
			fn:          Unknown,
			name:        "Unknown Error",
			args:        args{msg: "Unknown"},
			wantErrcode: codes.Unknown,
			wantErrmsg:  "Unknown",
		},
		{
			fn:          NotFound,
			name:        "NotFound Error",
			args:        args{msg: "NotFound"},
			wantErrcode: codes.NotFound,
			wantErrmsg:  "NotFound",
		},
		{
			fn:          PermDenied,
			name:        "PermDenied Error",
			args:        args{msg: "PermDenied"},
			wantErrcode: codes.PermissionDenied,
			wantErrmsg:  "PermDenied",
		},
		{
			fn:          Canceled,
			name:        "Canceled Error",
			args:        args{msg: "Canceled"},
			wantErrcode: codes.Canceled,
			wantErrmsg:  "Canceled",
		},
		{
			fn:          AlreadyExists,
			name:        "Already Exists Error",
			args:        args{msg: "Already Exists"},
			wantErrcode: codes.AlreadyExists,
			wantErrmsg:  "Already Exists",
		},
		{
			fn:          ResourceExhausted,
			name:        "Resource Exhausted Error",
			args:        args{msg: "Resource Exhausted"},
			wantErrcode: codes.ResourceExhausted,
			wantErrmsg:  "Resource Exhausted",
		},
		{
			fn:          FailedPreCondition,
			name:        "FailedPreCondition Error",
			args:        args{msg: "FailedPreCondition"},
			wantErrcode: codes.FailedPrecondition,
			wantErrmsg:  "FailedPreCondition",
		},
		{
			fn:          Aborted,
			name:        "Aborted Error",
			args:        args{msg: "Aborted"},
			wantErrcode: codes.Aborted,
			wantErrmsg:  "Aborted",
		},
		{
			fn:          OutOfRange,
			name:        "OutOfRange",
			args:        args{msg: "OutOfRange"},
			wantErrcode: codes.OutOfRange,
			wantErrmsg:  "OutOfRange",
		},
		{
			fn:          DataLoss,
			name:        "DataLoss",
			args:        args{msg: "DataLoss"},
			wantErrcode: codes.DataLoss,
			wantErrmsg:  "DataLoss",
		},
		{
			fn:          Unavailable,
			name:        "Unavailable",
			args:        args{msg: "Unavailable"},
			wantErrcode: codes.Unavailable,
			wantErrmsg:  "Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(tt.args.msg, tt.args.details...); err != nil {
				s := status.Convert(err)
				if s.Code() != tt.wantErrcode { //Verify if expected error code is setup
					t.Errorf("(%s) failed due to Error Code. have = %d, want = %d", tt.name, s.Code(), tt.wantErrcode)
				}
				if s.Message() != tt.wantErrmsg { //Verify if expected error message is setup
					t.Errorf("(%s) failed due to error message, have = %s, want = %s", tt.name, s.Message(), tt.wantErrmsg)
				}
				d := s.Details()
				for i, v := range d { //Verify if the details passed are properly setup
					val := v.(proto.Message)
					if val.String() != tt.wantDetails[i].String() {
						t.Errorf("(%s) failed due to details, have = %v, want = %v", tt.name, val, tt.wantDetails[i])
					}
				}
			} else {
				t.Errorf("(%s) failed with error = %v, want valid Err", tt.name, err)
			}
		})
	}
}

func ExampleInvalidArgs() {

	err := InvalidArgs("Additional Message") //Create Invalid Args error with additional description
	fmt.Println(err)                         //Print err info - code and message
	h := ConvHTTP(err)                       //Convert to HTTP err format, can be used with json encoding
	fmt.Printf("%+v\n", h.Rest)              //Print err in HTTP-Json Format

	fmt.Println("Creating Error with Details")
	v := Violation{
		Field:       "Field Name",                   //Pass field name that has violation
		Description: "mandatory field not provided", //Details of the violation
	}
	br := BadRequest{
		FieldViolations: []*Violation{&v}, //Pass the violations to BadRequest message
	}

	derr := InvalidArgs("Additional Message", &br) //Create Invalid Args error by passing addition description and violation details
	fmt.Println(derr)                              //Print err info - code and message
	fmt.Println(String(derr))                      //Print err with Details using String func
	h = ConvHTTP(derr)                             //Convert to HTTP err format, can be used with json encoding
	fmt.Printf("%+v", h.Rest)                      //Print err in HTTP-Json Format
	// Output:
	// rpc error: code = InvalidArgument desc = Additional Message
	// {Code:400 Desc:InvalidArgument Message:Additional Message Details:[]}
	// Creating Error with Details
	// rpc error: code = InvalidArgument desc = Additional Message
	// Code = InvalidArgument, Message = Additional Message, Details = [field_violations:<field:"Field Name" description:"mandatory field not provided" > ]
	// {Code:400 Desc:InvalidArgument Message:Additional Message Details:[field_violations:<field:"Field Name" description:"mandatory field not provided" > ]}
}

func TestConvHTTP(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		args   args
		wantSe SvcErr
	}{
		{
			name: "Invalid Args",
			args: args{err: InvalidArgs("Invalid Args")},
			wantSe: SvcErr{
				Rest: RestErr{
					Code:    http.StatusBadRequest,
					Desc:    codes.InvalidArgument.String(),
					Message: "Invalid Args",
					Details: "[]", //expect empty slice if no details provided
				},
			},
		},
		{
			name: "Not Found",
			args: args{err: NotFound("Not Found")},
			wantSe: SvcErr{
				Rest: RestErr{
					Code:    http.StatusNotFound,
					Desc:    codes.NotFound.String(),
					Message: "Not Found",
					Details: "[]", //expect empty slice if no details provided
				},
			},
		},
		{
			name: "Generic Error",
			args: args{err: errors.New("Generic")},
			wantSe: SvcErr{
				Rest: RestErr{
					Code: http.StatusInternalServerError,
					Desc: "Generic",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSe := ConvHTTP(tt.args.err); !reflect.DeepEqual(gotSe.Rest, tt.wantSe.Rest) {
				t.Errorf("ConvHTTP() = %v, want %v", gotSe, tt.wantSe)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid gRPC status error",
			args: args{err: ResourceExhausted("Resource Exhausted")},
			want: true,
		},
		{
			name: "InValid gRPC status error",
			args: args{err: errors.New("Test Error")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.args.err); got != tt.want {
				t.Errorf("IsValid(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "String",
			args: args{InvalidArgs("test")},
			want: "Code = InvalidArgument, Message = test, Details = []",
		},
		{
			name: "Generic String",
			args: args{err: errors.New("test")},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.err); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
