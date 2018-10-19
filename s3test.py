import boto3
import botocore
from io import BytesIO
import time

BUCKET_NAME = "kadfs"

session = boto3.session.Session()

s3 = boto3.resource(service_name='s3', endpoint_url="http://130.240.200.39",
                    aws_access_key_id="bla", aws_secret_access_key="hehe")

kadfs = s3.Bucket(BUCKET_NAME)
print("putting object")
kadfs.put_object(Key="test/test.txt", Body='yolo'.encode("utf8"))
print("put object")
b = BytesIO()
print("waiting 10 secs")
time.sleep(10)
print("getting object")
kadfs.download_fileobj("test/test.txt", b)
print("get object")
print(b.getvalue().decode('utf8'))
