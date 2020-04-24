# vim: expandtab tabstop=4 shiftwidth=4

from pathlib import Path

import sys

input_path = Path(sys.argv[1])
output_path = Path(sys.argv[2])
metadata_path = list(input_path.glob('*.properties'))[0]
data_path = metadata_path.with_suffix('')

metadata = metadata_path.read_text().split('\n')
metadata = [tuple(x.split('=')) for x in metadata if len(x.strip()) > 0]
metadata = {x[0].strip():x[1].strip() for x in metadata}

for field in ['year', 'month', 'day', 'hour', 'minute', 'second']:
    metadata[field] = int(metadata[field])

output_metadata_path = output_path.joinpath('formatted_date.properties')
output_metadata_path.write_text('type = date')
output_data_path = output_path.joinpath('formatted_date')
output_data_path.write_text('{year:04d}-{month:02d}-{day:02d}T{hour:02d}:{minute:02d}:{second:02d}'.format(**metadata))
