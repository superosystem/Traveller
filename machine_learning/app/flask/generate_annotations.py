import os
import pandas as pd


def generate_annotation(df, annotation_base="/data", separator=" "):
    annotations = []
    for _, row in df.iterrows():
        indice, label = row
        if label == "-":
            continue
        left = os.path.join(annotation_base, "gambar%04d.jpg" % indice)
        right = label
        annotation = left + separator + right
        annotations.append(annotation)
    return annotations


def generate_annotation_file(annotations, file_name="annotations.txt", rewrite=False):
    with open(file_name, "w" if rewrite else "a") as file:
        for annotation in annotations:
            file.write(annotation + "\n")


def main():
    """Generate annotation for training from excel/spreadsheets"""

    # Configuration
    excel_path = "./labels.xlsx"
    annotation_base = "/train_name_preprocessed"
    sheet_name = "Name"
    nrows = None  # Change None for default
    annotations_separator = "\\t"

    df = pd.read_excel(excel_path, sheet_name=sheet_name, header=None, nrows=nrows)

    annotations = generate_annotation(df, annotation_base, annotations_separator)
    generate_annotation_file(annotations, "gt_name.txt", rewrite=True)
    # for a in annotations:
    #     print(a)


if __name__ == "__main__":
    main()
