############
Satellite v2
############
Satellite is Censored Planet's tool to detect DNS interference. Refer to the following papers for more details:

* `Global Measurement of DNS Manipulation <https://censoredplanet.org/assets/Pearce2017b.pdf>`_
* `Satellite: Joint Analysis of CDNs and Network-Level Interference <https://censoredplanet.org/assets/Scott2016a.pdf>`_

The published data has the following directory structure: ::

    CP_Satellite-YYYY-MM-DD-HH-MM-SS/
    |-- log.json
    |-- raw/
        |-- blockpages.json
        |-- dns.pkt
        |-- resolvers_err.json
        |-- resolvers_ip.json
        |-- resolvers.json
        |-- resolvers_ptr.json
        |-- resolvers_raw.json
        |-- responses_control.json
        |-- responses_ip.json
        |-- responses.json
        |-- responses_raw.json
        |-- results.json
        |-- tagged_responses.json
        |-- tagged_resolvers.json


*******
Output
*******

The relevant output is located in the `raw/` directory.

------
Probe
------

1. Generate a DNS A query packet for a controlled domain (`dns.pkt`).

2. Perform a ZMap scan with the probe packet for open DNS resolvers.

    :code:`resolvers_raw.json` contains the ZMap output:

    * :code:`saddr` : String
        IP address of a DNS resolver.
    * :code:`data` : String
        Raw response to probe domain.

------
Filter
------

1. Check the probe responses of the resolvers found by ZMap.

    :code:`resolvers_ip.json` contains resolvers that returned the correct probe response:

    * :code:`vp` : String
        The IP address of the potential vantage point (a DNS resolver).
    * :code:`response` : String
        The resolver's response (IP address) to the probe domain.

2. Perform PTR queries on the IPs of resolvers with the correct probe response.

    :code:`resolvers_err.json` contains resolvers with failed PTR queries:

    * :code:`vp` : String
        The IP address of the potential vantage point (a DNS resolver).
    * :code:`error` : JSON Object
        Contains error information.

    :code:`resolvers_ptr.json` contains resolvers with succesful PTR queries:

    * :code:`vp` : String
        The IP address of the potential vantage point (a DNS resolver).
    * :code:`names` : Array
        Result from PTR query (the hostname).

3. Identify infrastructure resolvers from successful PTR queries and add predefined "control" and "special" resolvers to form final set of vantage points.

    :code:`resolvers.json` contains the infrastructure, "control", and "special" resolvers.

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`name` : String
        Result from PTR query (if infrastructure), "control", or "special".

------
Query
------

1. Make DNS queries for each test domain to each resolver.

    :code:`responses_raw.json` contains raw responses from successful queries:

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`test_url` : String
        The test domain being queried.
    * :code:`data` : String
        Raw query response.

    **Note:**

        * The query for the test domain is attempted up to four times in case of non Type A response. To check the status of the resolver, a control domain is queried before and after the queries for the test domain.

2. Parse and separate responses from control resolvers and non-control resolvers.

    :code:`responses_control.json` contains responses for queries to control resolvers and :code:`responses.json` contains responses for queries to non-control resolvers:

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`test_url` : String
        The test domain being queried.
    * :code:`response` : Array
        The resolver's responses for the control and test domain - in the order control domain, test domain (up to 4 attempts), control domain.

        * :code:`url` : String
            The domain being queried in this trial (either the control domain or :code:`test_url`)
        * :code:`has_type_a` : Boolean
            Equals true if the query returned a valid A resource record.
        * :code:`answer` : Array
            The resolver's response for the queried domain in this trial (list of answer IPs if successful).
        * :code:`error` : String
            Contains error information.
        * :code:`rcode` : Integer
            Response code mapping to success (0) or errors (>0).
        * :code:`start_time` : String
            The start time of the measurement.
        * :code:`end_time` : String
            The end time of the measurement.
    * :code:`resolver_status` : Boolean
        Equals true if the resolver succesfully responds to the two control queries.
    * :code:`raw` : Array
        The resolver's unparsed responses (corresponding to the respective index in :code:`response`).

3. Determine set of IP addresses that appeared across all query responses for tagging.

    :code:`responses_ip.json` contains these IPs, one IP per line:

    * :code:`response` : String
        An IP address from a query response.

------
Tag
------

1. Tag each answer IP with information from Censys.

    :code:`tagged_responses.json` contains the answer IPs and their HTTP, TLS, and AS tags:

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

    **Note:**

        * Fields may have null values if the information was not available on Censys.

2. Tag each resolver with the location from Maxmind.

    :code:`tagged_resolvers.json` contains the resolvers and their countries:

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`location`: JSON object
        * :code:`country_name` : String
            The full name of the country where the resolver is located.
        * :code:`country_code` : String
            The two-letter ISO 3166 code of the country where the resolver is located.

------
Detect
------

1. Compare query responses between non-control resolvers and control resolvers to identify interference.

    :code:`results.json` contains the interference assessment for the query responses:

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`location`: JSON object
        * :code:`country_name` : String
            The full name of the country where the resolver is located.
        * :code:`country_code` : String
            The two-letter ISO 3166 code of the country where the resolver is located.
    * :code:`test_url` : String
        The domain being queried.
    * :code:`response` : JSON object
        The resolver's returned answer IPs for the queried domain are the keys. Each answer IP is mapped to an array of its tags that matched the control tags - if the IP is in the control set, "ip" is appended and if the IP has no tags, "no_tags" is appended. Also has an :code:`rcode` field mapping to a list of response codes for the trials.
    * :code:`passed_control` : Boolean
        Equals true if both control queries were successful.
    * :code:`in_control_group` : Boolean
        Equals true if at least one control resolver had a valid response for this test domain.
    * :code:`connect_error` : Boolean
        Equals true if all test domain query attempts returned errors.
    * :code:`anomaly` : Boolean
        Equals true if an anomaly is detected.
    * :code:`start_time` : String
        The start time of the measurement.
    * :code:`end_time` : String
        The end time of the measurement.
    * :code:`confidence` : JSON object
        * :code:`average` : Float
            Average percentage of tags matching the control set for the answers (average of :code:`matches`).
        * :code:`matches` : Array
            Contains the percentage of tags matching the control set for each answer. If an answer IP is in the control set, the percentage for that answer is 100 even if the IP has no tags.
        * :code:`untagged_controls` : Boolean
            Equals true if all control IPs for the query have no tags.
        * :code:`untagged_answers` : Boolean
            Equals true if all answer IPs have no tags.

    **Note:**

        * For each response, the answer IPs and their tags are compared to the set of answer IPs and tags from all the control resolvers for the same query domain. A response is classified as an anomaly if there is no overlap between the two.

        * Cases where the control answer IPs have no tags will be considered anomalies if the resolver's answer IPs are not in the control set.

------
Fetch
------

1. Perform HTTP(S) GET requests to the IPs identified as anomalies.

    :code:`blockpages.json` contains the responses:

    * :code:`ip` : String
        The IP address from an anomalous DNS response.
    * :code:`keyword` : String
        The domain queried for the anomalous DNS response.
    * :code:`http` : Object
        HTTP response.
    * :code:`https` : Object
        HTTPS response.
    * :code:`fetched` : Boolean
        Equals true if page is successfully fetched.
    * :code:`start_time` : String
        The start time of the measurement.
    * :code:`end_time` : String
        The end time of the measurement.


*******
Modules
*******

This is a brief tour of the modules in satellite.

All files mentioned are under the :code:`rawDir` designated in :code:`config.go`, unless specified.

* :code:`probe`:  probes IPv4 address space for resolvers using zmap.
    * input: N/A
    * output:
        * list of resolver candidates (:code:`resolvers_raw.json`)
* :code:`filter`: removes resolvers that aren't infrastructure (runs PTR queries on resolvers).
    * input:
        * :code:`resolvers_raw.json`
    * output:
        * list of filtered public open resolvers (:code:`resolvers.json`)
        * (:code:`resolvers_ip.json`)
        * list of PTR query results (:code:`resolvers_PTR.json`)
        * list of erroneous PTR query results (:code:`resolvers_err.json`)
* :code:`query`:  queries public open resolvers with a list of domains.
    * input:
        * list of resolvers to query (:code:`resolvers.json`)
        * list of domains for querying (:code:`assets/input_lists/test_domains`)
        * control resolvers (:code:`assets/satellite/control_resolvers.txt`)
        * special resolvers (:code:`assets/satellite/special_resolvers.txt`)
    * output:
        * answers from control resolvers (:code:`responses_control.json`)
        * IP list of answers (:code:`responses_ip.json`)
        * non control resolver answers (:code:`responses.json`)
        * raw response packets (:code:`responses_raw.json`)
* :code:`tag`:  tags resolvers with MaxMind (country) and IPs with censys (certificate, AS number and AS name).
    * input:
        * list of resolvers to query (:code:`resolvers.json`)
        * list of answered IPs (:code:`responses_ip.json`)
    * output:
        * list of tagged IPs with Censys (:code:`tagged_responses.json`)
        * list of tagged resolvers with Maxmind (:code:`tagged_resolvers.json`)
* :code:`detect`: detects interference by comparing DNS query responses to control set.
    * input:
        * :code:`tagged_responses.json`
        * answers from non control resolvers (:code:`responses.json`)
        * answers from control resolvers (:code:`responses_control.json`)
        * :code:`assets/satellite/control_resolvers.txt`
    * output:
        * list of interference result (:code:`results.json`)
* :code:`fetch`: fetches pages hosted on the IPs identified as interference for future blockpage analysis.
    * input:
        * :code:`results.json`
    * output:
        * list of tampered IPs, and results of HTTP(S) GET (:code:`blockpages.json`)
* :code:`stat`:   data analysis.
* :code:`full`:   all aforementioned modules combined.
* :code:`upload`