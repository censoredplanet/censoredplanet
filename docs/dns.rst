####################
DNS Data - Satellite
####################

`Satellite/Iris <https://censoredplanet.org/projects/satellite>`_ is Censored Planet’s remote measurement technique that detects DNS interference using Open DNS resolvers. Below, we provide an overview of Satellite and its data format. Refer to our `academic papers <https://censoredplanet.org/assets/Pearce2017b.pdf>`_ for in-depth details about Satellite.

************
Satellite-v1
************

.. image:: images/satellite-v1.png
  :width: 600
  :alt: Figure - Overview of Satellite-v1

Figure - Overview of Satellite-v1

Satellite-v1 is the first version of Satellite that we operated from August 2018 - February 2021. The primary function of Satellite is to detect incorrect DNS resolutions from open DNS resolvers in many countries.

* From a measurement machine at the University of Michigan, we send a DNS query for a website whose reachability we’re interested in, to an open DNS resolver in a country of interest (1). The response from the DNS resolver is our Test IP (2).

* We also send a DNS query for the same website to trusted control resolvers (3), and record their response as the control IP (4).

* We then compare the test and control responses using several heuristics, including a direct IP address comparison, and comparison of the AS number, AS names, HTTP content hashes, and TLS certificates associated with the test and control IP addresses (5). Satellite-v1 only labels a measurement as an anomaly when all of the heuristics mismatch.

Our various `publications <http://censoredplanet.org/publications>`_ and `reports <http://censoredplanet.org/reports>`_ have used Satellite-v1 to detect many cases of DNS manipulation. For instance, in our `recent investigation into the filtering of COVID-19 websites <https://censoredplanet.org/assets/covid.pdf>`_, Satellite-v1 found many networks using website filtering products to manipulate DNS responses of COVID-related websites.

==============
Satellite-v1.0
==============


Data Format
***********


The published data has the following directory structure: ::

    CP_Satellite-YYYY-MM-DD-HH-MM-SS/
    |-- answers_control.json
    |-- answers_err.json
    |-- answers_ip.json
    |-- answers.json
    |-- answers_raw.json
    |-- dns.pkt
    |-- interference_err.json
    |-- interference.json
    |-- resolvers_err.json
    |-- resolvers_ip.json
    |-- resolvers.json
    |-- resolvers_ptr.json
    |-- resolvers_raw.json
    |-- tagged_answers.json
    |-- tagged_resolvers.json


Output
======

The relevant output is located in the `raw/` directory.


Probe
~~~~~~~~~~~

1. Generate a DNS A query packet for a controlled domain (`dns.pkt`).

2. Perform a ZMap scan with the probe packet for open DNS resolvers.

	:code:`resolvers_raw.json` contains the ZMap output:

	* :code:`saddr` : String
		IP address of a DNS resolver.
	* :code:`data` : String
		Raw response to probe domain.


Filter
~~~~~~~~~~~

1. Check the probe responses of the resolvers found by ZMap. 

	:code:`resolvers_ip.json` contains resolvers that returned the correct probe response:

	* :code:`resolver` : String
	    IP address of a DNS resolver.
	* :code:`answer` : String
		The resolver's response (IP address) to the probe domain.

2. Perform PTR queries on the IPs of resolvers with the correct probe response.

	:code:`resolvers_err.json` contains resolvers with failed PTR queries:

	* :code:`resolver` : String
	    IP address of a DNS resolver.
	* :code:`error` : JSON Object
		Contains error information.

	:code:`resolvers_ptr.json` contains resolvers with succesful PTR queries:

	* :code:`resolver` : String
	    IP address of a DNS resolver.
	* :code:`names` : Array
	    Result from PTR query (the hostname).

3. Identify infrastructure resolvers from successful PTR queries and add predefined "control" and "special" resolvers to form final set of vantage points.

	:code:`resolvers.json` contains the infrastructure, "control", and "special" resolvers.

	* :code:`resolver` : String
	    IP address of a DNS resolver.
	* :code:`name` : String
		Result from PTR query (if infrastructure), "control", or "special".


Query
~~~~~~~~~~~

1. Make DNS queries for each test domain to each resolver.

	:code:`answers_err.json` contains erroneous queries:

	* :code:`resolver` : String
	    The IP address of the vantage point (a DNS resolver).
	* :code:`query` : String
	    The domain being queried.
	* :code:`error` : String / JSON Object
	    Either "no_answer" or a dictionary with additional error information.

	**Note:**

		* In some cases, the :code:`resolver` field may be replaced by :code:`ip` - both are referring to the resolver's IP.

		* "no_answer" appears in the :code:`error` field if no A resource records (IPs) are returned - this includes the :code:`NXDOMAIN` response.

		* Responses with :code:`NXDOMAIN` or other errors may indicate censorship. However, these cases are not analyzed further in Satellite-v1. 

	:code:`answers_raw.json` contains raw responses from successful queries:

	* :code:`resolver` : String
	    The IP address of the vantage point (a DNS resolver).
	* :code:`query` : String
	    The domain being queried.
	* :code:`data` : String
	    Raw query response.

2. Separate responses (converted to IP addresses) from control resolvers and non-control resolvers.

	:code:`answers_control.json` contains responses for queries to control resolvers:

	* :code:`resolver` : String
	    The IP address of the vantage point (a DNS resolver).
	* :code:`query` : String
	    The domain being queried.
	* :code:`answers` : Array
	    The resolver's response for the queried domain (list of answer IPs).

	:code:`answers.json` contains responses for queries to non-control resolvers:

	* :code:`resolver` : String
	    The IP address of the vantage point (a DNS resolver).
	* :code:`query` : String
	    The domain being queried.
	* :code:`answers` : Array
	    The resolver's response for the queried domain (list of answer IPs).

3. Determine set of IP addresses that appeared across all query responses for tagging.

	:code:`answers_ip.json` contains these IPs, one IP per line:

	* :code:`answer` : String
		An IP address from a query response.


Tag
~~~~~~~~~~~

1. Tag each answer IP with information from Censys.

	:code:`tagged_answers.json` contains the answer IPs and their HTTP, TLS, and AS tags: 

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

	* :code:`resolver` : String
		The IP address of the vantage point (a DNS resolver).
	* :code:`country` : String
		The full name of the country where the resolver is located.


Detect
~~~~~~~~~~~

1. Compare query responses between non-control resolvers and control resolvers to identify interference.

	:code:`interference_err.json` contains resolver responses for queries with no control response:

	* :code:`resolver` : String
	    The IP address of the vantage point (a DNS resolver).
	* :code:`query` : String
	    The domain being queried.
	* :code:`answers` : Array
	    The resolver's response for the queried domain (list of answer IPs).

	:code:`interference.json` contains the interference assessment for the remaining resolver responses:

	* :code:`resolver` : String
	    The IP address of the vantage point (a DNS resolver).
	* :code:`query` : String
	    The domain being queried.
	* :code:`answers` : JSON object
	    The resolver's returned answer IPs for the queried domain are the keys. Each answer IP is mapped to an array of its tags that matched the control tags - if the IP is in the control set, "ip" is appended and if the IP has no tags, "no_tags" is appended.
	* :code:`passed` : Boolean
	    Equals true if interference is not detected. Note that if this field is set to false, it may indicate either DNS interference, or an unexpected answer for the resolution. Further manual confirmation is required to confirm censorship.

	**Note:**

		* For each response, the answer IPs and their tags are compared to the set of answer IPs and tags from all the control resolvers for the same query domain. A response is classified as interference if there is no overlap between the two. 

		* Cases where the control answer IPs have no tags will be considered interference if the resolver's answer IPs are not in the control set.

		* Satellite-v1 anomalies (interference detected) need to be explicitly confirmed by fetching pages hosted at the resolved IPs in post-processing. This functionality is included by default in Satellite v2.


Limitations
***********

Although Satellite-v1 was extremely useful in detecting DNS interference at large scale, it suffered from several limitations, which form the improvements in Satellite-v2.x.

* Satellite-v1 could not detect DNS censorship where A records were not available i.e. Satellite-v1 primarily focused on detecting incorrect DNS resolutions through the resolved IP address, and did not contain heuristics to measure DNS manipulation which manifested through timeouts, NXDOMAIN responses, SERVFAIL responses, etc.

* Satellite-v1 required post-processing to remove false positives and confirm the presence of anomalies, such as through using post-measurement heuristics and blockpage regexes. Satellite-v2 has the inbuilt capability to perform most post-processing measurements.

* Satellite-v1’s data format provided results in several files, which were hard to parse. With Satellite-v2, our aim is also to present an easier, more intuitive data format.


**************
Satellite-v2
**************

.. image:: images/satellite-v2.png
  :width: 600
  :alt: Figure - Overview of Satellite-v2.0

Figure - Overview of Satellite-v2

Satellite-v2 is our brand new version of Satellite, where we’ve made several modifications to the measurement technique and data format for facilitating accurate and efficient remote DNS interference measurements. Below, we detail the major changes we’ve made in Satellite-v2.

* **Measuring DNS interference without A records** - In Satellite-v2, we have added a sandwiched retry mechanism to our Satellite measurements in order to detect DNS interference that results in a non-zero R code response. A description of the method is shown in the figure below. We first make a control query to the open DNS resolver, providing a domain name that we do not expect to be blocked (eg. www.example.com). After the control query, we make up to 4 retries of the test DNS query, providing the test domain name. In case an A record is detected, we stop the test measurement. At the end, we perform another control query similar to the first measurement. The control queries ensure that the resolver is behaving correctly for an innocuous domain, and the multiple retry mechanism accounts for temporary errors in the network. With the help of the sandwiched retry mechanism, Satellite-v2 is able to detect DNS interference that manifests as timeouts, NXDOMAIN, SERVFAIL etc. From our preliminary analysis of Satellite-v2 data, we’ve already found several cases of DNS interference that can be identified using this method. For example, from the Satellite-v2 scan performed on 2021-03-17, we are able to identify 174,795 responses that have non-zero R codes from China, which makes up 15.6% out of the responses marked as interference. This kind of DNS interference was previously omitted by satellite v1. Shown below is an example measurement that passed the sandwich control tests, but received server failure R code. This could be an indicator of censorship or geoblocking.

* **Fetching HTML pages hosted at resolved IPs marked as an anomaly** -  Satellite-v2 has an in-built fetch feature that performs HTTP and HTTPS GET requests to resolved IPs that fail our heuristics, and we store the HTML responses in blockpages.json. Satellite-v2 data files available on our website contain this file for easier confirmation of DNS censorship, while this step was being performed as a post-processing step in Satellite-v1. This addition helps in quickly identifying blockpages such as the example shown in the figure below.

* **Adding scan-level heuristics to exclude false positives** - Another step part of the post-processing pipeline of Satellite-v1 that is inbuilt in Satellite-v2. We exclude potentially false positive anomalies by using scan-level heuristics, such as the number of domains resolving to the anomalous IP address, or the anomalous IP address being part of a big CDN. Note that this step may lead to Satellite-v2 missing certain censorship. This output can be found in results_verified.json.

* **Other changes** - We updated the heuristics to determine whether a DNS response is interfered - Satellite-v2 now includes a new “confidence” field, which addresses the certainty of interference according to the state of comparison between responses from the test resolvers and the control resolvers. We also make sure that IPs with no metadata information from Censys are not marked as interference.

  We have also reorganized our output files so that they are easier to read. The primary output files containing DNS interference data are results.json and results_verified.json. Satellite-v2 integrates more information in the results.json file, like the country name and country code of the target resolver, and start time and end time of each measurement. We hope this modification makes processing of the satellite data easier for our users.

==============
Satellite-v2.0
==============

Data Format
***********

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
        |-- results_verified.json
        |-- tagged_responses.json
        |-- tagged_resolvers.json


Probe
~~~~~

1. Generate a DNS A query packet for a controlled domain (`dns.pkt`).

2. Perform a `ZMap <https://github.com/zmap/zmap>`_ (Internet-wide) scan with the probe packet for open DNS resolvers.

    :code:`resolvers_raw.json` contains the ZMap output:

    * :code:`saddr` : String
        IP address of a DNS resolver.
    * :code:`data` : String
        Raw response to probe domain.


Filter
~~~~~~

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


Query
~~~~~

1. Make DNS queries for each test domain to each resolver.

    :code:`responses_raw.json` contains raw responses from successful queries:

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`test_url` : String
        The test domain being queried.
    * :code:`data` : String
        Raw query response.

    **Note:**

        * NEW: The query for the test domain is attempted up to four times in case of non Type A response. To check the status of the resolver, a control domain is queried before and after the queries for the test domain.

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


Tag
~~~

1. Tag each answer IP with information from `Censys <https://about.censys.io/>`_.

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


Detect
~~~~~~

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
        Equals true if an anomaly is detected. In case there are no tags for the answers or control, then this field is conservatively marked as false. 
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


Fetch
~~~~~

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


Verify
~~~~~~

1. New hueristics to exclude possible cases of erroneous answers from resolvers. Currently, verify excludes answer IPs that are part of big CDNs (Note: this could lead to false negatives) and answer IPs that appear for a low number of domains (<=2). 
    :code:`results_verified.json` contains only the rows that were earlier marked as anomalies:

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
    * :code:`excluded` : Boolean
        Should this observation be excluded from being counted as an anomaly?
    * :code:`exclude_reason` : String Array
        If observation should be excluded, why? (eg. "is_CDN")


==============
Satellite-v2.1
==============
Satellite-v2.1 incorporates minor changes from Satellite-v2.0, starting after April 14, 2021. The changes include,
* In the filter module, we removed resolvers from resolvers.json if they can not resolve the root server
* We removed the liveness test response from results.json and results_verified.json.
* We removed the error messages from results.json and results_verified.json, if any.

==============
Satellite-v2.2
==============
Satellite-v2.2 incorporates major changes in code and data structure from Satellite-v2.1, but no major changes in the functionality of Satellite. The changes are made after June 7, 2021 and they include,
* Store information generated from the query, tag, detect, and verify module in memory, producing only one file (results.json) as output, instead of generating outputs for every module. Renamed query-tag-detect-verify as “test” module, and probe-filter as “discovery”.
* Updated test module so that it first conducts queries for control resolvers, and then query, tag and detect test resolvers in batches.

Data Format
***********

The published data has the following directory structure: ::

    CP_Satellite-YYYY-MM-DD-HH-MM-SS/
    |-- resolvers_raw.json
    |-- dns.pkt
    |-- resolvers.json
    |-- results_verified.json
    |-- blockpages.json

Satellite v2 is divided into three parts: 

1. :code:`discovery`: consist of :code:`probe` and :code:`filter` modules.

2. :code:`test`: consist of :code:`query`, :code:`tag` and :code:`detect` modules.

3. verification and blockpage fetching: consist of :code:`fetch` and :code:`verify`.


Probe
~~~~~

1. Generate a DNS A query packet for a controlled domain (:code:`dns.pkt`).

2. Perform a `ZMap <https://github.com/zmap/zmap>`_ (Internet-wide) scan with the probe packet for open DNS resolvers.

    :code:`resolvers_raw.json` contains the ZMap output:

    * :code:`saddr` : String
        IP address of a DNS resolver.
    * :code:`data` : String
        Raw response to probe domain.


Filter
~~~~~~

1. Perform PTR queries on the IPs of resolvers found by ZMap and filter out the ones without PTR records.

2. Perform Liveness test on the infrastructural resolvers and filter out the ones that fail.

3. Add predefined "control" and "special" resolvers to form the final set of vantage points.

4. Tag each resolver with the location from Maxmind.

    :code:`resolvers.json` contains the infrastructure, "control", and "special" resolvers.

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`name` : String
        Result from PTR query (if infrastructure), "control", or "special".
    * :code:`location`: JSON object
        * :code:`country_name` : String
            The full name of the country where the resolver is located.
        * :code:`country_code` : String
            The two-letter ISO 3166 code of the country where the resolver is located.


Query
~~~~~

1. Make DNS queries for each test domain to each resolver. The query for the test domain is attempted up to four times in case of connection error. To check the status of the resolver, a control measurement is conducted before the queries for the test domain. If the first control measurement fails, no further measurements will be conducted for the same :code:`(resolver, domain)` pair. If all 4 trials for the test domain fail, another control measurement will be conducted.

2. Parse and separate responses from control resolvers and non-control resolvers.


Tag
~~~

1. Tag each answer IP with information from `Censys <https://about.censys.io/>`_.
    **Note:**

        * Fields may have empty strings if the information was not available on Censys.


Detect
~~~~~~

1. Compare query responses between non-control resolvers and control resolvers to identify interference. When running satellite v2 as a whole module, :code:`detect` does not output any files. However, when run separately, :code:`detect` outputs :code:`results.json` with the :code:`excluded` field set to :code:`false` and the :code:`excluded_reason` field set to :code:`null` by default. (See the output structure in :code:`verify` section)

    **Note:**

        * For each response, the answer IPs and their tags are compared to the set of answer IPs and tags from all the control resolvers for the same query domain. A response is classified as an anomaly if there is no overlap between the two.


Fetch
~~~~~

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
        Equals true if a page is successfully fetched.
    * :code:`start_time` : String
        The start time of the measurement.
    * :code:`end_time` : String
        The end time of the measurement.


Verify
~~~~~~

1. New heuristics to exclude possible cases of erroneous answers from resolvers. Currently, :code:`verify` excludes answer IPs that are part of big CDNs (Note: this could lead to false negatives) and answer IPs that appear for a low number of domains (<=2). 
    :code:`results_verified.json` contains all the information when running :code:`full` mode.

    * :code:`vp` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`test_url` : String
        The test domain being queried.
    * :code:`location`: JSON object
        * :code:`country_name` : String
            The full name of the country where the resolver is located.
        * :code:`country_code` : String
            The two-letter ISO 3166 code of the country where the resolver is located.
    * :code:`passed_liveness` : Boolean
            Equals :code:`false` if both control queries were unsuccessful.
    * :code:`in_control_group` : Boolean
            Equals true if at least one control resolver had a valid response for this test domain.
    * :code:`connect_error` : Boolean
            Equals true if all test domain query attempts returned errors. This field is also set to be :code:`true` if the first control measurement fails, and no further measurements for the test domain are conducted. Use this field in conjunction with the :code:`passed_liveness` field to find anomalies.
    * :code:`anomaly` : Boolean
            Equals true if an anomaly is detected. In case there are no tags for the answers or control, then this field is conservatively marked as false. 
    * :code:`start_time` : String
            The start time of the measurement.
    * :code:`end_time` : String
            The end time of the measurement.
    * :code:`response` : JSON object

        The resolver's returned answers for the queried domain are the keys.

        * :code:`url`: String
            The domain being queried in this trial, either the control domain for liveness test or :code:`test_url`. The liveness test DNS responses are only recorded if they do not contain a type-A RR.
        * :code:`has_type_a`: Boolean
            Equals :code:`true` if the query returned a valid A resource record.
        * :code:`error`: String
            Contains error information.
        * :code:`rcode`: Integer
            Response code mapping to success (0) or errors (-1 for connection error, > 0 for errors specified in `RFC 2929 <https://tools.ietf.org/html/rfc2929#section-2.3>`_).
        * :code:`response`: JSON Object
            Consist of a map between IPs the resolver responded for the queried domain and tags from Maxmind:

            * :code:`http` : String
                The hash of the HTTP body.
            * :code:`cert` : String
                The hash of the TLS certificate.
            * :code:`asname` : String
                The autonomous system (AS) name.
            * :code:`asnum` : Integer
                The autonomous system (AS) number.
            * :code:`matched` : Array
                An array of its tags that matched the control tags - if the IP is in the control set, "ip" is appended and if the IP has no tags, "no_tags" is appended.

    * :code:`confidence` : JSON object
        * :code:`average` : Float
            The average percentage of tags matching the control set for the answers (average of :code:`matches`).
        * :code:`matches` : Array
            Contains the percentage of tags matching the control set for each answer. If an answer IP is in the control set, the percentage for that answer is 100 even if the IP has no tags.
        * :code:`untagged_controls` : Boolean
            Equals true if all control IPs for the query have no tags.
        * :code:`untagged_answers` : Boolean
            Equals true if all answer IPs have no tags.
    * :code:`excluded` : Boolean
        Equals :code:`true` if this observation should be excluded from being counted as an anomaly.
    * :code:`exclude_reason` : String Array
        The reasons that this observation should be excluded (eg. "is_CDN")

*****
Notes
*****
While Satellite includes multiple control resolvers intended to avoid false inferences there is still a 
possibility that certain measurements are marked as anomalies incorrectly. To confirm censorship, it is
recommended that the raw DNS responses are compared to known blockpage fingerprints. The blockpage fingerprints
currently recorded by Censored Planet are available `here <https://assets.censoredplanet.org/blockpage_signatures.json>`_.
Moreover, aggregations can be used to avoid anomalous vantage points and domains.  
Please refer to our sample `analysis scripts <https://github.com/censoredplanet/censoredplanet>`_ for a guide on processing 
the data. 

Censored Planet detects network interference of websites using remote measurements to infrastructural vantage points 
within networks (eg. institutions). Note that this raw data cannot determine the entity responsible for the blocking 
or the intent behind it. Please exercise caution when using the data, and reach out to us at `censoredplanet@umich.edu` 
if you have any questions.