package validation

import (
	"bytes"
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
	"log"
	"regexp"
	"strconv"
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

func buildMessageSpecPartRegexpStr(msgSpecPart msgspec.MessageSpecPart) string {
	var inner string
	switch msgPart := msgSpecPart.(type) {
	case *msgspec.MessageSpecSegmentPart:
		inner = msgPart.SegmentSpec.Id + ":"
	case *msgspec.MessageSpecSegmentGroupPart:
		groupPartRegexpStrs := []string{}
		for _, groupChild := range msgPart.Children() {
			groupPartRegexpStrs = append(groupPartRegexpStrs, buildMessageSpecPartRegexpStr(groupChild))
		}
		inner = strings.Join(groupPartRegexpStrs, "")
	default:
		panic("Not implemented")
	}

	// Regexp engine allows max. repeat count of 1000, whereas the UNCE
	// EDIFACT spec allows 9999
	var maxRegexpRepeatCountStr string
	maxSpecRepeatCount := msgSpecPart.MaxCount()
	if maxSpecRepeatCount == 9999 {
		maxRegexpRepeatCountStr = "" //  unlimited
	} else if maxSpecRepeatCount > 1000 {
		log.Printf("Clamping max repeat count in regexp (msg part %s)", msgSpecPart.String())
		maxRegexpRepeatCountStr = "1000"
	} else {
		maxRegexpRepeatCountStr = strconv.Itoa(maxSpecRepeatCount)
	}

	minRegexpRepeatCount := msgSpecPart.MinCount()
	if minRegexpRepeatCount == 1 && maxSpecRepeatCount == 1 {
		return fmt.Sprintf("(%s)", inner)
	} else {
		return fmt.Sprintf("(%s){%d,%s}", inner, minRegexpRepeatCount, maxRegexpRepeatCountStr)
	}
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
	match := v.segmentValidationRegexp.FindStringSubmatch(buildSegmentListStr(segmentIDs))
	return match != nil, nil
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
