package validation

import (
	"bytes"
	"errors"
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
	"log"
	"regexp"
	"strings"
)

const topLevelStr = "toplevel"

// To correlate matching segments with segments in message
type IndexSegmentMap map[int]*msg.Segment

// For looking up segments e.g. in query
type segExprToSegMap map[string]*msg.SegmentOrGroup

type MessageValidator struct {
	messageSpec                *msgspec.MessageSpec
	segmentValidationRegexpStr string
	segmentValidationRegexp    *regexp.Regexp
}

func (v *MessageValidator) String() string {
	return fmt.Sprintf("MessageValidator (%s)", v.String())
}

// Creates a repeat str for a regular expression. If possible, a
// simple representation like the Kleene star *, or + is chosen
// to avoid a golang ErrInvalidRepeatSize error in nested groups
func getRegexpRepeatStr(minSpecRepeat int, maxSpecRepeat int, isGroup bool) (result string) {
	var maxRepeat int

	// TODO: postprocessing step to check specified counts; actual counts can be higher
	// than what is representable with regexp repeat counts (which uses a hard limit
	// to avoid performance problem with quadratic behavior due to nesting of repeated groups)
	if isGroup {
		maxRepeat = 5
	} else {
		maxRepeat = 5
	}

	if minSpecRepeat == 1 && maxSpecRepeat == 1 {
		return ""
	} else if minSpecRepeat == 0 && maxSpecRepeat == 1 {
		return "*"
	} else {
		if maxSpecRepeat >= maxRepeat {
			if maxSpecRepeat > maxRepeat {
				log.Printf(
					"Segment spec count %d exceeds implementation limit %d",
					maxSpecRepeat, maxRepeat)
			}
			if minSpecRepeat == 0 {
				return "*"
			} else if minSpecRepeat == 1 {
				return "+"
			}
		} else {
			return fmt.Sprintf("{%d,%d}", minSpecRepeat, maxSpecRepeat)
		}
	}
	panic("unreachable")
}

func buildMessageSpecPartRegexpStr(msgSpecPart msgspec.MessageSpecPart, groupName string) string {
	var inner string
	subMatchNameStr := ""
	var regexpRepeatStr string
	specMinCount := msgSpecPart.MinCount()
	specMaxCount := msgSpecPart.MaxCount()

	switch msgSpecPart := msgSpecPart.(type) {
	case *msgspec.MessageSpecSegmentPart:
		inner = msgSpecPart.SegmentSpec.Id + `-[0-9]+:`
		regexpRepeatStr = getRegexpRepeatStr(specMinCount, specMaxCount, false)
		subMatchNameStr = groupName
	case *msgspec.MessageSpecSegmentGroupPart:
		groupPartRegexpStrs := []string{}
		for _, groupChild := range msgSpecPart.Children() {
			groupPartRegexpStrs = append(
				groupPartRegexpStrs,
				buildMessageSpecPartRegexpStr(groupChild, msgSpecPart.Name()))
		}
		inner = strings.Join(groupPartRegexpStrs, "")
		regexpRepeatStr = getRegexpRepeatStr(specMinCount, specMaxCount, true)
		// subMatchNameStr = msgSpecPart.Name()
	default:
		panic("Not implemented")
	}

	if subMatchNameStr != "" {
		subMatchNameStr = fmt.Sprintf("?P<%s>", subMatchNameStr)
	}
	return fmt.Sprintf("(%s%s)%s", subMatchNameStr, inner, regexpRepeatStr)
}

func buildMessageSpecPartsRegexpStr(msgSpecParts []msgspec.MessageSpecPart) string {
	var buf bytes.Buffer
	buf.WriteString("^")
	for _, part := range msgSpecParts {
		buf.WriteString(buildMessageSpecPartRegexpStr(part, topLevelStr))
	}
	buf.WriteString("$")
	return buf.String()
}

// Build a regular expression for validating sequences of segments
// against the message specification.
// Each segment is encoded as XXX: where XXX is the segment ID and ':'
// is a separator to avoid misaligned matches
func buildSegmentSeqValidationRegexp(msgSpec *msgspec.MessageSpec) (msgRegexpStr string, msgRegexp *regexp.Regexp, err error) {
	regexpStr := buildMessageSpecPartsRegexpStr(msgSpec.Parts)
	log.Printf("regexp str: '%s'", regexpStr)
	msgRegexp = regexp.MustCompile(regexpStr)
	return regexpStr, msgRegexp, nil
}

// Contructs a string out of the segment sequence, so regexp matching can be
// used to validate the segment sequence.
// Also constructs a mapping of integer indices to segments, which can be
// used to locate a particular segment by index
func buildSegmentListStr(segments []*msg.Segment) (segmentListStr string, indexSegmentMap IndexSegmentMap) {
	var buf bytes.Buffer
	indexSegmentMap = map[int]*msg.Segment{}
	for index, segment := range segments {
		buf.WriteString(fmt.Sprintf("%s-%d:", segment.Id(), index))
		indexSegmentMap[index] = segment
	}
	return buf.String(), indexSegmentMap
}

func (v *MessageValidator) Validate(message msg.Message) (isValid bool, err error) {
	panic("Not implemented")
}

/*// Builds a map for looking up either stand-alone segments or segment groups
func (v *MessageValidator) buildSegExprToSegMap(indexSegmentMap IndexSegmentMap, segSeqMatch []string) segExprToSegMap {
	result = IndexSegmentMap{}
	segMap := map[string]
	for matchIndex, matchName := range segSeqMatch {
		if matchName == "" {
			panic("Internal error: unnamed submatch")
		}
		if !strings.HasPrefix(matchName, "Group") {


		}
	}
	return result
}*/

// Validate a list of segment names as they occur in a message
func (v *MessageValidator) ValidateSegmentList(segments []*msg.Segment) (isValid bool, err error) {
	if len(segments) == 0 {
		return false, errors.New("No segments")
	}
	segmListStr, indexSegmentMap := buildSegmentListStr(segments)
	// log.Printf("segmListStr: %s\n", segmListStr)
	// log.Printf("regexp str: '%q'\n", v.segmentValidationRegexpStr)
	log.Printf("indexSegmentMap: %v\n", indexSegmentMap)
	log.Printf("group names: %q\n", v.segmentValidationRegexp.SubexpNames())
	_ = indexSegmentMap
	match := v.segmentValidationRegexp.FindStringSubmatch(segmListStr)
	log.Printf("match: %#v\n", match)
	return len(match) != 0, nil
}

func NewMessageValidator(messageSpec *msgspec.MessageSpec) (validator *MessageValidator, err error) {
	msgRegexpStr, msgRegexp, err := buildSegmentSeqValidationRegexp(messageSpec)
	if err != nil {
		return nil, err
	}
	return &MessageValidator{
		messageSpec:                messageSpec,
		segmentValidationRegexpStr: msgRegexpStr,
		segmentValidationRegexp:    msgRegexp,
	}, nil
}
