package media

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	// WxOpenComment 打开评论
	WxOpenComment = "cgi-bin/comment/open"
	// wxCloseComment 关闭评论
	wxCloseComment = "cgi-bin/comment/close"
)

// JSONContentType http method=POST, Content-Type
const JSONContentType = "application/json; charset=utf-8"

// changeComment 修改评论的操作
func changeComment(host, accessToken, action string, msgdataid, index uint32) (*Response, error) {
	JSON := fmt.Sprintf(`{"msg_data_id":"%d", "index":"%d"}`, msgdataid, index)
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, action, accessToken)
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := readResponse(res)
	if err != nil {
		return nil, err
	}
	var resp Response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// OpenComment 打开评论功能
func OpenComment(host, accessToken string, msgdataid, index uint32) (*Response, error) {
	return changeComment(host, accessToken, WxOpenComment, msgdataid, index)
}

// CloseComment 关闭评论功能
func CloseComment(host, accessToken string, msgdataid, index uint32) (*Response, error) {
	return changeComment(host, accessToken, wxCloseComment, msgdataid, index)
}

// CommentResponse 请求回复数据时，服务器的响应结构
type CommentResponse struct {
	ErrCode int        `json:",errcode,omitempty"`
	ErrMsg  string     `json:"errmsg,omitempty"`
	Total   int        `json:"total,omitempty"`
	Comment []*Comment `json:"comment,omitempty"`
}

// Comment 回应中评论
type Comment struct {
	UserCommentID int           `json:"user_comment_id,omitempty"`
	OpenID        int           `json:"openid,omitempty"`
	CreateTime    int           `json:"create_time,omitempty"`
	Content       string        `json:"content,omitempty"`
	CommentType   int           `json:"comment_type,omitempty"`
	Reply         *CommentReply `json:"reply,omitempty"`
}

// CommentReply 回复评论的内容
type CommentReply struct {
	Content    string `json:"content,omitempty"`
	CreateTime int    `json:"create_time,omitempty"`
}

// WxCommentList 请求评论列表的路径
const WxCommentList = `cgi-bin/comment/list`

// GetCommentList 获取评论列表
func GetCommentList(host, accessToken string, msgdataid, index, begin, count, typ uint32) (*CommentResponse, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, WxCommentList, accessToken)
	JSON := fmt.Sprintf(`{"msg_data_id":%d,"index":%d,"begin":%d,"count":%d,"type":%d}`, msgdataid, index, begin, count, typ)
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := readResponse(res)
	if err != nil {
		return nil, err
	}
	var cres CommentResponse
	if err = json.Unmarshal(b, &cres); err != nil {
		return nil, err
	}
	return &cres, nil
}

const (
	// WxMarkElect 将评论设置为精选
	WxMarkElect = `cgi-bin/comment/markelect`
	// WxUnMarkElect 取消评论精选
	WxUnMarkElect = `cgi-bin/comment/unmarkelect`
)

// ChangeElect 修改评论
func ChangeElect(host, accessToken, path string, msgdataid, index, usercommentid uint32) (*Response, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, path, accessToken)
	JSON := fmt.Sprintf(`{"msg_data_id":%d, "index":%d, "user_comment_id":%d}`, msgdataid, index, usercommentid)
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := readResponse(res)
	if err != nil {
		return nil, err
	}
	var resp Response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MarkElect 标记评论为精选
func MarkElect(host, accessToken string, msgdataid, index, usercommentid uint32) (*Response, error) {
	return ChangeElect(host, accessToken, WxMarkElect, msgdataid, index, usercommentid)
}

// UnMarkElect 撤销评论精选
func UnMarkElect(host, accessToken string, msgdataid, index, usercommentid uint32) (*Response, error) {
	return ChangeElect(host, accessToken, WxUnMarkElect, msgdataid, index, usercommentid)
}

// WxDelComment 删除评论
const WxDelComment = `cgi-bin/comment/delete`

// DeleteComment 删除评论
func DeleteComment(host, accessToken string, msgdataid, index, usercommentid uint32) (*Response, error) {
	return ChangeElect(host, accessToken, WxDelComment, msgdataid, index, usercommentid)
}

// WxReplyComment 回复评论的路径
const WxReplyComment = `cgi-bin/comment/reply/add`

// ReplyComment 回复评论
func ReplyComment(host, accessToken string, msgdataid, index, usercommentid uint32, content string) (*Response, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, WxReplyComment, accessToken)
	JSON := fmt.Sprintf(`{"msg_data_id":%d, "index":%d, "user_comment_id":%d, "content":%s}`, msgdataid, index, usercommentid, content)
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := readResponse(res)
	if err != nil {
		return nil, err
	}
	var resp Response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// WxDelCommentReply 删除评论回复的api
const WxDelCommentReply = `cgi-bin/comment/reply/delete`

// DeleteCommentReply 删除回复评论的内容
func DeleteCommentReply(host, accessToken string, msgdataid, index, usercommentid uint32) (*Response, error) {
	return ChangeElect(host, accessToken, WxDelCommentReply, msgdataid, index, usercommentid)
}
