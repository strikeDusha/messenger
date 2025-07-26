package handlers

import (
	"messenger/database"
	"sort"

	"github.com/gin-gonic/gin"
)

func (d *DB) History(c *gin.Context) {
	re := c.Param("id") // id of receiver
	se, _ := c.Get("username")
	ss, _ := se.(string)
	var msgs []*database.MsgMini
	mine, err := d.s.GetChat(ss, re) // []*MsgMini { text,time} of mine messages
	if err != nil {
		c.AbortWithStatus(500)
	}
	his, err := d.s.GetChat(re, ss)
	if err != nil {
		c.AbortWithStatus(500)
	}
	msgs = append(msgs, mine...)
	msgs = append(msgs, his...)
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Time.Before(msgs[j].Time)
	})
	c.JSON(200, msgs)

}
