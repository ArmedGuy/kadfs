import boto3
import botocore
from io import BytesIO

BUCKET_NAME = "kadfs"

session = boto3.session.Session()

s3 = boto3.resource(service_name='s3', endpoint_url="http://localhost:8080",
                    aws_access_key_id="bla", aws_secret_access_key="hehe")

try:
    kadfs = s3.Bucket(BUCKET_NAME)
    kadfs.put_object(Key="test/test/test.txt", Body='yolo'.encode("utf8"))
    b = BytesIO()
    kadfs.download_fileobj("test/test/test.txt", b)
    print(b.getvalue().decode('utf8'))
except botocore.exceptions.ClientError as e:
    if e.response["Error"]["Code"] == "404":
        print("The object does not exist.")
    else:
        raise
