load('url', 'url')
load('assert.star', 'assert')


assert.eq(url.query_escape("/foo&bar qux"), "%2Ffoo%26bar+qux")
assert.eq(url.query_unescape("%2Ffoo%26bar+qux"), "/foo&bar qux")
assert.fails(lambda: url.query_unescape("%ssf"), 'invalid URL escape "%ss"')

assert.eq(url.path_escape("/foo&bar qux"), "%2Ffoo&bar%20qux")
assert.eq(url.path_unescape("%2Ffoo&bar%20qux"), "/foo&bar qux")
assert.fails(lambda: url.path_unescape("%ssf"), 'invalid URL escape "%ss"')

r = url.parse("http://qux:bar@bing.com/search?q=dotnet#foo")
assert.eq(r.scheme, "http")
assert.eq(r.opaque, "")
assert.eq(r.username, "qux")
assert.eq(r.password, "bar")
assert.eq(r.host, "bing.com")
assert.eq(r.path, "/search")
assert.eq(r.raw_query, "q=dotnet")
assert.eq(r.fragment, "foo")

r = url.parse("http://qux:@bing.com/search?q=dotnet#foo")
assert.eq(r.username, "qux")
assert.eq(r.password, "")

r = url.parse("http://qux@bing.com/search?q=dotnet#foo")
assert.eq(r.username, "qux")
assert.eq(r.password, None)

r = url.parse("http://bing.com/search?q=dotnet#foo")
assert.eq(r.username, None)
assert.eq(r.password, None)

assert.fails(lambda: url.parse("%ssf"), 'invalid URL escape "%ss"')
