package url

import (
	"net/url"
	"sync"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	// ModuleName defines the expected name for this Module when used
	// in starlark's load() function, eg: load('io/ioutil', 'json')
	ModuleName = "url"

	pathEscapeFuncName    = "path_escape"
	pathUnescapeFuncName  = "path_unescape"
	queryEscapeFuncName   = "query_escape"
	queryUnescapeFuncName = "query_unescape"
	parseFuncName         = "parse"
)

var (
	once         sync.Once
	ioutilModule starlark.StringDict
)

// LoadModule loads the url module.
// It is concurrency-safe and idempotent.
//
//   outline: url
//     url parses URLs and implements query escaping.
//     path: url
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		ioutilModule = starlark.StringDict{
			"url": &starlarkstruct.Module{
				Name: "url",
				Members: starlark.StringDict{
					pathEscapeFuncName:    starlark.NewBuiltin(pathEscapeFuncName, PathEscape),
					pathUnescapeFuncName:  starlark.NewBuiltin(pathUnescapeFuncName, PathUnescape),
					queryEscapeFuncName:   starlark.NewBuiltin(queryEscapeFuncName, QueryEscape),
					queryUnescapeFuncName: starlark.NewBuiltin(queryUnescapeFuncName, QueryUnescape),
					parseFuncName:         starlark.NewBuiltin(parseFuncName, Parse),
				},
			},
		}
	})

	return ioutilModule, nil
}

// PathEscape escapes the string so it can be safely placed inside a URL path
// segment, replacing special characters (including /) with %XX sequences as
// needed.
//
//   outline: url
//     functions:
//       path_escape(s)
//         escapes the string so it can be safely placed inside a URL path
//         segment, replacing special characters (including /) with %XX
//         sequences as needed.
//         params:
//           s string
func PathEscape(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var s string

	err := starlark.UnpackArgs(pathEscapeFuncName, args, kwargs, "s", &s)
	if err != nil {
		return nil, err
	}

	return starlark.String(url.PathEscape(s)), nil
}

// PathUnescape does the inverse transformation of PathEscape, converting each
// 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB. It
// returns an error if any % is not followed by two hexadecimal digits.
// PathUnescape is identical to QueryUnescape except that it does not unescape
// '+' to ' ' (space).
//
//   outline: url
//     functions:
//       path_unescape(s)
//         does the inverse transformation of path_escape, converting each
//         3-byte encoded substring of the form "%AB" into the hex-decoded byte
//         0xAB. It returns an error if any % is not followed by two hexadecimal
//         digits. path_unescape is identical to query_unescape except that it
//         does not unescape '+' to ' ' (space).
//         params:
//           s string
func PathUnescape(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var s string

	err := starlark.UnpackArgs(pathUnescapeFuncName, args, kwargs, "s", &s)
	if err != nil {
		return nil, err
	}

	output, err := url.PathUnescape(s)
	return starlark.String(output), err
}

// QueryEscape escapes the string so it can be safely placed inside a URL query.
//
//   outline: url
//     functions:
//       path_escape(s)
//         escapes the string so it can be safely placed inside a URL query.
//         params:
//           s string
func QueryEscape(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var s string

	err := starlark.UnpackArgs(queryEscapeFuncName, args, kwargs, "s", &s)
	if err != nil {
		return nil, err
	}

	return starlark.String(url.QueryEscape(s)), nil
}

// QueryUnescape does the inverse transformation of QueryEscape, converting each
// 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB.
// It returns an error if any % is not followed by two hexadecimal digits.
//
//   outline: url
//     functions:
//       path_unescape(s)
//         does the inverse transformation of query_escape, converting each
//         3-byte encoded substring of the form "%AB" into the hex-decoded byte
//         0xAB. It returns an error if any % is not followed by two hexadecimal
//         digits.
//         params:
//           s string
func QueryUnescape(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var s string

	err := starlark.UnpackArgs(queryUnescapeFuncName, args, kwargs, "s", &s)
	if err != nil {
		return nil, err
	}

	output, err := url.QueryUnescape(s)
	return starlark.String(output), err
}

type sString = starlark.String

// URL represents a parsed URL (technically, a URI reference).
//
//   outline: url
//     types:
//       URL
//         Represents a parsed URL (technically, a URI reference).
//
//         fields:
//           scheme string
//           opaque string
//             Encoded opaque data.
//           username string
//             Username information.
//           password string
//             Password information.
//           host string
//             Host or host:port.
//           path string
//             Path (relative paths may omit leading slash).
//           raw_query string
//             Encoded query values, without '?'.
//           fragment string
//             Fragment for references, without '#'.
//
type URL struct {
	url url.URL
	sString
}

// Parse parses rawurl into a URL structure.
//
//   outline: url
//     functions:
//       parse(rawurl) URL
//         Parse parses rawurl into a URL structure.
//
//         params:
//           rawurl string
//              rawurl may be relative (a path, without a host) or absolute
//              (starting with a scheme). Trying to parse a hostname and path
//              without a scheme is invalid but may not necessarily return an
//              error, due to parsing ambiguities.
func Parse(
	thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var rawurl string
	err := starlark.UnpackArgs(parseFuncName, args, kwargs, "rawurl", &rawurl)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(rawurl)
	if err != nil {
		return starlark.None, err
	}

	return &URL{
		url:     *url,
		sString: starlark.String(url.String()),
	}, nil
}

func (u *URL) Attr(name string) (starlark.Value, error) {
	switch name {
	case "scheme":
		return starlark.String(u.url.Scheme), nil
	case "opaque":
		return starlark.String(u.url.Opaque), nil
	case "username":
		if u.url.User == nil {
			return starlark.None, nil
		}

		return starlark.String(u.url.User.Username()), nil
	case "password":
		if u.url.User == nil {
			return starlark.None, nil
		}

		password, provided := u.url.User.Password()
		if !provided {
			return starlark.None, nil
		}

		return starlark.String(password), nil
	case "host":
		return starlark.String(u.url.Host), nil
	case "path":
		return starlark.String(u.url.Path), nil
	case "raw_query":
		return starlark.String(u.url.RawQuery), nil
	case "fragment":
		return starlark.String(u.url.Fragment), nil
	}

	return nil, nil
}

func (*URL) AttrNames() []string {
	return []string{
		"scheme", "opaque", "username", "password", "host", "path",
		"raw_query", "fragment",
	}
}
