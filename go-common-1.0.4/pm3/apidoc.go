package pm3

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/eolinker/go-common/auto"

	"github.com/gin-gonic/gin"
)

type readHandler func(ginCtx *gin.Context) (interface{}, error)

type Response struct {
	Data     map[string]any `json:"data,omitempty"`
	Code     int            `json:"code"`
	Success  any            `json:"success,omitempty"`
	Message  string         `json:"msg,omitempty"`
	Total    int64          `json:"tatal,omitempty"`
	PageSize int            `json:"pageSize,omitempty"`
	Page     int            `json:"page,omitempty"`
}
type IApiError interface {
	error
	Code() int
	Message() string
	Success() any
}

type apiDoc struct {
	Method      string
	Path        string
	IN          []string
	OUT         []string
	HandlerFunc any

	restfulSet map[string]struct{}
}

func CreateApiWidthDoc(method string, path string, IN []string, OUT []string, handlerFunc any) Api {
	return Gen(&apiDoc{Method: method, Path: path, IN: IN, OUT: OUT, HandlerFunc: handlerFunc})
}
func CreateApiSimple(method string, path string, handler gin.HandlerFunc) Api {
	return &formApi{
		method:  method,
		path:    path,
		handler: handler,
	}
}
func CreateApiWidthError(method string, path string, handlerFunc any) Api {
	return Gen(&apiDoc{Method: method, Path: path, IN: nil, OUT: nil, HandlerFunc: handlerFunc})
}

// HandlerFunc 为实现方法, 方法签名需要和 IN 对应
// IN 的格式为 {type}:{key}, 框架会自动转换类型, type 为 header, query,
//  type 类型有:
//	context=> gin.Context,
//  header:{key} = > http的请求头
//  query:{key} => query参数
//  restful:{key}  => restful参数, 框架会检查key是否为path中的有效restful参数, restful 参数可以用 :{key} 来表示
//  body:{encode} => 输入bogy, 使用 {encode}对body进行目标类型的反序列化,默认为 json
// HandlerFunc的返回值至少要有一个,最后一个必须是 error, OUT 字段标识返回的值在 response.body下data下的字段名,
// 如想 IN=["context",":id,"body"] ,OUT = ["list"], 适用签名  func(ctx *gin.Context,id int,body *BodyStruct) (list []*Result,err error)

// Handler 构造api handler
func (a *apiDoc) Handler() gin.HandlerFunc {

	readHandlers, err := a.check()
	if err != nil {
		log.Fatal(fmt.Errorf("check api [%s %s] error:%w", a.Method, a.Path, err))
	}

	fv := reflect.ValueOf(a.HandlerFunc)
	return func(context *gin.Context) {
		in := make([]reflect.Value, 0, len(readHandlers))
		for _, rh := range readHandlers {
			v, err := rh(context)
			if err != nil {
				context.JSON(200, &Response{
					Data:    nil,
					Code:    -1,
					Success: "fail",
					Message: fmt.Sprintf("invald request:%s", err.Error()),
				})
				return
			}
			in = append(in, reflect.ValueOf(v))
		}

		rs := fv.Call(in)

		errV := rs[len(rs)-1]
		if !errV.IsNil() {

			err := errV.Interface().(error)
			var ae IApiError
			if errors.As(err, &ae) {
				context.JSON(200, &Response{
					Data:    nil,
					Code:    ae.Code(),
					Success: ae.Success(),
					Message: ae.Message(),
				})
			} else {
				context.JSON(200, &Response{

					Code:    -1,
					Success: "fail",
					Message: err.Error(),
				})
			}
			return
		}
		resp := &Response{
			Data:    make(map[string]any),
			Code:    0,
			Success: "success",
			Message: "",
		}
		for i, field := range a.OUT {

			rv := rs[i]
			if field != "" {
				value := rv.Interface()
				auto.CompleteLabels(context, value)
				resp.Data[field] = value
			}

		}
		context.JSON(200, resp)

	}
}

func (a *apiDoc) check() ([]readHandler, error) {
	vt := reflect.TypeOf(a.HandlerFunc)
	if vt.Kind() != reflect.Func {
		return nil, fmt.Errorf("api handler must func but get :%s", vt.Kind())
	}

	if vt.NumIn() != len(a.IN) {
		return nil, fmt.Errorf("api hanlder[%s] need %d input arg but require  %d", vt.String(), vt.NumIn(), len(a.IN))
	}

	if vt.NumOut() < 1 {
		return nil, fmt.Errorf("api handler at least one  error out :[%s]", vt.String())
	}

	if vt.Out(vt.NumOut()-1).String() != "error" {
		return nil, fmt.Errorf("api handler the last one out must be error but not on:[%s]", vt.String())
	}

	if vt.NumOut() != len(a.OUT)+1 {
		return nil, fmt.Errorf("api handler[%s] has %d out but require %d", vt.String(), vt.NumOut(), len(a.OUT)+1)
	}
	readHanlers := make([]readHandler, 0, len(a.IN))

	for i, in := range a.IN {
		handler, err := a.parseReadHandler(in, vt.In(i))
		if err != nil {
			return nil, err
		}
		readHanlers = append(readHanlers, handler)
	}
	return readHanlers, nil

}
func (a *apiDoc) parseReadHandler(pattern string, vt reflect.Type) (readHandler, error) {

	ps := strings.Split(pattern, ":")
	if len(ps) == 1 {
		ps = append(ps, "")
	}
	tn := ps[0]
	switch strings.ToLower(tn) {
	case "", "rest", "path", "restful":
		var name = ps[1]
		if name != "" {
			if !a.checkRestFul(name) {
				return nil, fmt.Errorf("restful parame not exist[%s] for api require", name)
			}
			vparse, err := getStringPareHandler(vt)
			if err != nil {
				return nil, err
			}
			return func(ginCtx *gin.Context) (interface{}, error) {
				v := ginCtx.Param(name)
				return vparse(v)
			}, nil
		}
	case "body":
		encodeType := strings.ToLower(ps[1])
		if encodeType == "json" || encodeType == "" {
			return func(ginCtx *gin.Context) (interface{}, error) {
				if vt.Kind() == reflect.Ptr {
					vt = vt.Elem()
				}
				value := reflect.New(vt).Interface()
				err := ginCtx.BindJSON(value)
				if err != nil {
					return nil, err
				}
				err = auto.SearchIDCheck(ginCtx, value)
				if err != nil {
					return nil, err
				}
				return value, nil
			}, nil
		}
		if encodeType == "yaml" {
			return func(ginCtx *gin.Context) (interface{}, error) {
				value := reflect.New(vt).Interface()
				err := ginCtx.BindYAML(value)
				if err != nil {
					return nil, err
				}
				return value, nil
			}, nil
		}

	case "context":
		return contextReadHandler, nil
	case "header":
		var name = ps[1]
		if name != "" {
			vparse, err := getStringPareHandler(vt)
			if err != nil {
				return nil, err
			}
			return func(ginCtx *gin.Context) (interface{}, error) {
				v := ginCtx.GetHeader(name)
				return vparse(v)
			}, nil
		}
	case "query":
		var name = ps[1]
		if name != "" {
			vparse, err := getStringPareHandler(vt)
			if err != nil {
				return nil, err
			}
			return func(ginCtx *gin.Context) (interface{}, error) {
				v, _ := ginCtx.GetQuery(name)
				return vparse(v)
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid pattern  for input reader [%s]", pattern)

}

func contextReadHandler(ginCtx *gin.Context) (interface{}, error) {
	return ginCtx, nil
}
func getStringPareHandler(vt reflect.Type) (func(v string) (interface{}, error), error) {

	if vt.Kind() == reflect.Ptr {
		return nil, fmt.Errorf("not support type:%s", vt.String())
	}

	switch vt.Kind() {
	case reflect.String:
		return parseString, nil
	case reflect.Int:
		return parseInt[int], nil
	case reflect.Int8:
		return parseInt[int8], nil
	case reflect.Int16:
		return parseInt[int64], nil
	case reflect.Int32:
		return parseInt[int32], nil

	case reflect.Int64:
		return parseInt[int64], nil

	case reflect.Uint:
		return parseUInt[uint], nil
	case reflect.Uint8:
		return parseUInt[uint8], nil

	case reflect.Uint16:
		return parseUInt[uint16], nil

	case reflect.Uint32:
		return parseUInt[uint32], nil

	case reflect.Uint64:
		return parseUInt[uint64], nil

	case reflect.Bool:
		return parseBool, nil
	case reflect.Float32:
		return parseFloat32, nil
	case reflect.Float64:
		return parseFloat64, nil
	default:

	}

	return nil, fmt.Errorf("not support kind %s for type:%s", vt.Kind(), vt.String())

}

func (a *apiDoc) checkRestFul(name string) bool {
	if a.restfulSet == nil {
		a.restfulSet = parseRestfulParams(a.Path)
	}
	_, ok := a.restfulSet[name]
	return ok

}
func parseRestfulParams(path string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, n := range strings.Split(path, "/") {
		if len(n) > 1 {
			if n[0] == ':' || n[0] == '*' {
				set[n[1:]] = struct{}{}
			}
		}
	}
	return set
}

type Integer interface {
	int | int8 | int16 | int32 | int64
}

func parseString(s string) (interface{}, error) {
	return s, nil
}
func parseInt[T Integer](v string) (interface{}, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

type UInteger interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func parseUInt[T UInteger](v string) (interface{}, error) {
	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}
func parseBool(v string) (interface{}, error) {
	return strconv.ParseBool(v)
}

func parseFloat32(v string) (interface{}, error) {
	vo, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return nil, err
	}
	return float32(vo), err
}
func parseFloat64(v string) (interface{}, error) {
	return strconv.ParseFloat(v, 64)
}
