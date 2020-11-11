############
Satellite
############
Satellite is Censored Planet's tool to detect DNS interference. Refer to the following papers for more details:

* `Global Measurement of DNS Manipulation <https://censoredplanet.org/assets/Pearce2017b.pdf>`_
* `Satellite: Joint Analysis of CDNs and Network-Level Interference <https://censoredplanet.org/assets/Scott2016a.pdf>`_

*******
Output
*******

The detect module has two output files: :code:`interference.json` and :code:`interference_err.json`.

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