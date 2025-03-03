package transformations

func TransformIdTitle[T any](models []T, idFunc func(T) interface{}, titleFunc func(T) string) []map[string]interface{} {
	transformed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		transformed[i] = map[string]interface{}{
			"id":    idFunc(model),
			"title": titleFunc(model),
		}
	}
	return transformed
}

func TransformLoans[T any](models []T, id, loanA, loanT, loanC, loanU, agentF, agentL, memberF, memberL func(T) interface{}) []map[string]interface{} {
	transformed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		transformed[i] = map[string]interface{}{
			"id":      id(model),
			"amount":  loanA(model),
			"term":    loanT(model),
			"created": loanC(model),
			"updated": loanU(model),
			"agentF":  agentF(model),
			"agentL":  agentL(model),
			"memberF": memberF(model),
			"memberL": memberL(model),
		}
	}
	return transformed
}

func TransformSeven[T any](models []T, one, two, three, four, five, six func(T) interface{}, seven func(T) []map[string]interface{}) []map[string]interface{} {
	transformed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		transformed[i] = map[string]interface{}{
			"one":   one(model),
			"two":   two(model),
			"three": three(model),
			"four":  four(model),
			"five":  five(model),
			"six":   six(model),
			"seven": seven(model),
		}
	}
	return transformed
}

func TransformIdSlugUpdatedAt[T any](models []T, idFunc func(T) interface{}, slugFunc func(T) string, updatedAtFunc func(T) interface{}) []map[string]interface{} {
	transformed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		transformed[i] = map[string]interface{}{
			"id":         idFunc(model),
			"slug":       slugFunc(model),
			"updated_at": updatedAtFunc(model),
		}
	}
	return transformed
}

func TransformIdTitleImage[T any](models []T, idFunc func(T) interface{}, titleFunc func(T) string, imageFunc func(T) string) []map[string]interface{} {
	transformed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		transformed[i] = map[string]interface{}{
			"id":    idFunc(model),
			"title": titleFunc(model),
			"image": imageFunc(model),
		}
	}
	return transformed
}

func TransformIdTitleImageSlug[T any](models []T, idFunc func(T) interface{}, titleFunc func(T) string, imageFunc func(T) string, slugFunc func(T) string) []map[string]interface{} {
	transformed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		transformed[i] = map[string]interface{}{
			"id":    idFunc(model),
			"title": titleFunc(model),
			"image": imageFunc(model),
			"slug":  slugFunc(model),
		}
	}
	return transformed
}

func TransformModelBlogPosts[T any](
	models []T,
	idFunc func(T) interface{},
	titleFunc func(T) string,
	snippetFunc func(T) string,
	imageFunc func(T) string,
	slugFunc func(T) string,
	authorFunc func(T) string,
	categoriesFunc func(T) []map[string]interface{},
	tagsFunc func(T) []map[string]interface{},
	createdAtFunc func(T) interface{},
) []map[string]interface{} {
	transformed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		transformed[i] = map[string]interface{}{
			"id":         idFunc(model),
			"title":      titleFunc(model),
			"snippet":    snippetFunc(model),
			"image":      imageFunc(model),
			"slug":       slugFunc(model),
			"author":     authorFunc(model),
			"categories": categoriesFunc(model),
			"tags":       tagsFunc(model),
			"created_at": createdAtFunc(model),
		}
	}
	return transformed
}

func TransformModelBlogPost[T any](
	model T,
	idFunc func(T) interface{},
	titleFunc func(T) string,
	snippetFunc func(T) string,
	contentFunc func(T) string,
	imageFunc func(T) string,
	slugFunc func(T) string,
	authorFunc func(T) string,
	categoriesFunc func(T) []map[string]interface{},
	tagsFunc func(T) []map[string]interface{},
	createdAtFunc func(T) interface{},
) map[string]interface{} {
	return map[string]interface{}{
		"id":         idFunc(model),
		"title":      titleFunc(model),
		"snippet":    snippetFunc(model),
		"content":    contentFunc(model),
		"image":      imageFunc(model),
		"slug":       slugFunc(model),
		"author":     authorFunc(model),
		"categories": categoriesFunc(model),
		"tags":       tagsFunc(model),
		"created_at": createdAtFunc(model),
	}
}

func TransformModelCategory[T any](
	model T,
	idFunc func(T) interface{},
	titleFunc func(T) string,
	snippetFunc func(T) string,
	slugFunc func(T) string,
) map[string]interface{} {
	return map[string]interface{}{
		"id":      idFunc(model),
		"title":   titleFunc(model),
		"snippet": snippetFunc(model),
		"slug":    slugFunc(model),
	}
}

func TransformRelatedBlogs[T any](models []T,
	idFunc func(*T) interface{},
	titleFunc func(*T) string,
	snippetFunc func(*T) string,
	contentFunc func(*T) string,
	imageFunc func(*T) string,
	slugFunc func(*T) string,
	authorFunc func(*T) string,
	categoriesFunc func(*T) []map[string]interface{},
	tagsFunc func(*T) []map[string]interface{},
	createdAtFunc func(*T) interface{}) []map[string]interface{} {

	var transformed []map[string]interface{}
	for _, model := range models {
		transformedModel := map[string]interface{}{
			"id":         idFunc(&model),
			"title":      titleFunc(&model),
			"snippet":    snippetFunc(&model),
			"content":    contentFunc(&model),
			"image":      imageFunc(&model),
			"slug":       slugFunc(&model),
			"author":     authorFunc(&model),
			"categories": categoriesFunc(&model),
			"tags":       tagsFunc(&model),
			"created_at": createdAtFunc(&model),
		}
		transformed = append(transformed, transformedModel)
	}
	return transformed
}
