// Code generated by codegen; DO NOT EDIT.

package {{.Service}}

import (
	"context"
	"fmt"
	"testing"
	"google.golang.org/grpc"
	"github.com/cloudquery/plugin-sdk/faker"
	"github.com/cloudquery/plugins/source/gcp/client"
	{{if .ProtobufImport}}
  pb "{{.ProtobufImport}}"
  {{end}}{{ if .RelationsTestData.ProtobufImport }}
  "{{.RelationsTestData.ProtobufImport}}"
  {{end}}
)


func create{{.SubService | ToCamel}}(gsrv *grpc.Server) error {
	fakeServer := &fake{{.SubService | ToCamel}}Server{}
	pb.{{.RegisterServerName}}(gsrv, fakeServer){{ if .RelationsTestData.RegisterServerName }}
	fakeRelationsServer := &fake{{.SubService | ToCamel}}RelationsServer{}
	{{.RelationsTestData.RegisterServerName}}(gsrv, fakeRelationsServer)
	{{ end }}
	return nil
}

type fake{{.SubService | ToCamel}}Server struct {
	pb.{{.UnimplementedServerName}}
}

func (f *fake{{.SubService | ToCamel}}Server) {{.ListFunctionName}}(context.Context, *pb.{{.RequestStructName}}) (*pb.{{.ResponseStructName}}, error) {
	resp := pb.{{.ResponseStructName}}{}
	if err := faker.FakeObject(&resp); err != nil {
		return nil, fmt.Errorf("failed to fake data: %w", err)
	}
	resp.NextPageToken = ""
	return &resp, nil
}

func Test{{.SubService | ToCamel}}(t *testing.T) {
	client.MockTestGrpcHelper(t, {{.SubService | ToCamel}}(), create{{.SubService | ToCamel}}, client.TestOptions{})
}


{{ if .RelationsTestData.UnimplementedServerName }}
type fake{{.SubService | ToCamel}}RelationsServer struct {
	{{.RelationsTestData.UnimplementedServerName}}
}

{{ range .RelationsTestData.ListFunctions }}
func (f *fake{{$.SubService | ToCamel}}RelationsServer) {{.Signature}} {
	resp := {{.ResponseStructName}}{}
	if err := faker.FakeObject(&resp); err != nil {
		return nil, fmt.Errorf("failed to fake data: %w", err)
	}
	resp.NextPageToken = ""
	return &resp, nil
}
{{ end }}
{{ end }}