load("assert.star", "assert")

ignition = provider("ignition", "1.1.0")

user = ignition.data.user("test")
user.name = "foo"
user.uid = 42
user.groups = ["foo", "bar"]
user.system = True


assert.eq(user.__dict__, {
    "name": "foo", 
    "uid": 42,
    "groups": ["foo", "bar"],
    "system": True, 
})

disk = ignition.data.disk("foo")

root = disk.partition()
root.label = "root"
root.start = 2048
root.size = 4 * 1024 * 1024 

home = disk.partition()
home.label = "home"
home.start = root.size + root.start 
home.size = 4 * 1024 * 1024 

assert.eq(disk.__dict__, {
    "partition": [{
        "label": "root", 
        "start": 2048, 
        "size": 4194304
    }, {
        "start": 4196352, 
        "size": 4194304, 
        "label": "home"
    }]
})