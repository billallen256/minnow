# vim: expandtab tabstop=4 shiftwidth=4

from pathlib import Path
from time import sleep

import sys

input_path = Path(sys.argv[1])
output_path = Path(sys.argv[2])
metadata_path = list(input_path.glob('*.properties'))[0]
data_path = metadata_path.with_suffix('')

metadata = metadata_path.read_text().split('\n')
metadata = [tuple(x.split('=')) for x in metadata if len(x.strip()) > 0]
metadata = {x[0].strip():x[1].strip() for x in metadata}
value_sum = sum([int(metadata[field]) for field in ['year', 'month', 'day', 'hour', 'minute', 'second']])

output_metadata_path = output_path.joinpath('summed_date.properties')
metadata_str = 'type = summed_date\nsum = {}'.format(value_sum)
sleep(600)
output_metadata_path.write_text(metadata_str)
output_data_path = output_path.joinpath('summed_date')
output_data_path.touch()
