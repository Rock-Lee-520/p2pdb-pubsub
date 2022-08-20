package protocol

import (
	"github.com/Rock-liyi/p2pdb-pubsub/chat/internal/entity"
)

const (
	FileType  = "file"
	ImageType = "image"
	TextType  = "text"
	EventType = "event"

	GoingStatus   = "going"
	FailedStatus  = "failed"
	SucceedStatus = "succeed"
)

type Message struct {
	Id          string       `json:"id,omitempty"`
	Status      string       `json:"status"`
	Type        string       `json:"type"`
	SendTime    int64        `json:"sendTime,omitempty"`
	Content     string       `json:"content"`
	FileSize    int64        `json:"fileSize,omitempty"`
	FileName    string       `json:"fileName,omitempty"`
	ToContactId string       `json:"toContactId,omitempty"`
	FromUser    *entity.User `json:"fromUser,omitempty"`
}
