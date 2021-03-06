TODO
====

- Code examples (as Example* functions and / or comments; see codegangsta/cli)

- Absolute & relative navigation (relative e.g. with respect to a repeat group, to avoid overhead)

- Support repeating data elements; example: PAXLST_1.txt and PAXLST_2.txt, first COM data element
  COM+703-555-1212:TE+703-555-4545:FX'

- Handle multiple segments of same ID in same group correctly
  Example MOVINS message, group 6: LOC-RFF-FTX-MEA-DIM-LOC-NAD-SG7-SG8-SG9-SG10-SG11
  (LOC segment occurs twice)

- Propagate annotated error messages (http://rogpeppe.neocities.org/error-loving-talk/#26)

- Obtain separator chars from interchange UNA segment (UNA:+.? ')

- Maintain segment group information for addressing via path
  - segment sequence regex: extract group information

- Repeating data elements (r1241.txt, ch. 8.7.3; other occurences)
  - "repetition separator"

- Separate postprocessing step or match callback to validate higher repeat
  counts of segments in message spec (repeat counts that are too high
  for the golang regexp parser)

- Obtain UNH and other missing definitions from EDIFACT spec? (missing from UNCE specs)
  - 2015-05-03: UNH, UNT, UNS, UGH, UGT missing
  - See testdata/d14b/_TextFiles/CONTENT.TXT:

    (2) Chapter 2 of Part 3 does not include syntax service messages
    (e.g., CONTRL) of UN/EDIFACT.

    (3) Chapter 3 of Part 3 does not include service segments (tags
    beginning with Uxx) defined in ISO 9735 (the EDIFACT syntax).
    Version control of these service segments is reflected in data
    element 0002 in the interchange header segment and is based on change
    to ISO 9735. Therefore the usual UN/EDIFACT directory version/release
    procedures for UN/EDIFACT messages (using data elements 0052 and 0054
    in the message header segment and the functional group header segment)
    is NOT applicable to those segments.

    (4) Chapter 4 of Part 5 does not include service composite data
    elements (the "Sxxx" series) which are defined in ISO 9735 (the
    EDIFACT syntax). Version control for these composite data elements
    is reflected in data element 0002 in the interchange header segment
    and is based on changes to ISO 9735. Therefore the usual UN/EDIFACT 
    directory version/release procedures for UN/EDIFACT messages (using 
    data elements 0052 and 0054 in composite S008, S009 and S306) is NOT
    applicable to those composite data elements.

    (5) Chapter 5 of Part 3 does not include service data elements (the
    "0xxx" series) which are defined in ISO 9735 (the EDIFACT syntax).
    Version control for these data elements is reflected in data element
    0002 in the interchange header segment and is based on changes to 
    ISO 9735. Therefore the usual UN/EDIFACT directory version/release 
    procedures for UN/EDIFACT messages (using data elements 0052 and 0054)
    is NOT applicable to those data elements.