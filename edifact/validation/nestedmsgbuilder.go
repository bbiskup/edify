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
		b.nestedMsg.AppendPart(msg.NewRepSeg(segment))
	} else {
		gc.repeatSegGroup.GetLast().AppendSeg(segment)
	}
}

func (b *NestedMsgBuilder) RepSeg(segment *msg.Segment) {
	log.Printf("BUILD: RepSeg %s NOT IMPLEMENTED", segment.Id())
	//panic("Not implemented")
	var repeatSeg *msg.RepSeg
	if b.isAtTopLevel() {
		repeatSeg = b.nestedMsg.GetLastPart().(*msg.RepSeg)
		/*if !ok {
			panic(fmt.Sprintf(
				"Incorrect type %T; expected *RepSeg", repeatSeg))
		}*/
	} else {
		//repeatSeg = b.currentGroupContext().repeatSegGroup.GetLast().(*msg.RepSeg)
	}
	repeatSeg.AddSegment(segment)
}

func (b *NestedMsgBuilder) AddSegGrp(name string) *msg.RepeatSegGrp {
	log.Printf("BUILD: AddSegGrp %s", name)
	newRepeatMsgPart := []msg.RepeatMsgPart{}
	repeatSegGroup := msg.NewRepeatSegGrp(
		msg.NewSegGrp(name, newRepeatMsgPart))
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

func (b *NestedMsgBuilder) RepeatSegGrp(segGrp *msg.SegGrp) {
	log.Printf("BUILD: RepeatSegGrp %s NOT IMPLEMENTED", segGrp.Id())
}

func NewNestedMsgBuilder(msgName string, groupStack *util.Stack) *NestedMsgBuilder {
	return &NestedMsgBuilder{
		groupStack: groupStack,
		nestedMsg:  msg.NewNestedMsg(msgName, []msg.RepeatMsgPart{}),
	}
}
