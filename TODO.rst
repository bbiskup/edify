TODO
====

- Parallel parsing of segment specs

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