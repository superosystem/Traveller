import json


def main():
    with open("./input.json", "r") as openfile:
        file_json = json.load(openfile)

    print(file_json["image"])


if __name__ == "__main__":
    main()
