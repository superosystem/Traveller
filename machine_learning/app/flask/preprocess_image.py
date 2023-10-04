import glob
from typing import Union

import cv2  # type: ignore
import numpy as np
from PIL import Image
from scipy.ndimage import interpolation as inter


def convert_to_binary(image: np.array) -> np.array:
    preprocessed_image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    return preprocessed_image


def thresholding_adaptive(image: np.array) -> np.array:
    preprocessed_image = cv2.adaptiveThreshold(
        image, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY, 11, 2
    )
    return preprocessed_image


def find_score(arr: np.array, angle: np.array) -> Union[np.array, np.array]:
    data = inter.rotate(arr, angle, reshape=False, order=0)
    hist = np.sum(data, axis=1)
    score = np.sum((hist[1:] - hist[:-1]) ** 2)
    return hist, score


def rotate_image(image: np.array) -> np.array:
    delta = 1
    limit = 5
    angles = np.arange(-limit, limit + delta, delta)
    scores = []
    for angle in angles:
        hist, score = find_score(image, angle)
        scores.append(score)
    best_score = max(scores)
    best_angle = angles[scores.index(best_score)]
    print("Best angle: {}".format(best_angle))
    # correct skew
    data = inter.rotate(image, best_angle, reshape=False, order=0)
    preprocessed_image = (255 * data).astype("uint8")

    return preprocessed_image


def noise_removal(image: np.array) -> np.array:
    preprocessed_image = cv2.fastNlMeansDenoisingColored(image, None, 10, 10, 7, 15)
    return preprocessed_image


def thinning(image: np.array) -> np.array:
    image = np.abs(image - 255)
    kernel = np.ones((5, 5), np.uint8)
    preprocessed_image = cv2.erode(image, kernel, iterations=1)
    preprocessed_image = np.abs(preprocessed_image - 255)
    return preprocessed_image


def PIL_to_cv2(image: Image) -> np.array:
    np_image = np.array(image.convert("RGB"))
    return np_image[:, :, ::-1]


def cv2_to_PIL(image: np.array) -> Image:
    pil_image = Image.fromarray(image)
    return pil_image


def crop_image(image: Image, left: int, top: int, right: int, bottom: int) -> Image:
    return image.crop((left, top, right, bottom))


def thresholding_trunc(image: np.array) -> np.array:
    _, preprocessed_image = cv2.threshold(image, 127, 255, cv2.THRESH_TRUNC)
    return preprocessed_image


def thresholding_otsu(image: np.array) -> np.array:
    _, preprocessed_image = cv2.threshold(
        image, 127, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU
    )
    return preprocessed_image


def all_preprocessing(
    image: Image, left: int, top: int, right: int, bottom: int
) -> Image:
    np_image = crop_image(image, left, top, right, bottom)
    np_image = PIL_to_cv2(np_image)
    # np_image = noise_removal(np_image)
    np_image = convert_to_binary(np_image)
    # np_image = thinning(np_image)
    # np_image = thresholding_trunc(np_image)
    np_image = thresholding_otsu(np_image)
    final_result = cv2_to_PIL(np_image)

    return final_result


def main():
    list_image = glob.glob("./data/train_cut/*.jpg")
    for image_path in list_image:
        file_name = image_path.split("\\")[-1]
        print(file_name)
        # left_image, top_image, right_image, bottom_image = 115, 86, 390, 131
        image = Image.open(image_path)
        # image = crop_image(
        #     image,
        #     left=left_image,
        #     top=top_image,
        #     right=right_image,
        #     bottom=bottom_image,
        # )
        np_image = PIL_to_cv2(image)
        # np_image = noise_removal(np_image)
        np_image = convert_to_binary(np_image)
        # np_image = thresholding_trunc(np_image)
        np_image = thresholding_otsu(np_image)
        final_result = cv2_to_PIL(np_image)
        final_result = final_result.save(f"./data/train_preprocess/{file_name}")


if __name__ == "__main__":
    main()
