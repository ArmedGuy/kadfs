import boto3
import botocore

BUCKET_NAME = "kadfs"
KEY = "blabla/blabla.txt"

session = boto3.session.Session()

s3 = boto3.resource(service_name='s3', endpoint_url="http://localhost:8080",
                    aws_access_key_id="bla", aws_secret_access_key="hehe")

try:
    s3.Bucket(BUCKET_NAME).download_file(KEY, 'local.txt')
except botocore.exceptions.ClientError as e:
    if e.response["Error"]["Code"] == "404":
        print("The object does not exist.")
    else:
        raise
