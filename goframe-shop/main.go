package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "goframe-shop/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"goframe-shop/internal/cmd"

	_ "goframe-shop/internal/logic"

	_ "goframe-shop/internal/packed"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
