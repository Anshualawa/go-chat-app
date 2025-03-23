package golang_social_chat

import (
	"encoding/json"
	"io"
	L "log/slog"
	"net/http"
	"time"
)

const (
	InTimeKey     = "InTime"
	RequestIDKey  = "RequestID"
	ResponseIDKey = "ResponseID"
)

// _Response is the base response with the essential fields of envery response.
type _Response struct {
	StatusCode          uint   `json:"statusCode"`
	Status              string `json:"status"`
	Message             string `json:"message"`
	RequestID           string `json:"requestId"`
	ResponseID          string `json:"responseId"`
	CreatedAt           string `json:"createdAt"` // YYYYY-MM-DDTHH:MM:SS.ddddd
	TimeToProcessMillis uint32 `json:"timeToProcessMillis"`
	Data                any    `json:"data,omitempty"`
}

// newResponse answers a response filled with the given basic details.
func newResponse(r *http.Request, status uint, msg string) *_Response {
	var br _Response

	// Compute the time taken to process the request.
	start, _ := r.Context().Value("InTime").(time.Time)
	br.TimeToProcessMillis = uint32(time.Since(start).Milliseconds())

	// Set the request and the response IDs.
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	br.RequestID = reqID
	resID, _ := r.Context().Value(ResponseIDKey).(string)
	br.ResponseID = resID

	// Set the current time up to microsecond precision.
	br.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000000")

	br.StatusCode = status
	br.Status = http.StatusText(int(status))
	br.Message = msg

	return &br
}

// SendError prepares and writes an error JSON response based on the given error data
func SendError(w http.ResponseWriter, r *http.Request, iErr error) {
	var _res *_Response
	switch _err := iErr.(type) {
	case *UserError:
		_res = newResponse(r, _err.Status, _err.Message)
		_res.Data = _err.Data
	case *SystemError:
		_res = newResponse(r, http.StatusInternalServerError, _err.Message)
		L.Error(_err.Message,
			L.String("requestId", _res.RequestID),
			L.String("responseId", _res.ResponseID),
			L.Int("status", int(_err.Status)),
			L.Any("data", _err.Data))
	default:
		_res = newResponse(r, http.StatusInternalServerError, _err.Error())
	}

	buf, err := json.Marshal(_res)
	if err != nil {
		L.Error("error marshaling error response!",
			L.String("requestId", _res.RequestID),
			L.String("responseId", _res.ResponseID),
			L.String("error", err.Error()))
		io.WriteString(w, `{"status":500,"message":"response marshaling failed"}`)
		return
	}
	w.Write(buf)
}

// SendSuccess prepares and write a success JSON response with possible result
// data, based on the given SuccessResponse structure.
func SendSuccess(w http.ResponseWriter, r *http.Request, res *Result) {
	_res := newResponse(r, res.status, res.message)
	_res.Data = res.data

	buf, err := json.Marshal(_res)
	if err != nil {
		L.Error("error marshaling success response!",
			L.String("requestId", _res.RequestID),
			L.String("responseId", _res.ResponseID),
			L.String("error", err.Error()))
		io.WriteString(w, `{"status":500, "message":"response marshaling failed"}`)
		return
	}
	w.Write(buf)
}
