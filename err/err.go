// Package svcerr Defines the Error codes and functions to create Error messages.
// Only use the defined errors, so the list of errors are kept to minimal.
package svcerr

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	/* grpc error codes to http codes mapping */
	httpMap = map[codes.Code]int{
		codes.OK:                 http.StatusOK,
		codes.Canceled:           499, //Need to verify the http status code
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusGatewayTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.ResourceExhausted:  http.StatusTooManyRequests,
		codes.FailedPrecondition: http.StatusBadRequest,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
		codes.Unauthenticated:    http.StatusUnauthorized,
	}
)

//Aliases to googelapi errdetails, to avoid errdetails import in other packages.
type (
	//Violation - Used to represent request Field Violation
	Violation = errdetails.BadRequest_FieldViolation
	//BadRequest - Details on Request Violations
	BadRequest = errdetails.BadRequest
	//DebugInfo - Use for Sending Stack Entries
	DebugInfo = errdetails.DebugInfo
	//RetryInfo - Info for clients
	RetryInfo = errdetails.RetryInfo
	//QuotaFailure - Quota checks and failures
	QuotaFailure = errdetails.QuotaFailure
	//QuotaViolation -- Violations
	QuotaViolation = errdetails.QuotaFailure_Violation
	//PreConditionFailure  - Any Pre Conditiona failures and violations
	PreConditionFailure   = errdetails.PreconditionFailure
	PreConditionViolation = errdetails.PreconditionFailure_Violation
	//RequestInfo - meteadata about the request that clients can use
	RequestInfo = errdetails.RequestInfo
	//ResourceInfo - resources that is being accessed
	ResourceInfo = errdetails.ResourceInfo
	//Help Messages
	Help     = errdetails.Help
	HelpLink = errdetails.Help_Link
	//LocalizedMessage -- Localized error messages
	LocalizedMessage = errdetails.LocalizedMessage
)

//SvcErr - Service Error
type SvcErr struct {
	Rest      RestErr   `json:"Error"`
	LocalTime time.Time `json:"LocalTime"`
}

//RestErr - Error for REST API's
type RestErr struct {
	Code    int    `json:"Code"`
	Desc    string `json:"Desc"`
	Message string `json:"Message"`
	Details string `json:"Details"`
}

//Wrap - Function to wrap an error within a new error
func Wrap(text string, err error) error {
	return fmt.Errorf("%s:(%w)", text, err)
}

//Below are the Helper functions to create errors
//Args are error Message and optional reference to the details

//InvalidArgs - Creates a new error of type InvalidArgument
//will return err holding code for invalid arg and the message passed.
func InvalidArgs(msg string, details ...proto.Message) error {

	return newErr(codes.InvalidArgument, msg, details...)
}

//InternalErr Creates a new error of type Internal.
func InternalErr(msg string, details ...proto.Message) error {
	return newErr(codes.Internal, msg, details...)
}

//Unknown Creates a new error of type Unknown
func Unknown(msg string, details ...proto.Message) error {
	return newErr(codes.Unknown, msg, details...)
}

//NotFound Creates a new error of type NotFound
func NotFound(msg string, details ...proto.Message) error {
	return newErr(codes.NotFound, msg, details...)
}

//PermDenied Creates a new error of type Permission Denied
func PermDenied(msg string, details ...proto.Message) error {
	return newErr(codes.PermissionDenied, msg, details...)
}

//Canceled - Returns Request Canceled error
func Canceled(msg string, details ...proto.Message) error {
	return newErr(codes.Canceled, msg, details...)
}

//DeadlineExceeded - returns Deadline errors
func DeadlineExceeded(msg string, details ...proto.Message) error {
	return newErr(codes.DeadlineExceeded, msg, details...)
}

//AlreadyExists -- Returns error to client
func AlreadyExists(msg string, details ...proto.Message) error {
	return newErr(codes.AlreadyExists, msg, details...)
}

//ResourceExhausted - Returns Error to client
func ResourceExhausted(msg string, details ...proto.Message) error {
	return newErr(codes.ResourceExhausted, msg, details...)
}

//FailedPreCondition -- Returns Pre cond failure errors
func FailedPreCondition(msg string, details ...proto.Message) error {
	return newErr(codes.FailedPrecondition, msg, details...)
}

//Aborted - Request Aborted errors
func Aborted(msg string, details ...proto.Message) error {
	return newErr(codes.Aborted, msg, details...)
}

//OutOfRange - errors retured to client
func OutOfRange(msg string, details ...proto.Message) error {
	return newErr(codes.OutOfRange, msg, details...)
}

//Unimplemented - errors returned to client
func Unimplemented(msg string, details ...proto.Message) error {
	return newErr(codes.Unimplemented, msg, details...)
}

//DataLoss - errors to be returned to Client
func DataLoss(msg string, details ...proto.Message) error {
	return newErr(codes.DataLoss, msg, details...)
}

//Unavailable - errors to be returned to Client
func Unavailable(msg string, details ...proto.Message) error {
	return newErr(codes.Unavailable, msg, details...)
}

//Unauthenticated Creates a new error of type Unauthenticated
func Unauthenticated(msg string, details ...proto.Message) error {
	return newErr(codes.Unauthenticated, msg, details...)
}

//internal function creates new error and adds message details
func newErr(c codes.Code, message string, details ...proto.Message) error {
	err := status.Errorf(c, message)
	if len(details) != 0 {
		s := status.Convert(err)
		if errDetail, lerr := s.WithDetails(details...); lerr != nil {
			return err
		} else {
			return errDetail.Err()
		}
	}
	return err
}

//ConvHTTP - converts grpc error into Http Error structure
func ConvHTTP(err error) (se SvcErr) {
	st, ok := status.FromError(err)
	if ok {
		se.Rest = RestErr{
			Code:    httpMap[st.Code()],
			Desc:    st.Code().String(),
			Message: st.Message(),
			Details: fmt.Sprintf("%s", st.Details()),
		}
	} else {
		se.Rest = RestErr{
			Code: http.StatusInternalServerError,
			Desc: err.Error(),
		}
	}
	se.LocalTime = time.Now()
	return se
}

//IsValid - returns true if error is of rpc status type
func IsValid(err error) bool {
	_, ok := status.FromError(err)
	return ok
}

//String - returns string representation of rpc status error structure
func String(err error) string {
	s, ok := status.FromError(err)
	if ok {
		return fmt.Sprintf("Code = %s, Message = %s, Details = %s",
			s.Code().String(),
			s.Message(),
			s.Details(),
		)
	}
	return s.Message()
}

//Code - returns internal code from error
func Code(err error) codes.Code {
	s, ok := status.FromError(err)
	if ok {
		return s.Code()
	}
	return codes.Unknown
}
