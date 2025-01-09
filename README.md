# godm

ODM handling in Go

Implements the CCSDS ODM standard, April 2023 Edition: https://public.ccsds.org/Pubs/502x0b3e1.pdf

## OPM - Orbit Parameter Message

> An OPM specifies the position and velocity of a single object at a specified epoch.
> Optionally, osculating Keplerian elements may be provided. It should be noted that a
> sequence of OPMs for either a single object or for multiple objects can be aggregated into a
> single Navigation Data Message (NDM) XML file as described in 8.12 and shown in
> annex G. This message is suited to exchanges that (1) involve automated interaction and/or
> human interaction, and (2) do not require high-fidelity dynamic modeling.
> The OPM requires the use of a propagation technique to determine the position and velocity
> at times different from the specified epoch, leading to a higher level of effort for software
> implementation than for the OEM.
>
> The OPM also contains an optional 6x6 position/velocity covariance matrix that reflects the
> uncertainty of the orbit state and may be used in the propagation process to estimate future
> uncertainties.
>
> The OPM allows for modeling of any number of maneuvers (as both finite and instantaneous
> events) and simple modeling of solar radiation pressure and atmospheric drag.
> Though primarily intended for use by computers, the attributes of the OPM also make it
> suitable for applications such as exchanges by email, FAX, or voice, or applications in which
> the message is to be frequently interpreted by humans.

CCSDS ODM Standard, April 2023 Edition, Section 2.1
