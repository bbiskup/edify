package query

// import (
// 	"errors"
// 	"fmt"
// 	"github.com/bbiskup/edify/edifact/msg"
// )

// // Provides navigation within EDIFACT message
// type Navigator struct {
// }

// func (n *Navigator) navigate(queryStr string, message *msg.NestedMsg) (msgPart msg.NestedMsgPart, err error) {
// 	// return segment, segment group or data element
// 	panic("Not implemented")
// }

// func (n *Navigator) GetSegmentGroup(
// 	queryStr string, message *msg.NestedMsg) (group *msg.SegmentGroup, err error) {
// 	msgPart, err := n.navigate(queryStr, message)
// 	if err != nil {
// 		return nil, err
// 	}
// 	group, ok := msgPart.(*msg.SegmentGroup)
// 	if !ok {
// 		return nil, errors.New(fmt.Sprintf("Unexpected type %T", group))
// 	}
// 	return group, nil
// }

// func (n *Navigator) GetSegment(queryStr string, message *msg.NestedMsg) (*msg.Segment, error) {
// 	msgPart, err := n.navigate(queryStr, message)
// 	if err != nil {
// 		return nil, err
// 	}
// 	segment, ok := msgPart.(*msg.Segment)
// 	if !ok {
// 		return nil, errors.New(fmt.Sprintf("Unexpected type %T", segment))
// 	}
// 	return segment, nil
// }

// func (n *Navigator) GetSegmentDataElem(queryStr string, message *msg.NestedMsg) (*msg.DataElement, error) {
// 	msgPart, err := n.navigate(queryStr, message)
// 	if err != nil {
// 		return nil, err
// 	}
// 	dataElem, ok := msgPart.(*msg.DataElement)
// 	if !ok {
// 		return nil, errors.New(fmt.Sprintf("Unexpected type %T", dataElem))
// 	}
// 	return dataElem, nil
// }

// func NewNavigator() *Navigator {
// 	return &Navigator{}
// }
