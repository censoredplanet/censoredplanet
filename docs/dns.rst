####################
DNS Data - Satellite
####################

`Satellite/Iris <https://censoredplanet.org/projects/satellite>`_ is Censored Planet’s remote measurement technique that detects DNS interference using Open DNS resolvers. Below, we provide an overview of Satellite and its data format. Refer to our `academic papers <https://censoredplanet.org/assets/Pearce2017b.pdf>`_ for in-depth details about Satellite.

******************
Satellite-v2.2-raw
******************

To provide raw data for easy data analysis, we made the following changes:

1. Split data based on the country of resolvers so that it is easier to select and download data according to users' country of interest.

2. Separated the data collection phase and data analysis phase. Right now the Satellite data from our `raw measurement data website <https://data.censoredplanet.org/raw>`_ is truthful to the data collected without further analysis. We deprecated the “anomaly” field since there are misunderstandings that anomaly represents censorship.

3. Added new data containing further metadata fields and flattened nested data for easy analysis. Modified field names for disambiguation purposes.

    * :code:`domain` : String
        The test domain being queried.
    * :code:`domain_is_control` : Boolean
        Equals true if the queried domain is the root server for liveness test.
    * :code:`test_url` : String
        The URL of the queried domain.
    * :code:`date` : String
            The date of the measurement.
    * :code:`start_time` : String
            The start time of the measurement.
    * :code:`end_time` : String
            The end time of the measurement.
    * :code:`resolver_ip` : String
        The IP address of the vantage point (a DNS resolver).
    * :code:`resolver_name` : String
        The hostname of the vantage point.
    * :code:`resolver_is_trusted` : Boolean
        Equals true if the resolver is a control resolver.
    * :code:`resolver_netblock` : String
        The netblock the vantage point belongs to.
    * :code:`resolver_asn` : String
        The AS number of the AS the vantage point resides in.
    * :code:`resolver_as_name` : String
        The name of the AS the vantage point resides in.
    * :code:`resolver_as_full_name` : String
        The full name of the AS the vantage point resides in.
    * :code:`resolver_as_class` : String
        The class of the AS the vantage point resides in.
    * :code:`resolver_country` : String
        The country the vantage point resides in.
    * :code:`resolver_organization` : String
        The IP organization the vantage point resides in.
    * :code:`received_error` : String
        Flatten error messages from the received responses.
    * :code:`received_rcode` : Integer
        Flatten rcode from the received responses. Response code mapping to success (0) or errors (-1 for connection error, > 0 for errors specified in `RFC 2929 <https://tools.ietf.org/html/rfc2929#section-2.3>`_).
    * :code:`source` : String
        Tar file name of the measurement.
    * :code:`answers` : JSON object
        The resolver's returned answers for queried domain.

        * :code:`ip`: String
            Returned IP.
        * :code:`asn`: String
            The AS number of the AS the returned IP resides in.
        * :code:`as_name`: String
            The AS name of the AS the returned IP resides in.
        * :code:`censys_http_body_hash`: String
            The hash of the HTTP body from Censys.
        * :code:`censys_ip_cert`: String
            The hash of the TLS certificate from Censys.
        * :code:`http_error`: String
            Parsed HTTP page error message from :code:`fetch` module.
        * :code:`http_response_status`: String
            Parsed HTTP page status code from :code:`fetch` module.
        * :code:`http_response_headers`: String
            Parsed HTTP page headers from :code:`fetch` module.
        * :code:`http_response_body`: String
            Parsed HTTP page body from :code:`fetch` module.
        * :code:`https_error`: String
            Parsed HTTPS page error message from :code:`fetch` module.
        * :code:`https_response_status`: String
            Parsed HTTPS page status code from :code:`fetch` module.
        * :code:`https_response_headers`: String
            Parsed HTTPS page headers from :code:`fetch` module.
        * :code:`https_response_body`: String
            Parsed HTTPS page body from :code:`fetch` module.
        * :code:`https_tls_version`: String
            Parsed TLS version from :code:`fetch` module.
        * :code:`https_tls_cipher_suite`: String
            Parsed TLS cipher suite from :code:`fetch` module.
        * :code:`https_tls_cert`: String
            Parsed TLS certificate from :code:`fetch` module.
        * :code:`https_tls_cert_common_name`: String
            Parsed common name field from TLS certificate.
        * :code:`https_tls_cert_alternative_names`: String
            Parsed alternative name field from TLS certificate.
        * :code:`https_tls_cert_issuer`: String
            Parsed issuer field from TLS certificate.
        * :code:`https_tls_cert_start_date`: String
            Parsed start date of the TLS certificate.
        * :code:`https_tls_cert_end_date`: String
            Parsed end date of the TLS certificate.


*************************
Satellite-v1 (deprecated)
*************************

.. image:: images/satellite-v1.png
  :width: 600
  :alt: Figure - Overview of Satellite-v1

Figure - Overview of Satellite-v1

Satellite-v1 is the first version of Satellite that we operated from August 2018 - February 2021. The primary function of Satellite is to detect incorrect DNS resolutions from open DNS resolvers in many countries.

* From a measurement machine at the University of Michigan, we send a DNS query for a website whose reachability we’re interested in, to an open DNS resolver in a country of interest (1). The response from the DNS resolver is our Test IP (2).

* We also send a DNS query for the same website to trusted control resolvers (3), and record their response as the control IP (4).

* We then compare the test and control responses using several heuristics, including a direct IP address comparison, and comparison of the AS number, AS names, HTTP content hashes, and TLS certificates associated with the test and control IP addresses (5). Satellite-v1 only labels a measurement as an anomaly when all of the heuristics mismatch.

Our various `publications <http://censoredplanet.org/publications>`_ and `reports <http://censoredplanet.org/reports>`_ have used Satellite-v1 to detect many cases of DNS manipulation. For instance, in our `recent investigation into the filtering of COVID-19 websites <https://censoredplanet.org/assets/covid.pdf>`_ , Satellite-v1 found many networks using website filtering products to manipulate DNS responses of COVID-related websites.

Limitations
***********

Although Satellite-v1 was extremely useful in detecting DNS interference at large scale, it suffered from several limitations, which form the improvements in Satellite-v2.x.

* Satellite-v1 could not detect DNS censorship where A records were not available i.e. Satellite-v1 primarily focused on detecting incorrect DNS resolutions through the resolved IP address, and did not contain heuristics to measure DNS manipulation which manifested through timeouts, NXDOMAIN responses, SERVFAIL responses, etc.

* Satellite-v1 required post-processing to remove false positives and confirm the presence of anomalies, such as through using post-measurement heuristics and blockpage regexes. Satellite-v2 has the inbuilt capability to perform most post-processing measurements.


*************************
Satellite-v2 (deprecated)
*************************

.. image:: images/satellite-v2.png
  :width: 600
  :alt: Figure - Overview of Satellite-v2.0

Figure - Overview of Satellite-v2

Satellite-v2 is our brand new version of Satellite, where we’ve made several modifications to the measurement technique and data format for facilitating accurate and efficient remote DNS interference measurements. Below, we detail the major changes we’ve made in Satellite-v2.

* **Fetching HTML pages hosted at resolved IPs marked as an anomaly** -  Satellite-v2 has an in-built fetch feature that performs HTTP and HTTPS GET requests to resolved IPs that fail our heuristics. This step was being performed as a post-processing step in Satellite-v1. This addition helps in quickly identifying blockpages such as the example shown in the figure below. Moreover, we are in the process of developing a technique to use TLS certificates to detect DNS manipulation. Reach out to `censoredplanet@umich.edu` for more information.

* **Measuring DNS interference without A records** - In Satellite-v2, we have added a sandwiched retry mechanism to our Satellite measurements in order to detect DNS interference that results in a non-zero R code response. A description of the method is shown in the figure below. We first make a control query to the open DNS resolver, providing a domain name that we do not expect to be blocked (eg. www.example.com). After the control query, we make up to 4 retries of the test DNS query, providing the test domain name. In case an A record is detected, we stop the test measurement. At the end, we perform another control query similar to the first measurement. The control queries ensure that the resolver is behaving correctly for an innocuous domain, and the multiple retry mechanism accounts for temporary errors in the network. With the help of the sandwiched retry mechanism, Satellite-v2 is able to detect DNS interference that manifests as timeouts, NXDOMAIN, SERVFAIL etc. From our preliminary analysis of Satellite-v2 data, we’ve already found several cases of DNS interference that can be identified using this method. For example, from the Satellite-v2 scan performed on 2021-03-17, we are able to identify 174,795 responses that have non-zero R codes from China, which makes up 15.6% out of the responses marked as interference. This kind of DNS interference was previously omitted by satellite v1. Shown below is an example measurement that passed the sandwich control tests, but received server failure R code. This could be an indicator of censorship or geoblocking.

* **Adding scan-level heuristics to exclude false positives** - Another step part of the post-processing pipeline of Satellite-v1 that is inbuilt in Satellite-v2. We exclude potentially false positive anomalies by using scan-level heuristics, such as the number of domains resolving to the anomalous IP address, or the anomalous IP address being part of a big CDN. Note that this step may lead to Satellite-v2 missing certain censorship.

* **Other changes** - We updated the heuristics to determine whether a DNS response is interfered - Satellite-v2 now includes a new “confidence” field, which addresses the certainty of interference according to the state of comparison between responses from the test resolvers and the control resolvers. We also make sure that IPs with no metadata information from Censys are not marked as interference.

**************
Satellite-v2.1
**************
Satellite-v2.1 incorporates minor changes from Satellite-v2.0, starting after April 14, 2021. Most of these changes are related to change in data formats. 

**************
Satellite-v2.2
**************
Satellite-v2.2 incorporates major changes in code and data structure from Satellite-v2.1, but no major changes in the functionality of Satellite. The changes are made after June 7, 2021 and they include,

* Store information generated from the query, tag, detect, and verify module in memory, producing only one file (results.json) as output, instead of generating outputs for every module. Renamed query-tag-detect-verify as “test” module, and probe-filter as “discovery”.

* Updated test module so that it first conducts queries for control resolvers, and then query, tag and detect test resolvers in batches.


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
    :code:`results.json` contains all the information when running :code:`full` mode.

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


*****
Notes
*****
While Satellite includes multiple control resolvers intended to avoid false inferences there is still a 
possibility that certain measurements are marked as anomalies incorrectly. To confirm censorship, it is critical 
that the raw DNS responses are compared to known blockpage fingerprints. The blockpage fingerprints
currently recorded by Censored Planet are available `here <https://assets.censoredplanet.org/blockpage_signatures.json>`_.
Moreover, aggregations can be used to avoid anomalous vantage points and domains.  
Please use our `analysis pipeline <https://github.com/censoredplanet/censoredplanet-analysis>`_ 
to process the data before using it.

Censored Planet detects network interference of websites using remote measurements to infrastructural vantage points 
within networks (eg. institutions). Note that this raw data cannot determine the entity responsible for the blocking 
or the intent behind it. Please exercise caution when using the data, and reach out to us at `censoredplanet@umich.edu` 
if you have any questions.