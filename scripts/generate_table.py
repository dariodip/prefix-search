#!/usr/bin/python3
# This script execute our benchmark on all the datasets inside the directory resources/dataset (excluded w0.txt because
# it is empty). All the result will be saved in the directory resources/results.
# The benchmark will be executed on both the algorithms (lprc and psrc) with 5 different values for the parameter
# epsilon (1, 5, 10, 20 and 100)
import csv
import json
import os
import pandas as pd
import subprocess
import tempfile as tmp

EPSILON_LIST = [1, 5, 10, 20, 100]

BINARY_NAME = "prefix-search fullbenchmark"
ALGORITHM_PARAM = "-a {}"
INPUT_PARAM = "-i {}"
PREFIX_PARAM = "-p {}"
OUTPUT_PARAM = "-o {}"
EPSILON_PARAMS = " ".join(['-l {}'.format(epsilon) for epsilon in EPSILON_LIST])

PROJECT_DIR_NAME = "prefix-search"
SCRIPTS_DIR_NAME = "scripts"
DATASETS_DIR_PATH = "resources/dataset"
RESULTS_DIR = "resources/results"
PREFIX_FILE_PATH = "resources/prefixes/pref10k.txt"
EPSILON_DATASIZE_FILENAME = RESULTS_DIR + '/' + 'eps_{}_datasize.csv'
EPSILON_ALG_FILENAME = RESULTS_DIR + '/' + 'eps_{}_alg_{}.csv'

ALGORITHMS = ['lprc', 'psrc']

if os.getcwd().split('/')[-1] == SCRIPTS_DIR_NAME:  # We are calling the script inside scripts directory
    os.chdir('..')  # We will move into the root project directory

PROJECT_PATH = os.getcwd()

dataframe_structure = {
    'columns': [
        'dataset',
        'algorithm',
        'init time',
        'structure size',
        'uncompressed data size',
        'words count',
        'average search time',
    ],
}


def generate_cmd_argument(alg, dataset, output_filename):
    '''
    Generate a complete prefix-search command.
    :param alg: the algorithm to use
    :param dataset: the dataset containing all the words
    :param output_filename: the file that will contain the results
    '''
    return BINARY_NAME + " " + \
           ALGORITHM_PARAM.format(alg) + " " + \
           INPUT_PARAM.format(dataset) + " " + \
           PREFIX_PARAM.format(PREFIX_FILE_PATH) + " " + \
           OUTPUT_PARAM.format(output_filename) + " " + \
           EPSILON_PARAMS


def get_total_words_found(prefix_result):
    '''
    It returns all words found from a search on different prefixes
    :param prefix_result: list containing results about each prefix search
    :return: the total words found
    '''
    return sum([res['PrefixedWordCount'] for res in prefix_result])


def populate_df_from_dataset(dataset):
    '''
    Execute prefix-search fullbenchmark on the dataset dataset.
    Each epsilon has a dataframe containing the results of the benchmark corresponding to its value
    :param dataset: the dataset to test
    :return: a dict containing a pandas dataframe for each epsilon, containing all the results
    '''
    epsilon_results = {epsilon: pd.DataFrame(**dataframe_structure) for epsilon in EPSILON_LIST}
    for alg in ALGORITHMS:
        with tmp.NamedTemporaryFile() as f:
            print(generate_cmd_argument(alg, DATASETS_DIR_PATH + '/' + dataset, f.name))
            status = subprocess.run([generate_cmd_argument(alg, DATASETS_DIR_PATH + "/" + dataset, f.name)], shell=True)
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

                new_row = pd.DataFrame({
                    'dataset': [dataset],
                    'algorithm': [alg],
                    'init time': [init_time],
                    'structure size': [structure_size],
                    'uncompressed data size': [uncompressed_size],
                    'words count': [words_count],
                    'average search time': [round(average_search_time, 4)],
                })

                epsilon_results[epsilon] = epsilon_results[epsilon].append(new_row, ignore_index=True)

    return epsilon_results


def preparate_csv():
    '''
    Create (or erase the content) of the csv files that will contain all the final results
    It will also add an header to each file
    '''
    for epsilon in EPSILON_LIST:
        filename = EPSILON_DATASIZE_FILENAME.format(epsilon)

        with open(filename, 'w') as csv_file:
            writer = csv.writer(csv_file, delimiter=',')

            writer.writerow([
                'dataset',
                'file dimension',
                'lprc',
                'psrc',
            ])

        for alg in ALGORITHMS:
            alg_filename = EPSILON_ALG_FILENAME.format(epsilon, alg)

            with open(alg_filename, 'w') as alg_csv_file:
                alg_writer = csv.writer(alg_csv_file, delimiter=',')

                alg_writer.writerow([
                    'dataset',
                    'init time',
                    'match found',
                    'average search time',
                ])


def append_to_csv(results: dict):
    '''
    Given a dict of pandas dataframes, containing all the results of the benchmark for each epsilon, append the results
    on the final tables
    :param results: dict containing results for each epsilon
    '''
    for epsilon in results:
        filename = EPSILON_DATASIZE_FILENAME.format(epsilon)
        df = results[epsilon]
        with open(filename, 'a') as csv_file:
            writer = csv.writer(csv_file, delimiter=',')

            for dataset in df['dataset'].unique():
                rows = df[df['dataset'] == dataset]  # We have two row here, one for each algorithm
                csv_row = [
                    dataset,
                    rows['uncompressed data size'][0],
                    rows['structure size'][0],  # lprc structure size
                    rows['structure size'][1],  # psrc structure size
                ]
                writer.writerow(csv_row)

        for alg in ALGORITHMS:
            alg_filename = EPSILON_ALG_FILENAME.format(epsilon, alg)

            with open(alg_filename, 'a') as alg_csv_file:
                writer = csv.writer(alg_csv_file, delimiter=',')
                for dataset in df['dataset'].unique():
                    dataset_rows = df[df['dataset'] == dataset]
                    row = dataset_rows[dataset_rows['algorithm'] == alg]
                    csv_row = [
                        dataset,
                        row['init time'].values[0],
                        row['words count'].values[0],
                        row['average search time'].values[0],
                    ]
                    writer.writerow(csv_row)


if __name__ == '__main__':
    preparate_csv()
    for dataset in os.listdir(PROJECT_PATH + '/' + DATASETS_DIR_PATH):
        if dataset == "w0.txt":
            continue
        append_to_csv(populate_df_from_dataset(dataset))
