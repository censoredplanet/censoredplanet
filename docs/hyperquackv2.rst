#############
Hyperquack-v2
#############
Hyperquack-v2 measures application-layer keyword blocking,
making use of the Echo, Discard, HTTP, and HTTPS internet protocols to remotely
detect censorship. For a full description of how each of these protocols are
used for remote censorship detection, you can read the following papers:

* `Echo and Discard (Quack) <https://censoredplanet.org/assets/VanderSloot2018.pdf>`_
* `HTTP and HTTPS (Hyperquack) <https://censoredplanet.org/assets/filtermap.pdf>`_

*************
Trial Outputs
*************

Trial outputs are produced when Hyperquack-v2 attempts to measure whether or not
a vantage point observes the blocking of a given keyword by sending one or more
probes to a vantage point.

Fields
======

* :code:`vp` : String
    The IP address of the vantage point used in this trial.
* :code:`location`
    The location of the vantage point. This field has two
    subfields:
    
    * :code:`country_name` : String
        The name of the country the vantage point resides in, as provided by
        the MaxMind geolocation database.
    * :code:`country_code` : String
        The two letter code associated with the aforementioned country.

* :code:`service` : String
    The service of the vantage point we are using for this trial.
    This field is set to the name of the service. 
    If the service is running on a non-standard port,
    a colon and the port number are appended
    (i.e., *discard* for discard on port 9, or *echo:8080* for echo on port 8080).
* :code:`test_url` : String
    The URL being tested by this trial.
* :code:`response`
    An array representing the results of each probe sent to the vantage point.
    Each entry is a JSON object with six subfields:

    * :code:`matches_template` : Boolean
        Each trial performed by Hyperquack-v2 determines whether or not the
        probe was interfered with by comparing the response returned by the
        vantage point to an already known template. If the response does not
        match, that is potentially evidence of blocking. Set to :code:`true`
        if the response given by the vantage point matches the known template,
        and :code:`false` otherwise.
    * :code:`response` : JSON Object
        If the response given by the vantage point does not match the template,
        Hyperquack-v2 will add this field. Describes the response sent by the
        vantage point, including HTTP headers, the HTTP response code, the
        body of the response, and any TLS information. 
    * :code:`error` : String
        If the probe fails with an error, that is potential evidence of
        blocking. If this occurs, this field will be included. Describes the
        encountered error. Note that this field can be used to filter out TCP handshake and setup errors. 
    * :code:`control_url` : String
        During a trial, Hyperquack-v2 will sometimes send probes with
        non-sensitive URLs if all probes with sensitive URLs show
        evidence of being censored. If the probe described by this entry in the
        results array is a control probe, this field will be included. Contains
        the control URL used in the probe.
    * :code:`start_time` : Timestamp
        Th1e time when the probe was sent.
    * :code:`end_time` : Timestamp
        The time when the reponse to the probe response arrived.

* :code:`anomaly` : Boolean
    Indicates whether the probes to the vantage point show enough evidence to
    conclude that the vantage point has observed some sort of anomaly, potentially
    indicative of blocking.
* :code:`controls_failed` : Boolean
    Set to :code:`true` when all control probes sent to the vantage point fail to
    match the known template. This implies that the mismatching responses are
    due to an error in the vantage point or the network, not censorship. Rows with 
    :code:`controls_failed` set to :code:`true` should not be considered for analysis.
* :code:`stateful_block` : Boolean
    Certain methods of censorship will block all communication from a given IP
    address for a length of time after that IP sends a request containing a
    censored keyword. We call this type of censorship ‘Stateful Blocking’. We
    detect this by sending a control probe immediately after our sensitive
    probes, waiting for some time, then sending another control probe. If the
    first control is blocked but the second isn’t, there is potentially
    stateful blocking. If this trial shows evidence of stateful blocking,
    this field is set to :code:`true`.

******************
Evaluation Outputs
******************

Evaluation outputs are produced when the Hyperquack-v2 performs a health
evaluation of a vantage point's service. Services are evaluated by sending one
or more probes containing control keyoword to the vantage point.

Fields
======

* :code:`vp` : String
    The IP address of the vantage point being evaluated..
* :code:`service` : String
    The service of the vantage point that is being evaluated. This field is set
    to the name of the service. If the service is running on a non-standard
    port, a colon and the port number are appended
    (i.e., *discard* for discard on port 9, or *echo:8080* for echo on port 8080).
* :code:`response`
    An array representing the results of each probe sent to the vantage point.
    Each entry is a JSON object with five subfields:

    * :code:`test_url` : String
        The control URL used for this probe.
    * :code:`response` : JSON Object
        If the vantage point responds to the probe, this field is added.
        Describes the response sent by the vantage point, including the HTTP
        headers, the HTTP response code, and the body of the response.
    * :code:`error` : String
        If the probe fails with an error this field is included. Describes the
        encountered error.
    * :code:`start_time` : Timestamp
        The time when the probe was sent.
    * :code:`end_time` : Timestamp
        The time when the reponse to the probe finished arriving.

* :code:`template` : JSON Object
    If Hyperquack-v2 is able to generate a template from the probes, this field
    is included.
    Represents the expected response from the vantage point when sent a probe
    containing an uncensored keyword. If the service being tested is HTTP or 
    HTTPS, this field is an HTTP response, including HTTP headers, the HTTP
    response code, and the body of the response. If the service is Echo or
    Discard, this field is omitted. This template is gereated by the first
    probe during the health evaluation.
* :code:`issue` : String
    If there was an issue in generating the template for this service, this
    field will be included. Describes the issue encountered when generating the
    template or when comparing subsequent control probes to the template.

*************
Notes
*************
While Hyperquack-v2 includes multiple trials intended to avoid random network errors, there is still a 
possibility that certain measurements are marked as anomalies incorrectly. To confirm censorship, it is
recommended that the raw responses are compared to known blockpage fingerprints. The blockpage fingerprints
currently recorded by Censored Planet are available here <https://assets.censoredplanet.org/blockpage_signatures.json>.
Moreover, network errors (such as TCP handshake and Setup errors) must be filtered out to avoid false inferences. 
Please refer to our sample analysis scripts <https://github.com/censoredplanet/censoredplanet> for a guide on processing 
the data. 

Censored Planet detects network interference of websites using remote measurements to infrastructural vantage points 
within networks (eg. institutions). Note that this raw data cannot determine the entity responsible for the blocking 
or the intent behind it. Please exercise caution when using the data, and reach out to us at `censoredplanet@umich.edu` 
if you have any questions.