# vim: expandtab tabstop=4 shiftwidth=4

from pathlib import Path
from time import sleep

import sys

import minnow

class SumDateProcessor(minnow.Processor):
    def process(self, metadata_file_path, data_file_path, output_path):
        metadata = minnow.load_properties(metadata_file_path)
        value_sum = sum([int(metadata[field]) for field in ['year', 'month', 'day', 'hour', 'minute', 'second']])
        metadata = {'type': 'summed_date', 'sum': value_sum}
        minnow.save_properties(metadata, output_path.joinpath('summed_date.properties'))
        sleep(600)
        output_data_path = output_path.joinpath('summed_date')
        output_data_path.touch()

def main():
    input_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2])
    proc = SumDateProcessor(input_path, output_path)
    proc.run()

if __name__ == "__main__":
    main()
