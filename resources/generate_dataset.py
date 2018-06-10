from math import log2
import random


def main(path: str):
    # words.txt contains 466544
    # word_count contains all the powers of 2 from 8 to 2^log_2(466544), then 466543
    total_words = 466544
    word_count = [2**i for i in range(3, int(log2(total_words)))]
    # word_count.append(total_words - 1)

    word_list = []
    with open(path, "r") as f:
        word_list = f.readlines()  # load all the worlds

    # generate a random sample of indices in order to select unique words
    random_indices_samples = {i: random.sample(range(0, total_words - 1), i) for i in word_count}

    for k, v in random_indices_samples.items():
        # selecting k words indexes by the sample
        k_len_word_sample = [word_list[i].lower() for i in random_indices_samples[k]]
        random.shuffle(k_len_word_sample)
        with open('./dataset/w{}.txt'.format(k), 'w+') as f:
            f.writelines(k_len_word_sample)


if __name__ == "__main__":
    main('words.txt')
