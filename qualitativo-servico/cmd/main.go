package main

import (
	"git.sof.intra/siop/framework"
	_ "git.sof.intra/siop/orgaos-servico/client"
	"git.sof.intra/siop/qualitativo-servico/internal"
	_ "git.sof.intra/siop/qualitativo-servico/internal/gerados"
	"git.sof.intra/siop/qualitativo-servico/internal/graphql"
)

func main() {
	framework.Run(
		&internal.Service{},
		framework.WithGrpc(),
		framework.WithGraphQL(
			"/modulo/qualitativo",
			graphql.NewResolver,
			framework.WithGraphQLAuthBlacklist("IntrospectionQuery", "__schema"),
		))
}
