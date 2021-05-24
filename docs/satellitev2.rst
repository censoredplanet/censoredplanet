############
Satellite-v2
############
Satellite is Censored Planet's tool to detect DNS interference. Refer to the following papers for more details:

* `Global Measurement of DNS Manipulation <https://censoredplanet.org/assets/Pearce2017b.pdf>`_
* `Satellite: Joint Analysis of CDNs and Network-Level Interference <https://censoredplanet.org/assets/Scott2016a.pdf>`_

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

------
Probe
------

1. Generate a DNS A query packet for a controlled domain (:code:`dns.pkt`).

2. Perform a `ZMap <https://github.com/zmap/zmap>`_ (Internet-wide) scan with the probe packet for open DNS resolvers.

    :code:`resolvers_raw.json` contains the ZMap output:

    * :code:`saddr` : String
        IP address of a DNS resolver.
    * :code:`data` : String
        Raw response to probe domain.

------
Filter
------

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

------
Query
------

1. Make DNS queries for each test domain to each resolver. The query for the test domain is attempted up to four times in case of connection error. To check the status of the resolver, a control measurement is conducted before the queries for the test domain. If the first control measurement fails, no further measurements will be conducted for the same :code:`<resolver-domain>` pair. If all 4 trials for the test domain fail, another control measurement will be conducted.

2. Parse and separate responses from control resolvers and non-control resolvers.

------
Tag
------

1. Tag each answer IP with information from `Censys <https://about.censys.io/>`_.
    **Note:**

        * Fields may have empty strings if the information was not available on Censys.

------
Detect
------

1. Compare query responses between non-control resolvers and control resolvers to identify interference. When running satellite v2 as a whole module, :code:`detect` does not output any files. However, when run separately, :code:`detect` outputs :code:`results.json` with the :code:`excluded` field set to :code:`false` and the :code:`excluded_reason` field set to :code:`null` by default. (See the output structure in :code:`verify` section)

    **Note:**

        * For each response, the answer IPs and their tags are compared to the set of answer IPs and tags from all the control resolvers for the same query domain. A response is classified as an anomaly if there is no overlap between the two.

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
        Equals true if a page is successfully fetched.
    * :code:`start_time` : String
        The start time of the measurement.
    * :code:`end_time` : String
        The end time of the measurement.

------
Verify
------

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
        

*************
Notes
*************
While Satellite-v2 includes multiple control resolvers intended to avoid false inferences there is still a 
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