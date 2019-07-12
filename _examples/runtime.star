load("encoding/base64", "base64")
load("http", "http")

dec = base64.encode("ascode is amazing")

msg = http.get("https://httpbin.org/base64/%s" % dec)
print(msg.body())