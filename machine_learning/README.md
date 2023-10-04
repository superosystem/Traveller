# traveloka-ocr-ml
### Overview

We built machine learning models for OCR Application to validate and automatically retrieving KTP (Kartu Tanda Penduduk)/Indonesia User's ID data without being inputted manually.

Machine learning part of the Apps consists of 2 main models:
- Object Detection model for detecting KTP position in real time and getting the bounding box for each class.
- Optical Character Recognition model for extracting text from bounding box produced by Object Detection model.
<br>
<p align="center">
    <img src="./optical-character-recognition/misc/diagram.png" alt="Application model diagram" width="550" style="vertical-align:middle">
</p>

### Dataset

We collected the data for both models from many sources as mentioned below:
- Google Image
- OpenSea
- Pinterest
- Bing Image

We collected the images only for learning purposes. Our goal is to collect the best and less noisy KTP.

We completely annotate images manually and with the help of utilities like [Roboflow](https://app.roboflow.com/). We only took 5 information from the KTP in the form of NIK, Name, Gender, Marital Status, and Nationality which would be used as a class. After the annotations are done, we implement data augmentation and data preprocessing such as image thresholding, binarization, skew correction, and noise removal.



