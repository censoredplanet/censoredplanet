# Censored Planet Observatory 
[![Documentation Status](https://readthedocs.org/projects/censoredplanet/badge/?version=latest)](https://censoredplanet.readthedocs.io/en/latest/?badge=latest)
[![analyze-cp Build Status](https://github.com/censoredplanet/censoredplanet/workflows/analyze-cp/badge.svg)](https://github.com/censoredplanet/censoredplanet/actions)


This respository contains documentation about the raw data from the [Censored Planet Observatory](https://censoredplanet.org/data/raw) and includes code to analyze the data and run several useful observations. 

## Analysis (analyze-cp)
 Example analysis tool to parse the raw data files on the [Censored Planet Observatory Website](https://censoredplanet.org/data/raw). Currently, only analysis for `quack-v1` and `hyperquack-v1` data is supported. `satellite-v1` support is coming soon. The analysis tool converts raw data into digestible CSV files that contain aggregates of data at the country level. User is prompted to choose the type of output. 

 `analyze-cp` can be compiled using the makefile in the `analysis` directory or by using `go build` in `analysis/cmd/`. `analysis-cp` needs two REQUIRED inputs (tar.gz file downloaded from [Censored Planet Observatory Website](https://censoredplanet.org/data/raw) and Maxmind GeoLite2-City.mmdb file downloaded from the [Maxmind Website](https://maxmind.com)). 

 `analyze-cp` has the following flags: 
 ```
--input-file, REQUIRED, "Input tar.gz file (downloaded from censoredplanet.org)"
--output-file, Default - output.csv, "Output CSV file"
--mmdb-file, REQUIRED,  "Maxmind Geolocation MMDB file (Download from maxmind.com)"
--log-file, Default - '-'(STDERR), "file name for logging"
--verbosity, Default - 3, "level of log detail (increasing from 0-5)"
 ```


## Documentation 
The documentation is available in the `docs` directory and it is hosted on readthedocs at `https://censoredplanet.readthedocs.io`. 

Before generating the document, you must run `pip install sphinx`.

To generate the document, run `make html` in the `docs` directory.
The html files will be in the `_build` subdirectory.
