
import etcd3
import json

def restore(db_client, root, path=None):
    for key, raw_value in root.items():
        print("Write:", key, raw_value)
        value = json.dumps(raw_value)
        db_client.put(key, value)

def main():
    with open("./db.json", "r") as fd:
        db_document = json.load(fd)
        db_client = etcd3.client()
        restore(db_client, db_document)


if __name__ == "__main__":
    main()
