package validation

import (
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/util"
	"log"
)

// Builds a nested message
// Used in conjunction with SegSeqValidator
// Shares group stack with SegSeqValidator
type NestedMsgBuilder struct {
	groupStack *util.Stack

	// Nested (hierarchical) message under construction
	nestedMsg *msg.NestedMsg
}

func (s *NestedMsgBuilder) currentGroupContext() *SegSeqGroupContext {
	result := s.groupStack.Peek().(*SegSeqGroupContext)
	return result
}

func (s *NestedMsgBuilder) isAtTopLevel() bool {
	return s.groupStack.Len() < 2
}

func (b *NestedMsgBuilder) String() string {
	return fmt.Sprintf("NestedMsgBuilder msg: %s groupStack: %d elements",
		b.nestedMsg.Name, b.groupStack.Len())
}

func (b *NestedMsgBuilder) AddSegment(segment *msg.Segment) {
	log.Printf("BUILD: AddSegment %s", segment.Id())
	gc := b.currentGroupContext()
	if b.isAtTopLevel() {
		b.nestedMsg.AppendPart(msg.NewRepeatSegment(segment))
	} else {
		gc.repeatSegGroup.GetLast().AppendSeg(segment)
	}
}

func (b *NestedMsgBuilder) RepeatSegment(segment *msg.Segment) {
	log.Printf("BUILD: RepeatSegment %s NOT IMPLEMENTED", segment.Id())
	//panic("Not implemented")
	var repeatSeg *msg.RepeatSegment
	if b.isAtTopLevel() {
		repeatSeg = b.nestedMsg.GetLastPart().(*msg.RepeatSegment)
		/*if !ok {
			panic(fmt.Sprintf(
				"Incorrect type %T; expected *RepeatSegment", repeatSeg))
		}*/
	} else {
		//repeatSeg = b.currentGroupContext().repeatSegGroup.GetLast().(*msg.RepeatSegment)
	}
	repeatSeg.AddSegment(segment)
}

func (b *NestedMsgBuilder) AddSegmentGroup(name string) *msg.RepeatSegmentGroup {
	log.Printf("BUILD: AddSegmentGroup %s", name)
	newRepeatMsgPart := []msg.RepeatMsgPart{}
	repeatSegGroup := msg.NewRepeatSegmentGroup(
		msg.NewSegmentGroup(name, newRepeatMsgPart))
	gc := b.currentGroupContext()
	if b.isAtTopLevel() {
		log.Printf("Appending segment group %s to nested msg %s (%d parts)",
			repeatSegGroup.Id(), b.nestedMsg.Name, b.nestedMsg.Count())
		b.nestedMsg.AppendPart(repeatSegGroup)
	} else {
		log.Printf("Appending segment group %s to %d parts of %s",
			repeatSegGroup.Id(), gc.repeatSegGroup.Count(), gc.repeatSegGroup.Id())
		gc.repeatSegGroup.AppendSegGroupToLast(repeatSegGroup)
	}
	log.Printf("### msg parts count %d", b.nestedMsg.Count())
	return repeatSegGroup
}

func (b *NestedMsgBuilder) RepeatSegmentGroup(segmentGroup *msg.SegmentGroup) {
	log.Printf("BUILD: RepeatSegmentGroup %s NOT IMPLEMENTED", segmentGroup.Id())
}

func NewNestedMsgBuilder(msgName string, groupStack *util.Stack) *NestedMsgBuilder {
	return &NestedMsgBuilder{
		groupStack: groupStack,
		nestedMsg:  msg.NewNestedMsg(msgName, []msg.RepeatMsgPart{}),
	}
}
