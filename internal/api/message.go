package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"time"
	"whatsapp-client/internal/model"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"
)

type (
	MessagesReq struct {
		model.Pagination
	}
	MessagesRes struct {
		ID        uint      `json:"id,omitempty"`
		From      string    `json:"from,omitempty"`
		To        string    `json:"to,omitempty"`
		Type      int       `json:"type,omitempty"`
		Text      string    `json:"text,omitempty"`
		FileName  string    `json:"fileName,omitempty"`
		CreatedAt time.Time `json:"createdAt,omitempty"`
	}
)

type SendReq struct {
	JID   string               `form:"jid"`
	Phone string               `form:"phone"`
	Type  int                  `form:"type"`
	Text  string               `form:"text"`
	File  multipart.FileHeader `form:"file"`
}

func MessageSend(c *gin.Context) {
	var req SendReq
	if err := c.Bind(&req); err != nil {
		return
	}

	client, _ := whatsapp.GetClient(req.JID)

	if req.File.Size == 0 {
		client.SendTextMessage(req.Phone, req.Text)

		model.DB.Save(&model.WhatsappSendMessage{
			From: req.JID,
			To:   req.Phone,
			Type: req.Type,
			Text: req.Text,
		})
	} else {
		bytes := FormFileData(req.File)

		if req.Type == 1 {
			client.SendImageMessage(req.Phone, bytes, req.Text)
		} else if req.Type == 2 {
			client.SendDocumentMessage(req.Phone, bytes, req.Text)
		}
		model.DB.Save(&model.WhatsappSendMessage{
			From:     req.JID,
			To:       req.Phone,
			Type:     req.Type,
			Text:     req.Text,
			FileName: req.File.Filename,
		})
	}
	c.JSON(0, nil)
}

func MessageQuery(c *gin.Context) {
	var req MessagesReq
	if err := c.Bind(&req); err != nil {
		return
	}

	var list []MessagesRes

	var total int64
	model.DB.Model(&model.WhatsappSendMessage{}).
		Scopes(model.Paginate(req.Pagination)).
		Count(&total).
		Order("id desc").
		Find(&list)

	c.JSON(0, gin.H{
		"total": total,
		"list":  list,
	})
}

func FormFileData(f multipart.FileHeader) []byte {
	file, err := f.Open()
	defer utils.Close(file)
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return bytes
}
