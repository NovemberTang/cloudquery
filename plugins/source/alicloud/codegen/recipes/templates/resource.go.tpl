// Code generated by codegen; DO NOT EDIT.

package {{.Service}}

import (
	"github.com/cloudquery/plugin-sdk/schema"
    //"github.com/cloudquery/cloudquery/plugins/source/alicloud/client"
)

func {{.SubService | ToCamel}}() *schema.Table {
    return &schema.Table{{template "table.go.tpl" .Table}}
}