package dao

func ToArticle(req *CreateArticleRequest) *Article {
	return &Article{
		Title: req.Title,
		Content: req.Content,
		Is_active: req.Is_active,
	}
}

func ToUpdatedArticle(req *UpdateArticleRequest) *Article {
	return &Article{
		ID: req.ID,
		Title: req.Title,
		Content: req.Content,
		Is_active: req.Is_active,
	}
}
