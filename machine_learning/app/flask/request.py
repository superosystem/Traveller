import json

import requests

with open("./testing/input.json", "r") as openfile:
    json_file = json.loads(openfile.read())

headers = {"Content-type": "application/json", "Accept": "text/plain"}
response = requests.post(
    "http://localhost:5000", data=json.dumps(json_file), headers=headers
)
print(response)
print(response.text)
