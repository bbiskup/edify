Ideas for path syntax for referencing elements or other parts of msg
====================================================================

Specification
-------------

Main specification documents

- http://www.unece.org/trade/untdid/texts/d423.htm
- http://www.unece.org/trade/untdid/download/r1241.txt

Spec archive: http://www.unece.org/fileadmin/DAM/trade/untdid/d14b/d14b.zip

Information about files contained in spec archive: 
.edify/downloads/d14b/_TextFiles/CONTENT.TXT


Reference message
-----------------

(Indentation for demonstration only)

UNA:+.? '
UNB+UNOC:3+Senderkennung+Empfaengerkennung+060620:0931+1++1234567'
    UNH+1+ORDERS:D:96A:UN'
        BGM+220+B10001'
        DTM+4:20060620:102'
        NAD+BY+++Bestellername+Strasse+Stadt++23436+xx'
        LIN+1++Produkt Schrauben:SA'
        QTY+1:1000'
        UNS+S'
        CNT+2:1'
    UNT+9+1'
UNZ+1+1234567'

Query language
--------------

Querying using spec information

msg:ORDERS[0]|seg:BGM[comp:C104:0]|[simp:1004]

 --> "B10001"

msg:ORDERS[0]|seg:BGM[comp:C104:0]

 --> CompositeDataElement (with component elements)

msg:ORDERS[0]

  -> Message object

msg:ORDERS[*]

  -> [Message object]   (potentially many)

TODO: addressing segment groups

- Precompute paths?
  - only a small subset of data elements may ever get requested
  - message layout might change (e.g. by adding groups)
  --> better: on-the fly

- only fully qualified (abs.) paths; selection only by data element amost never makes sense


Web application examples
------------------------

e.g. http://www.truugo.com/edifact/d14b/

Validation & Navigation
=======================

Message segment sequence validation
-----------------------------------

Basic operations while matching message sequence:

- match segment       -> XX, increment current segment spec index
- repeat segment      -> increment repeat count 
- skip segment        -> increment current segment index
- descend into group  -> recursive call findSeg
- finish group (when: end of msg or segment not found in current group) -> return from findSeg
- repeat group        -> loop within findSeg
- skip group          -> increment current segment spec index
- (msg end ?          -> spec complete (no more mandatory segment/group) -> ok
                         incomplete:                                     -> error missing_mandatory_segment