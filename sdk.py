import os
import sys
import toml
from mcs import APIClient, BucketAPI


def bucket_list():
    mcs_api = APIClient(api_key, access_token, chain_name)
    bucket_client = BucketAPI(mcs_api)
    buckets = []
    for i in bucket_client.list_buckets():
        buckets.append({
            "name": i.bucket_name,
            "file_number": i.file_number,
            "size": i.size,
            "max_size": i.max_size,
        })
    return buckets


def check_bucket(bucket_name):
    mcs_api = APIClient(api_key, access_token, chain_name)
    bucket_client = BucketAPI(mcs_api)
    bucket = bucket_client.get_bucket(bucket_name)
    if bucket is None:
        bucket_client.create_bucket(bucket_name)


def upload_file(bucket_name, object_name, file_path):
    mcs_api = APIClient(api_key, access_token, chain_name)
    bucket_client = BucketAPI(mcs_api)
    paths = object_name.split("/")
    try:
        bucket_client.create_folder(bucket_name, paths[0])
    except :
        print()
    bucket_client.upload_file(bucket_name, object_name, file_path)


def get_file(bucket_name, object_name):
    mcs_api = APIClient(api_key, access_token, chain_name)
    bucket_client = BucketAPI(mcs_api)
    file_info = bucket_client.get_file(bucket_name, object_name).to_json()
    return file_info


if __name__ == "__main__":
    parsed_toml = toml.load(os.getcwd()+"/config.toml")
    chain_name = ""
    api_key = ""
    access_token = ""
    for key, value in parsed_toml.get('mcs').items():
        if key == "ChainName":
            chain_name = value
        if key == "ApiKey":
            api_key = value
        if key == "AccessToken":
            access_token = value

    args = len(sys.argv)
    if sys.argv[1] == "upload_file":
        bucketName, objectname, filePath = sys.argv[2:5]
        upload_file(bucketName, objectname, filePath)
    elif sys.argv[1] == "get_file":
        bucketName, objectname = sys.argv[2:4]
        print(get_file(bucketName, objectname))
    elif sys.argv[1] == "bucket_list":
        print(bucket_list())
    elif sys.argv[1] == "check_bucket":
        bucketName = sys.argv[2]
        check_bucket(bucketName)
    else:
        print("Not supported function name:", sys.argv[1])
