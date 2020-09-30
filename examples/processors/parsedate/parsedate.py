# vim: expandtab tabstop=4 shiftwidth=4

from datetime import datetime
from pathlib import Path

import sys

import minnow

class ParseDateProcessor(minnow.Processor):
    def process(self, metadata_file_path, data_file_path, output_path):
        date_str = data_file_path.read_text().strip()
        parsed_date = datetime.strptime(date_str, '%Y-%m-%dT%H:%M:%S')
        metadata = {'type':'parsed_date',
                    'year': parsed_date.year,
                    'month': parsed_date.month,
                    'day': parsed_date.day,
                    'hour': parsed_date.hour,
                    'minute': parsed_date.minute,
                    'second': parsed_date.second}

        minnow.save_properties(metadata, output_path.joinpath('parsed_date.properties'))
        output_data_path = output_path.joinpath('parsed_date')
        output_data_path.touch()

def main():
    input_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2])
    proc = ParseDateProcessor(input_path, output_path)
    proc.run()

if __name__ == "__main__":
    main()
