ignition = provider("ignition", "1.1.0")

user = ignition.data.user()
user.name = "foo"
user.uid = 42
user.groups = ["foo", "bar"]
user.system = True

disk = ignition.data.disk()
disk.device = "/dev/sda"

root = disk.partition()
root.start = 2048
root.size = 4 * 1024 * 1024 

home = disk.partition()
home.start = root.size + root.start 
home.size = 4 * 1024 * 1024 
