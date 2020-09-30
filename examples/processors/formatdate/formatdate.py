# vim: expandtab tabstop=4 shiftwidth=4

from pathlib import Path

import sys

import minnow

class FormatDateProcessor(minnow.Processor):
    def process(self, metadata_file_path, data_file_path, output_path):
        metadata = minnow.load_metadata(metadata_file_path)

        for field in ['year', 'month', 'day', 'hour', 'minute', 'second']:
            metadata[field] = int(metadata[field])

        output_metadata_path = output_path.joinpath('formatted_date.properties')
        output_metadata_path.write_text('type = date')
        output_data_path = output_path.joinpath('formatted_date')
        output_data_path.write_text('{year:04d}-{month:02d}-{day:02d}T{hour:02d}:{minute:02d}:{second:02d}'.format(**metadata))

def main():
    input_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2])
    proc = FormatDateProcessor(input_path, output_path)
    proc.run()

if __name__ == "__main__":
    main()
