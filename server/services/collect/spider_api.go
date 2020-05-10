package collect

import (
	"bbs-go/common/uploader"
	"errors"
	"github.com/mlogclub/simple"

	"bbs-go/common"
	"bbs-go/common/baiduai"
	"bbs-go/model"
	"bbs-go/services"
)

type SpiderApi struct {
}

func NewSpiderApi() *SpiderApi {
	return &SpiderApi{}
}

func (api *SpiderApi) Publish(article *Article) (articleId int64, err error) {
	if article.Summary == "" {
		article.Summary = common.GetSummary(article.ContentType, article.Content)
	}

	if len(article.Tags) == 0 {
		article.Tags = api.AnalyzeTags(article)
	}

	article.UserId = api.GetUserId(article.UserId, article.Nickname, article.Avatar, article.UserDescription)

	t, err := services.ArticleService.Publish(article.UserId, article.Title, article.Summary, article.Content,
		article.ContentType, article.Tags, article.SourceUrl, true)
	if err == nil {
		articleId = t.Id

		if article.PublishTime > 0 {
			_ = services.ArticleService.UpdateColumn(articleId, "create_time", article.PublishTime)
		}
	}
	return
}

func (api *SpiderApi) PublishComment(comment *Comment) (commentId int64, err error) {
	if len(comment.Content) == 0 {
		err = errors.New("评论内容不能为空")
		return
	}

	comment.UserId = api.GetUserId(comment.UserId, comment.Nickname, comment.Avatar, comment.UserDescription)

	c, err := services.CommentService.Publish(comment.UserId, &model.CreateCommentForm{
		EntityType:  comment.EntityType,
		EntityId:    comment.EntityId,
		Content:     comment.Content,
		ContentType: model.ContentTypeHtml,
	})
	if err == nil {
		commentId = c.Id

		if comment.PublishTime > 0 {
			_ = services.CommentService.UpdateColumn(commentId, "create_time", comment.PublishTime)
		}
	}
	return
}

func (api *SpiderApi) GetUserId(userId int64, nickname, avatar, description string) int64 {
	if userId > 0 {
		return userId
	}
	if simple.IsNotBlank(avatar) {
		avatar, _ = uploader.CopyImage(avatar)
	}
	user := &model.User{
		Nickname:    nickname,
		Avatar:      avatar,
		Description: description,
		Status:      model.StatusOk,
		Type:        model.UserTypeGzh,
		CreateTime:  simple.NowTimestamp(),
		UpdateTime:  simple.NowTimestamp(),
	}
	if err := services.UserService.Create(user); err == nil {
		return user.Id
	}
	return 0
}

func (api *SpiderApi) AnalyzeTags(article *Article) []string {
	var analyzeRet *baiduai.AiAnalyzeRet
	if article.ContentType == model.ContentTypeMarkdown {
		analyzeRet, _ = baiduai.GetAi().AnalyzeMarkdown(article.Title, article.Content)
	} else if article.ContentType == model.ContentTypeHtml {
		analyzeRet, _ = baiduai.GetAi().AnalyzeHtml(article.Title, article.Content)
	}
	var tags []string
	if analyzeRet != nil {
		tags = analyzeRet.Tags
		if article.Summary == "" {
			article.Summary = analyzeRet.Summary
		}
	}
	return tags
}

type Article struct {
	UserId          int64    `json:"userId" form:"userId"` // 发布用户编号
	Nickname        string   `json:"nickname"`
	Avatar          string   `json:"avatar"`
	UserDescription string   `json:"userDescription"`
	Title           string   `json:"title" form:"title"`
	Summary         string   `json:"summary" form:"summary"`
	Content         string   `json:"content" form:"content"`
	ContentType     string   `json:"contentType" form:"contentType"`
	Tags            []string `json:"tags" form:"tags"`
	SourceUrl       string   `json:"sourceUrl" form:"sourceUrl"`
	PublishTime     int64    `json:"publishTime" form:"publishTime"`
}

type Comment struct {
	UserId          int64  `json:"userId" form:"userId"`
	Nickname        string `json:"nickname"`
	Avatar          string `json:"avatar"`
	UserDescription string `json:"userDescription"`
	Content         string `json:"content" form:"content"`
	EntityType      string `json:"entityType" form:"entityType"`
	EntityId        int64  `json:"entityId" form:"entityId"`
	PublishTime     int64  `json:"publishTime" form:"publishTime"`
}
