import boto3
import botocore
from io import BytesIO
import time
import sys

BUCKET_NAME = "kadfs"

session = boto3.session.Session()

s3 = boto3.resource(service_name='s3', endpoint_url="http://130.240.200.39",
                    aws_access_key_id="bla", aws_secret_access_key="hehe")

KEY = "test/test.txt"

kadfs = s3.Bucket(BUCKET_NAME)
if sys.argv[1] == "put":
    print("putting object")
    kadfs.put_object(Key=sys.argv[2], Body=sys.argv[3].encode("utf8"))
    print("put object")
elif sys.argv[1] == "get":
    b = BytesIO()
    print("getting object")
    kadfs.download_fileobj(sys.argv[2], b)
    print("get object")
    print(b.getvalue().decode('utf8'))
elif sys.argv[1] == "del":
    print("deleting object")
    s3.Object(BUCKET_NAME, sys.argv[2]).delete()
    print("deleted object")
else:
    print("unknown command")
