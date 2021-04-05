############
Satellite-v1
############
Satellite is Censored Planet's tool to detect DNS interference. Refer to the following papers for more details:

* `Global Measurement of DNS Manipulation <https://censoredplanet.org/assets/Pearce2017b.pdf>`_
* `Satellite: Joint Analysis of CDNs and Network-Level Interference <https://censoredplanet.org/assets/Scott2016a.pdf>`_

Satellite-v1 corresponds to measurements from 2018 to February 2021. See Satellite v2 for documentation on measurements taken after this date.

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

------
Query
------

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

------
Tag
------

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

------
Detect
------

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


*************
Notes
*************
While Satellite-v1 includes multiple control resolvers intended to avoid false inferences there is still a 
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
