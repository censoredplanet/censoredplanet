############
Hyperquack
############
Quack and Hyperquack
make use of the Echo, Discard, HTTP, and HTTPS internet protocols to remotely
detect censorship. For a full description of how each of these protocols are
used for remote censorship detection, you can read the following papers:

* `Echo and Discard (Quack) <https://censoredplanet.org/assets/VanderSloot2018.pdf>`_
* `HTTP and HTTPS (Hyperquack) <https://censoredplanet.org/assets/filtermap.pdf>`_

*************
Outputs
*************

Outputs are produced when Hyperquack attempts to measure whether or not
a vantage point observes the blocking of a given keyword by sending one or more
probes to a vantage point.

Fields
======

* :code:`Server` : String
    The IP address of the vantage point used in this trial.
* :code:`Keyword` : String
    The URL being tested by this trial.
* :code:`Retries` : Integer
    The number of times Hyperquack had to resend the test packet in the course
    of the trial. For example, if the first packet sent to the vantage point
    returned the correct response, this field will be set to 0. If that first
    packet does not yield the correct response, every other packet sent by the
    system will increment this field by 1.
* :code:`Results`
    An array representing the results of each probe sent to the vantage point.
    Each entry is a JSON object with six subfields:
    
    * :code:`Sent` : String
        The contents of the packet sent to the vantage point.
    * :code:`Received` : JSON Object
        If the response given by the vantage point does not match the template,
        Hyperquack will add this field. Describes the response sent by the
        vantage point, including HTTP headers, the HTTP response code, and the
        body of the response.
    * :code:`Success` : Boolean
        Each trial performed by Hyperquack determines whether or not the
        probe was interfered with by comparing the response returned by the
        vantage point to an already known template. If the response does not
        match, that is potentially evidence of censorship. Set to :code:`true`
        if the response given by the vantage point matches the known template,
        and :code:`false` otherwise.
    * :code:`Error` : String
        If the probe fails with an error, that is potential evidence of
        censorship. If this occurs, this field will be included. Describes the
        encountered error.
    * :code:`StartTime` : Timestamp
        The time when the probe was sent.
    * :code:`EndTime` : Timestamp
        The time when the reponse to the probe finished arriving.

* :code:`Blocked` : Boolean
    Indicates whether the probes to the vantage point show enough evidence to
    conclude that the vantage point has observed some sort of anomaly, potentially
    indicative of blocking.
* :code:`FailSanity` : Boolean
    Set to :code:`true` when all control probes sent to the vantage point fail to
    match the known template. This implies that the mismatching responses are
    due to an error in the vantage point or the network, not censorship.
* :code:`StatefulBlock` : Boolean
    Certain methods of censorship will block all communication from a given IP
    address for a length of time after that IP sends a request containing a
    censored keyword. We call this type of censorship ‘Stateful Blocking’. We
    detect this by sending a control probe immediately after our sensitive
    probes, waiting for some time, then sending another control probe. If the
    first control is blocked but the second isn’t, there is potentially
    stateful blocking. If this trial shows evidence of stateful blocking,
    this field is set to :code:`true`.
