package impl

import (
	"context"
	"kagari/entity"
	"kagari/service"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

func NewArticleAccessor(ctx context.Context, driver neo4j.DriverWithContext) service.ArticleAccessor {

	return &ArticleAccessor{driver}
}

type ArticleAccessor struct{ driver neo4j.DriverWithContext }

func (aa *ArticleAccessor) GetOne(ctx context.Context, id string) (*entity.Article, error) {
	result, err := neo4j.ExecuteQuery(ctx, aa.driver, "MATCH (a:Article {id: $id}) RETURN a", map[string]any{
		"id": id,
	}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	first := result.Records[0]
	if first == nil {
		return nil, nil
	}
	return (&entity.Article{}).FromMap(first.AsMap()["a"].(dbtype.Node).GetProperties()), nil
}
