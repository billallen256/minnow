# vim: expandtab tabstop=4 shiftwidth=4

from datetime import datetime
from pathlib import Path

import sys

input_path = Path(sys.argv[1])
output_path = Path(sys.argv[2])
metadata_path = list(input_path.glob('*.properties'))[0]
data_path = metadata_path.with_suffix('')

date_str = data_path.read_text().strip()
parsed_date = datetime.strptime(date_str, '%Y-%m-%dT%H:%M:%S')
metadata = {'type':'parsed_date',
            'year': parsed_date.year,
            'month': parsed_date.month,
            'day': parsed_date.day,
            'hour': parsed_date.hour,
            'minute': parsed_date.minute,
            'second': parsed_date.second}

output_metadata_path = output_path.joinpath('parsed_date.properties')
metadata_str = '\n'.join(['{} = {}'.format(k, v) for k, v in metadata.items()])
output_metadata_path.write_text(metadata_str)
output_data_path = output_path.joinpath('parsed_date')
output_data_path.touch()
