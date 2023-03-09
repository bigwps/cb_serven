package utils

import "C"
import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"xc7/log"
	"xc7/utils/celpro"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"

	// cel官方就是这么写的。
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type CelOption struct {
	envoptions     []cel.EnvOption
	programoptions []cel.ProgramOption
}

// 声明自定义变量和函数。

func NewEnvOption() *CelOption {
	// 初始化结构体
	options := CelOption{}

	options.envoptions = []cel.EnvOption{
		// Container容器
		cel.Container("celpro"),
		cel.Types(
			&celpro.Request_C{},
			&celpro.Response_C{},
			&celpro.Reverse_C{},
			&celpro.UrlType_C{},
		),

		// 声明自定义变量
		// https://pkg.go.dev/github.com/google/cel-go/cel#Declarations
		// NewObjectType 为类型名称创建一个对象类型。

		cel.Declarations(
			// response -> Response_C
			decls.NewVar("response", decls.NewObjectType("celpro.Response_C")),
			decls.NewVar("request", decls.NewObjectType("celpro.Request_C")),
		),

		// 声明自定义函数
		cel.Declarations(
			// md5
			decls.NewFunction("md5", decls.NewOverload("md5_string",
				[]*exprpb.Type{decls.String}, decls.String)),
			// bcontains
			decls.NewFunction("bcontains", decls.NewInstanceOverload("bytes_bcontains_bytes",
				[]*exprpb.Type{decls.Bytes, decls.Bytes}, decls.Bool)),
			// contains
			decls.NewFunction("icontains", decls.NewInstanceOverload("icontains_string",
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool)),
			// wait
			decls.NewFunction("wait", decls.NewInstanceOverload("reverse_wait_int",
				[]*exprpb.Type{decls.Any, decls.Int}, decls.Bool)),
			// 随机小写字母
			decls.NewFunction("randomLowercase", decls.NewOverload("randomLowercase_int",
				[]*exprpb.Type{decls.Int}, decls.String)),
			// 正则匹配
			decls.NewFunction("bmatches", decls.NewInstanceOverload("string_bmatches_bytes",
				[]*exprpb.Type{decls.String, decls.Bytes}, decls.Bool)),
			// 随机数字
			decls.NewFunction("randomInt", decls.NewOverload("randomInt_int_int",
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Int)),
			//字符串切片
			decls.NewFunction("substr", decls.NewOverload("substr_string_int_int",
				[]*exprpb.Type{decls.String, decls.Int, decls.Int}, decls.String)),
			//base64
			decls.NewFunction("base64", decls.NewOverload("base64_string",
				[]*exprpb.Type{decls.String}, decls.String)),
			// //base64传入[]bytes
			decls.NewFunction("base64", decls.NewOverload("base64_bytes",
				[]*exprpb.Type{decls.Bytes}, decls.String)),
			decls.NewFunction("base64Decode", decls.NewOverload("base64Decode_string",
				[]*exprpb.Type{decls.String}, decls.String)),
			decls.NewFunction("base64Decode", decls.NewOverload("base64Decode_bytes",
				[]*exprpb.Type{decls.Bytes}, decls.String)),
			decls.NewFunction("urlencode", decls.NewOverload("urlencode_string",
				[]*exprpb.Type{decls.String}, decls.String)),
			decls.NewFunction("urlencode", decls.NewOverload("urlencode_bytes",
				[]*exprpb.Type{decls.Bytes}, decls.String)),
			decls.NewFunction("urldecode", decls.NewOverload("urldecode_string",
				[]*exprpb.Type{decls.String}, decls.String)),
			decls.NewFunction("urldecode", decls.NewOverload("urldecode_bytes",
				[]*exprpb.Type{decls.Bytes}, decls.String)),
		),
	}

	// 函数实现
	options.programoptions = []cel.ProgramOption{

		// 后续修改Function()来声明函数、重载
		cel.Functions(

			&functions.Overload{
				Operator: "md5_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to md5_string", value.Type())
					}
					return types.String(fmt.Sprintf("%x", md5.Sum([]byte(v))))
				},
			},
			&functions.Overload{
				Operator: "bytes_bcontains_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bcontains", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bcontains", rhs.Type())
					}
					// bytes.Contains检查一个字节切片（[]byte 类型）是否包含另一个字节切片。
					return types.Bool(bytes.Contains(v1, v2))
				},
			},
			&functions.Overload{
				Operator: "icontains_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to String", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to String", rhs.Type())
					}
					// 不区分大小写包含
					return types.Bool(strings.Contains(strings.ToLower(string(v1)), strings.ToLower(string(v2))))
				},
			},
			&functions.Overload{
				Operator: "reverse_wait_int",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					reverse, ok := lhs.Value().(*celpro.Reverse_C)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to 'wait'", lhs.Type())
					}
					timeout, ok := rhs.Value().(int64)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to 'wait'", rhs.Type())
					}
					return types.Bool(ReverseCheck(reverse, timeout))
				},
			},
			&functions.Overload{
				Operator: "randomLowercase_int",
				Unary: func(value ref.Val) ref.Val {
					n, ok := value.(types.Int)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to randomLowercase", value.Type())
					}
					return types.String(RandomLowercase(int(n)))
				},
			},
			&functions.Overload{
				Operator: "string_bmatches_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bmatch", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bmatch", rhs.Type())
					}
					ok, err := regexp.Match(string(v1), v2)
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.Bool(ok)
				},
			},
			&functions.Overload{
				Operator: "randomInt_int_int",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					from, ok := lhs.(types.Int)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to randomInt", lhs.Type())
					}
					to, ok := rhs.(types.Int)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to randomInt", rhs.Type())
					}
					min, max := int(from), int(to)
					//  Intn作为一个int返回半开区间[0,n]中的一个非负的伪随机数。
					return types.Int(rand.Intn(max-min) + 2)
				},
			},
			&functions.Overload{
				Operator: "substr_string_int_int",
				Function: func(values ...ref.Val) ref.Val {
					if len(values) == 3 {
						str, ok := values[0].(types.String)
						if !ok {
							return types.NewErr("invalid string to 'substr'")
						}
						start, ok := values[1].(types.Int)
						if !ok {
							return types.NewErr("invalid start to 'substr'")
						}
						length, ok := values[2].(types.Int)
						if !ok {
							return types.NewErr("invalid length to 'substr'")
						}
						runes := []rune(str)
						if start < 0 || length < 0 || int(start+length) > len(runes) {
							return types.NewErr("invalid start or length to 'substr'")
						}
						return types.String(runes[start : start+length])
					} else {
						return types.NewErr("too many arguments to 'substr'")
					}
				},
			},
			&functions.Overload{
				Operator: "base64_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_string", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString([]byte(v)))
				},
			},
			&functions.Overload{
				Operator: "base64_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_bytes", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString(v))
				},
			},
			&functions.Overload{
				Operator: "base64Decode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_string", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "base64Decode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_bytes", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "urlencode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_string", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urlencode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_bytes", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urldecode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_string", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "urldecode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_bytes", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			//	--------
		),
	}
	return &options
}

// cel.Library

func (c *CelOption) CompileOptions() []cel.EnvOption {
	return c.envoptions
}

func (c *CelOption) ProgramOptions() []cel.ProgramOption {
	return c.programoptions
}

// set全局变量，注册到env中.

func (c *CelOption) Update_Set_Options(args map[string]string) {
	for key, v := range args {
		var d *exprpb.Decl
		if strings.HasPrefix(v, "randomInt") {
			d = decls.NewVar(key, decls.Int)
		} else if strings.HasPrefix(v, "newReverse") {
			// key="reverse"
			d = decls.NewVar(key, decls.NewObjectType("celpro.Reverse_C"))
			log.InfoLog("注册reverse")
		} else {
			d = decls.NewVar(key, decls.String)
		}
		c.envoptions = append(c.envoptions, cel.Declarations(d))
	}
}

// 将r0 && r1这类的表达式结果，true or false加载到env

func (c *CelOption) UpdateFunctionOptions(rn string, isTrue ref.Val) {
	dec := decls.NewFunction(rn, decls.NewOverload(rn, []*exprpb.Type{}, decls.Bool))
	c.envoptions = append(c.envoptions, cel.Declarations(dec))
	function := &functions.Overload{
		Operator: rn,
		Function: func(values ...ref.Val) ref.Val {
			return isTrue
		},
	}
	c.programoptions = append(c.programoptions, cel.Functions(function))
}

func Evaluate(env *cel.Env, v string, params map[string]interface{}) (ref.Val, error) {

	// Compile结合了解析和检查阶段的CEL程序编译，返回Ast和iss。
	ast, iss := env.Compile(v)
	if iss.Err() != nil {
		log.ErrorLog("Compile解析出错!!")
		return nil, iss.Err()
	}
	// Program在环境（Env）中生成一个可评估的Ast实例。`cel.ProgramOption`选项进行自定义
	prg, err := env.Program(ast)
	if err != nil {
		log.ErrorLog("Program生成出错!!")
		return nil, err
	}
	// Eval返回Ast和环境对输入变量的评估结果。
	out, _, err := prg.Eval(params)
	if err != nil {
		log.ErrorLog("Evaluation err!!!")
		return nil, err
	}
	return out, nil
}
