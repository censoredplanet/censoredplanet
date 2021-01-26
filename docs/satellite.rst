############
Satellite
############
Satellite is Censored Planet's tool to detect DNS interference. Refer to the following papers for more details:

* `Global Measurement of DNS Manipulation <https://censoredplanet.org/assets/Pearce2017b.pdf>`_
* `Satellite: Joint Analysis of CDNs and Network-Level Interference <https://censoredplanet.org/assets/Scott2016a.pdf>`_

*******
Output
*******

The detect module is the module that provides final output. It has two output files: :code:`interference.json` and :code:`interference_err.json`.

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
* :code:`startTime` : String
    The start time of the measurement.
* :code:`endTime` : String
    The end time of the measurement
* :code:`confidence` : JSON object
    * :code:`average` : Float
        Average percentage of tags matching the control set for the answers (average of :code:`matches`).
    * :code:`matches` : Array
        Contains the percentage of tags matching the control set for each answer. If an answer IP is in the control set, the percentage for that answer is 100 even if the IP has no tags.
    * :code:`untaggedControls` : Boolean
        Equals true if all control IPs for the query have no tags. Cases where :code:`passed` is false (inteference detected) and :code:`untaggedControls` is true should be checked, since the interference classification is due to differing IPs.
    * :code:`untaggedAnswers` : Boolean
        Equals true if all answer IPs have no tags.


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
        * answers from control resolvers (:code:`answers_control.json`)
        * erroneous query responses (:code:`answers_err.json`)
        * IP list of answers (:code:`answers_ip.json`)
        * non-erroneous non control resolver answers (:code:`answers.json`)
        * non-erroneous raw response packets (:code:`answer_raw.json`)
* :code:`tag`:  tags resolvers with MaxMind (country) and IPs with censys (certificate, AS number and AS name).
    * input: 
        * list of resolvers to query (:code:`resolvers.json`)
        * list of answered IPs (:code:`answers_ip.json`)
    * output:
        * list of tagged IPs with Censys (:code:`tagged_answers.json`) 
        * list of tagged resolvers with Maxmind (:code:`tagged_resolvers.json`)
* :code:`detect`: detects interference by comparing DNS query responses to control set. 
    * input:
        * :code:`tagged_answers.json` 
        * answers from control resolvers (:code:`answers_control.json`)
        * :code:`assets/satellite/control_resolvers.txt` 
    * output:
        * list of interference result (:code:`interference.json`)
        * list of tuples where control set have no same queries (:code:`interference_err.json`) 
* :code:`fetch`: fetches pages hosted on the IPs identified as interference for future blockpage analysis.
    * input:
        * :code:`interference.json`
    * output:
        * list of tampered IPs, and results of HTTP(S) GET (:code:`blockpages.json`)
* :code:`stat`:   data analysis.
* :code:`full`:   all aforementioned modules combined.
* :code:`upload`