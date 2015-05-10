package validation

import (
	"bytes"
	"errors"
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
	"regexp"
	"strings"
)

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
		maxRepeat = 99
	} else {
		maxRepeat = 99
	}

	if minSpecRepeat == 1 && maxSpecRepeat == 1 {
		return ""
	} else if minSpecRepeat == 0 && maxSpecRepeat == 1 {
		return "*"
	} else {
		if maxSpecRepeat >= maxRepeat {
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

/*
func buildMessageSegmentSpecPartRegexpStr(msgSpecPart msgspec.MessageSpecSegmentPart) string {
	specMinCount := msgSpecPart.MinCount()
	specMaxCount := msgSpecPart.MaxCount()
	repeatStr := getRegexpRepeatStr(specMinCount, specMaxCount, msgSpecPart.IsGroup())
	return repeatStr
}*/

func buildMessageSpecPartRegexpStr(msgSpecPart msgspec.MessageSpecPart) string {
	var inner string
	var regexpRepeatStr string
	specMinCount := msgSpecPart.MinCount()
	specMaxCount := msgSpecPart.MaxCount()

	switch msgSpecPart := msgSpecPart.(type) {
	case *msgspec.MessageSpecSegmentPart:
		inner = msgSpecPart.SegmentSpec.Id + ":"
		regexpRepeatStr = getRegexpRepeatStr(specMinCount, specMaxCount, false)
	case *msgspec.MessageSpecSegmentGroupPart:
		groupPartRegexpStrs := []string{}
		for _, groupChild := range msgSpecPart.Children() {
			groupPartRegexpStrs = append(groupPartRegexpStrs, buildMessageSpecPartRegexpStr(groupChild))
		}
		inner = strings.Join(groupPartRegexpStrs, "")
		regexpRepeatStr = getRegexpRepeatStr(specMinCount, specMaxCount, true)
	default:
		panic("Not implemented")
	}

	return fmt.Sprintf("(%s)%s", inner, regexpRepeatStr)
}

func buildMessageSpecPartsRegexpStr(msgSpecParts []msgspec.MessageSpecPart) string {
	var buf bytes.Buffer
	buf.WriteString("^")
	for _, part := range msgSpecParts {
		buf.WriteString(buildMessageSpecPartRegexpStr(part))
	}
	buf.WriteString("$")
	return buf.String()
}

// Build a regular expression for validation sequences of segments
// against the message specification.
// Each segment is encoded as XXX: where XXX is the segment ID and ':'
// is a separator to avoid misaligned matches
func buildSegmentSeqValidationRegexp(msgSpec *msgspec.MessageSpec) (msgRegexpStr string, msgRegexp *regexp.Regexp, err error) {
	regexpStr := buildMessageSpecPartsRegexpStr(msgSpec.Parts)
	// log.Printf("regexp str: '%s'", regexpStr)
	msgRegexp = regexp.MustCompile(regexpStr)
	return regexpStr, msgRegexp, nil
}

func buildSegmentListStr(segmentIDs []string) string {
	var buf bytes.Buffer
	for _, id := range segmentIDs {
		buf.WriteString(fmt.Sprintf("%s:", id))
	}
	return buf.String()
}

func (v *MessageValidator) Validate(message msg.Message) (isValid bool, err error) {
	panic("Not implemented")
}

// Validate a list of segment names as they occur in a message
func (v *MessageValidator) ValidateSegmentList(segmentIDs []string) (isValid bool, err error) {
	if len(segmentIDs) == 0 {
		return false, errors.New("No segments")
	}
	match := v.segmentValidationRegexp.FindStringSubmatch(buildSegmentListStr(segmentIDs))
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
