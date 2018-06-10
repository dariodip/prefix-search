#!/usr/bin/python3

import os
import json
import pandas as pd
import subprocess
import tempfile as tmp

EPSILON_LIST = [1, 5]

LPRCCommand = "prefix-search fullbenchmark -a {} -i {} -p {} -o {} {}"
PSRCCommand = "prefix-search fullbenchmark -a {} -i {} -p {} -o {} {}"

BINARY_NAME = "prefix-search fullbenchmark"
COMMAND_NAME = ""
ALGORITHM_PARAM = "-a {}"
INPUT_PARAM = "-i {}"
PREFIX_PARAM = "-p {}"
OUTPUT_PARAM = "-o {}"
EPSILON_PARAMS = " ".join(['-l {}'.format(epsilon) for epsilon in EPSILON_LIST])

PROJECT_DIR_NAME ="prefix-search"
SCRIPTS_DIR_NAME ="scripts"
DATASETS_DIR_PATH = "resources/dataset/"
PREFIX_FILE_PATH = "resources/prefixes/pref100k.txt"

ALGORITHMS = ['lprc', 'psrc']


if os.getcwd().split('/')[-1] == SCRIPTS_DIR_NAME:  # We are calling the script inside scripts
    os.chdir('..')  # We will move into the root project directory

PROJECT_PATH = os.getcwd()

dataframe_structure = {
    'index': range((len(os.listdir(PROJECT_PATH + '/' + DATASETS_DIR_PATH)) - 1) * 2),
    'columns': [
        'dataset',
        'algorithm',
        'init time',
        'structure size',
        'uncompressed data size',
        'words count',
        'average search time',
    ]
}

epsilon_results = {epsilon: pd.DataFrame(**dataframe_structure) for epsilon in EPSILON_LIST}

def generate_cmd_argument(alg, dataset, output_filename):
    return BINARY_NAME + " " + \
           ALGORITHM_PARAM.format(alg) + " " + \
           INPUT_PARAM.format(dataset) + " " + \
           PREFIX_PARAM.format(PREFIX_FILE_PATH) + " " +\
           OUTPUT_PARAM.format(output_filename) + " " + \
           EPSILON_PARAMS


def get_total_words_found(prefix_result):
    return sum([res['PrefixedWordCount'] for res in prefix_result])

def populate_df_from_dataset(dataset):
    row_idx = 0
    for alg in ALGORITHMS:
        with tmp.NamedTemporaryFile() as f:
            print(generate_cmd_argument(alg, DATASETS_DIR_PATH + dataset, f.name))
            status = subprocess.run([generate_cmd_argument(alg, DATASETS_DIR_PATH + dataset, f.name)], shell=True)
            if status.returncode != 0:  # Ops! Something goes wrong
                print("Ops! Something goes wrong. Exit")
                exit(status.returncode)


            with open(f.name, "r") as res:
                result = json.load(res)

            for row in result:
                epsilon = row['Epsilon']
                init_time = row['InitTime']
                structure_size = sum([size for size in row['StructureSize'].values()])
                uncompressed_size = row['UncompressedDataSize']
                words_count = get_total_words_found(row['PrefixResult'])
                average_search_time = row['TotalSearchTime'] / len(row['PrefixResult'])

                epsilon_df = epsilon_results[epsilon]
                epsilon_df.loc[row_idx] = [
                    dataset,
                    alg,
                    init_time,
                    structure_size,
                    uncompressed_size,
                    words_count,
                    average_search_time
                ]
                row_idx += 1

            print(epsilon_df)


for dataset in os.listdir(PROJECT_PATH + '/' + DATASETS_DIR_PATH):
    if dataset == "w0.txt":
        continue
    populate_df_from_dataset(dataset)

