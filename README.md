# Censored Planet Observatory 
[![Documentation Status](https://readthedocs.org/projects/censoredplanet/badge/?version=latest)](https://censoredplanet.readthedocs.io/en/latest/?badge=latest)
[![analyze-cp Build Status](https://github.com/censoredplanet/censoredplanet/workflows/analyze-cp/badge.svg)](https://github.com/censoredplanet/censoredplanet/actions)


This respository contains documentation about the raw data from the [Censored Planet Observatory](https://censoredplanet.org/data/raw) and includes code to analyze the data and run several useful observations. 

## Analysis (analyze-cp)
 Example analysis tool to parse the raw data files on the [Censored Planet Observatory Website](https://censoredplanet.org/data/raw). Currently, only analysis for `quack-v1`, `hyperquack-v1`, and `satellite-v1` data is supported. `v2` support is coming soon. The analysis tool converts raw data into digestible CSV files that contain aggregates of data at the country level. User is prompted to choose the type of output. 

 `analyze-cp` can be compiled using the makefile in the `analysis` directory or by using `go build` in `analysis/cmd/`. `analysis-cp` needs two REQUIRED inputs (tar.gz file downloaded from [Censored Planet Observatory Website](https://censoredplanet.org/data/raw) and Maxmind GeoLite2-City.mmdb file downloaded from the [Maxmind Website](https://maxmind.com)). 

 `analyze-cp` has the following flags: 
 ```
--input-file, REQUIRED, "Input tar.gz file (downloaded from censoredplanet.org)"
--output-file, Default - output.csv, "Output CSV file"
--mmdb-file, REQUIRED,  "Maxmind Geolocation MMDB file (Download from maxmind.com)"
--satellitev1-html-file, OPTIONAL, Default - "", "JSON file that contains HTML responses for detecting blockpages from satellitev1 resolved IP addresses. The JSON file should have the following fields: 1) ip (resolved ip from satellitev1 that is marked as an anomaly), query (query performed by satellitev1), body (HTML body). If unspecified, the blockpage matching process will be skipped."
--log-file, Default - '-'(STDERR), "file name for logging"
--verbosity, Default - 3, "level of log detail (increasing from 0-5)"
 ```

## Documentation 
The documentation is available in the `docs` directory and it is hosted [here](https://docs.censoredplanet.org). 

Before generating the document, you must run `pip install sphinx`.

To generate the document, run `make html` in the `docs` directory.
The html files will be in the `_build` subdirectory.

## Paper
Take a look at the [Censored Planet CCS paper](https://censoredplanet.org/assets/censoredplanet.pdf) and the rest of our [publications](https://censoredplanet.org/publications) and [reports](https://censoredplanet.org/reports) for in-depth details about how Censored Planet works.

## Citation
Please use the following bibtex to refer to Censored Planet:
```
@inproceedings{sundararaman2020censoredplanet,
title ={Censored Planet: An Internet-Wide, Longitudinal Censorship Observatory},
author ={Sundara Raman, Ram and Shenoy, Prerana and Kohls, Katharina and Ensafi, Roya},
booktitle={In ACM SIGSAC Conference on Computer and Communications Security (CCS)},
year={2020}
}
```