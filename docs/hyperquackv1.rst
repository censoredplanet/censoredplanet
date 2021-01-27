#############
Hyperquack-v1
#############
Quack-v1 and Hyperquack-v1
make use of the Echo, Discard, HTTP, and HTTPS internet protocols to remotely
detect application-layer keyword blocking. For a full description of how each of these protocols are
used for remote censorship detection, you can read the following papers:

* `Echo and Discard (Quack) <https://censoredplanet.org/assets/VanderSloot2018.pdf>`_
* `HTTP and HTTPS (Hyperquack) <https://censoredplanet.org/assets/filtermap.pdf>`_

Hyperquack-v1 corresponds to measurements from July 2018 to February 2021. See Hyperquack-v2 for documentation on measurements taken after this date.

The published data has the following directory structure: ::

    CP_Quack-PROTOCOL-YYYY-MM-DD-HH-MM-SS/
    |-- servers.txt (List of vantage points used for measurements)
    |-- domains.txt (Test URLs)
    |-- results.json (Output file)
   
*************
Outputs
*************

The relevant output is located in `results.json`


Fields
======

* :code:`Server` : String
    The IP address of the vantage point used in this trial.
* :code:`Keyword` : String
    The URL being tested by this trial.
* :code:`Retries` : Integer
    The number of times Hyperquack had to resend the test packet in the course
    of the trial. For example, if the first packet sent to the vantage point
    returned the expected response, this field will be set to 0. If the first
    packet does not yield the expected response, every other packet sent by the
    system will increment this field by 1.
* :code:`Results`
    An array representing the results of each probe sent to the vantage point.
    Each entry is a JSON object with six subfields:
    
    * :code:`Sent` : String
        The contents of the packet sent to the vantage point. Note that this field 
        may not be populated in cases where there is no application-layer response 
        from the vantage point. In case the TCP handshake with the vantage point fails, 
        the error field will be set accordingly.  
    * :code:`Received` : JSON Object
        If the response given by the vantage point does not match the template,
        Hyperquack-v1 will add this field. Describes the response sent by the
        vantage point, including HTTP headers, the HTTP response code, the
        body of the response, and any TLS information. 

    * :code:`Success` : Boolean
        Each trial performed by Hyperquack determines whether or not the
        probe was interfered with by comparing the response returned by the
        vantage point to an already known template. If the response does not
        match, that is potentially evidence of interference. Set to :code:`true`
        if the response given by the vantage point matches the known template,
        and :code:`false` otherwise.
    * :code:`Error` : String
        If the probe fails with an error, that is potential evidence of
        blocking. If this occurs, this field will be populated. Describes the
        encountered error. Note that this field can be used to filter out TCP handshake and setup errors. 
    * :code:`StartTime` : Timestamp
        The time when the probe was sent.
    * :code:`EndTime` : Timestamp
        The time when the reponse to the probe response arrived.

* :code:`Blocked` : Boolean
    Indicates whether the probes to the vantage point show enough evidence to
    conclude that the vantage point has observed some sort of anomaly, potentially
    indicative of blocking.
* :code:`FailSanity` : Boolean
    Set to :code:`true` when all control probes sent to the vantage point fail to
    match the known template. This implies that the mismatching responses are
    due to an error in the vantage point or the network, not censorship. Rows with 
    :code:`FailSanity` set to :code:`true` should not be considered for analysis. 
* :code:`StatefulBlock` : Boolean
    Certain methods of censorship will block all communication from a given IP
    address for a length of time after that IP sends a request containing a
    censored keyword. We call this type of censorship ‘Stateful Blocking’. We
    detect this by sending a control probe immediately after our sensitive
    probes, waiting for some time (2 minutes), then sending another control probe. If the
    first control is blocked but the second isn’t, there is potentially
    stateful blocking. If this trial shows evidence of stateful blocking,
    this field is set to :code:`true`.

*************
Notes
*************
While Hyperquack-v1 includes multiple trials intended to avoid random network errors, there is still a 
possibility that certain measurements are marked as anomalies incorrectly. To confirm censorship, it is
recommended that the raw responses are compared to known blockpage fingerprints. The blockpage fingerprints
currently recorded by Censored Planet are available `here <https://assets.censoredplanet.org/blockpage_signatures.json>`_.
Moreover, network errors (such as TCP handshake and Setup errors) must be filtered out to avoid false inferences. 
Please refer to our sample `analysis scripts <https://github.com/censoredplanet/censoredplanet>`_ for a guide on processing 
the data. 

Censored Planet detects network interference of websites using remote measurements to infrastructural vantage points 
within networks (eg. institutions). Note that this raw data cannot determine the entity responsible for the blocking 
or the intent behind it. Please exercise caution when using the data, and reach out to us at `censoredplanet@umich.edu` 
if you have any questions.


