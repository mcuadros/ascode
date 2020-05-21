load("assert.star", "assert")

helm = tf.provider("helm", "1.0.0", "default")
helm.kubernetes.token = "foo"

# hcl
assert.eq(hcl(helm), "" +
'provider "helm" {\n' + \
'  alias   = "default"\n' + \
'  version = "1.0.0"\n\n' + \
'  kubernetes {\n' + \
'    token = "foo"\n' + \
'  }\n' + \
'}\n\n')

google = tf.provider("google", "3.16.0", "default")
sa = google.resource.service_account("sa")
sa.account_id = "service-account"

m = google.resource.storage_bucket_iam_member(sa.account_id+"-admin")
m.bucket = "main-storage"
m.role = "roles/storage.objectAdmin"
m.member = "serviceAccount:%s" % sa.email

addr = google.resource.compute_global_address("test")
addr.purpose = "VPC_PEERING"
addr.address_type = "INTERNAL"
addr.prefix_length = 16

# hcl with interpoaltion
assert.eq(hcl(google), "" + 
'provider "google" {\n' + \
'  alias   = "default"\n' + \
'  version = "3.16.0"\n' + \
'}\n' + \
'\n' + \
'resource "google_compute_global_address" "test" {\n' + \
'  provider      = google.default\n' + \
'  address_type  = "INTERNAL"\n' + \
'  prefix_length = 16\n' + \
'  purpose       = "VPC_PEERING"\n' + \
'}\n' + \
'\n' + \
'resource "google_service_account" "sa" {\n' + \
'  provider   = google.default\n' + \
'  account_id = "service-account"\n' + \
'}\n' + \
'\n' + \
'resource "google_storage_bucket_iam_member" "service-account-admin" {\n' + \
'  provider = google.default\n' + \
'  bucket   = "main-storage"\n' + \
'  member   = "serviceAccount:${google_service_account.sa.email}"\n' + \
'  role     = "roles/storage.objectAdmin"\n' + \
'}\n\n')

# hcl with prefixed provider
google = tf.provider("google", "3.16.0", "alias")
google.set_prefix(True)

sa = google.resource.service_account("sa")
sa.account_id = "service-account"
assert.eq(hcl(google), "" +
'provider "google" {\n' + \
'  alias   = "alias"\n' + \
'  version = "3.16.0"\n' + \
'}\n' + \
'\n' + \
'resource "google_service_account" "alias-sa" {\n' + \
'  provider   = google.alias\n' + \
'  account_id = "service-account"\n' + \
'}\n\n')