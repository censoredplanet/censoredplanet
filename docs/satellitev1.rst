############
Satellite v1
############
Satellite is Censored Planet's tool to detect DNS interference. Refer to the following papers for more details:

* `Global Measurement of DNS Manipulation <https://censoredplanet.org/assets/Pearce2017b.pdf>`_
* `Satellite: Joint Analysis of CDNs and Network-Level Interference <https://censoredplanet.org/assets/Scott2016a.pdf>`_

Satellite v1 corresponds to measurements from 2018 to December 2020. See Satellite for documentation on measurements taken after this date.

The published data has the following directory structure: ::

    CP_Satellite-YYYY-MM-DD-HH-MM-SS/
    |-- log.json
    |-- raw/
    |   |-- answers_control.json
    |   |-- answers_err.json
    |   |-- answers_ip.json
    |   |-- answers.json
    |   |-- answers_raw.json
    |   |-- dns.pkt
    |   |-- interference_err.json
    |   |-- interference.json
    |   |-- resolvers_err.json
    |   |-- resolvers_ip.json
    |   |-- resolvers.json
    |   |-- resolvers_ptr.json
    |   |-- resolvers_raw.json
    |   |-- tagged_answers.json
    |   |-- tagged_resolvers.json
    |-- stat/
        |-- stat_answers.json
        |-- stat_interference_agg.json
        |-- stat_interference_count.json
        |-- stat_interference_country_domain.json
        |-- stat_interference_country.json
        |-- stat_interference_country_percentage.json
        |-- stat_interference_err.json
        |-- stat_interference.json
        |-- stat_resolvers_country.json
        |-- stat_resolvers.json
        |-- stat_tagged.json


*******
Output
*******

:code:`resolvers_err.json`

:code:`resolvers_ip.json`

:code:`resolvers.json`

:code:`resolvers_ptr.json`

:code:`resolvers_raw.json`

:code:`answers_control.json`

* :code:`resolver` : String
    The IP address of the vantage point (a DNS resolver).
* :code:`query` : String
    The domain being queried.
* :code:`answers` : Array
    Contains the resolver's returned answer IPs for the queried domain.

:code:`answers_err.json`

* :code:`resolver` : String
    The IP address of the vantage point (a DNS resolver).
* :code:`query` : String
    The domain being queried.
* :code:`error` :
    Either "no_answer" or a dictionary with additional error information.

:code:`answers_ip.json`

* :code:`answer`: String
	An IP address from a query response.

:code:`answers.json`

* :code:`resolver` : String
    The IP address of the vantage point (a DNS resolver).
* :code:`query` : String
    The domain being queried.
* :code:`answers` : Array
    Contains the resolver's returned answer IPs for the queried domain.

:code:`answers_raw.json`

* :code:`resolver` : String
    The IP address of the vantage point (a DNS resolver).
* :code:`query` : String
    The domain being queried.
* :code:`data` : String
    Raw query response.

:code:`tagged_answers.json`

* :code:`ip` : String
	An IP address from a query response.
* :code:`http` : String
	The hash of the HTTP body.
* :code:`cert` : String
	The hash of the TLS certificate.
* :code:`asname` : String
	The autonomous system (AS) name.
* :code:`asnum` : Integer
	The autonomous system (AS) number.

:code:`tagged_resolvers.json`

* :code:`resolver` : String
	The IP address of the vantage point (a DNS resolver).
* :code:`country` : String
	The full name of the country where the resolver is located.

:code:`interference_err.json` contains resolver answers for queries with no control response, with the following fields:

* :code:`resolver` : String
    The IP address of the vantage point (a DNS resolver).
* :code:`query` : String
    The domain being queried.
* :code:`answers` : Array
    Contains the resolver's returned answer IPs for the queried domain.

:code:`interference.json` contains the interference assessment for the remaining resolver answers, with the following fields:

* :code:`resolver` : String
    The IP address of the vantage point (a DNS resolver).
* :code:`query` : String
    The domain being queried.
* :code:`answers` : JSON object
    The resolver's returned answer IPs for the queried domain are the keys. Each answer IP is mapped to an array of its tags that matched the control tags - if the IP is in the control set, "ip" is appended and if the IP has no tags, "no_tags" is appended.
* :code:`passed` : Boolean
    Equals true if interference is not detected.