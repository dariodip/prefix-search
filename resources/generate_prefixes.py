import os
import random
import itertools
import string

MAX_PREFIXES = 1000


def generate_prefixes():

    prefixes = []
    prefixes_count = 0
    i = 1
    while prefixes_count < MAX_PREFIXES:
        new_prefixes = ["".join(p) for p in itertools.product(string.ascii_lowercase, repeat=i)]
        i += 1
        new_prefixes = random.sample(new_prefixes, len(new_prefixes) // ((2 ** i) - 1))
        prefixes_count = len(new_prefixes)
        prefixes += new_prefixes
    prefixes = random.sample(prefixes, MAX_PREFIXES)
    return prefixes


if __name__ == "__main__":
    prefixes = generate_prefixes()
    prefix_file_path = os.path.join(".", "prefixes", "pref1k.txt")
    lenss = [len(p) for p in prefixes]
    with open(prefix_file_path, "w+") as f:
        for prefix in prefixes:
            f.write("%s\n" % prefix)
