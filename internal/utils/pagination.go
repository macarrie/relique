package utils

import (
	"github.com/macarrie/relique/internal/consts"
	"github.com/spf13/cobra"
)

func AddPaginationParams(cobraCmd *cobra.Command, pageSizeVar *int) {
	cobraCmd.Flags().IntVarP(pageSizeVar, "limit", "l", consts.DEFAULT_PAGE_SIZE, "Number of elements to show")
}
