# vim: expandtab tabstop=4 shiftwidth=4

from datetime import datetime
from pathlib import Path
from random import randint
from uuid import uuid4

import sys

def main():
    output_path = Path(sys.argv[1])

    for i in range(10):
        dt = datetime.fromtimestamp(randint(0, 2**32))
        name = str(uuid4())

        metadata_output_path = output_path.joinpath(name+'.properties')
        metadata_output_path.write_text('type = date')

        data_output_path = output_path.joinpath(name)
        data_output_path.write_text(dt.strftime('%Y-%m-%dT%H:%M:%S'))

if __name__ == "__main__":
    main()
