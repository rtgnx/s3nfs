package nfsd

import (
	"reflect"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/rtgnx/s3nfs/internal/fs/noop"
)

func TestMapperFS_Resolve(t *testing.T) {
	type fields struct {
		NoOP   noop.NoOP
		shares map[string]billy.Filesystem
	}
	type args struct {
		filename string
	}
	noopFS := noop.NewNoOP()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLfs billy.Filesystem
		wantP   string
		wantErr bool
	}{
		{
			name: "resolve valid share",
			fields: fields{shares: map[string]billy.Filesystem{
				"/fr-par": noopFS,
			}},

			args:    args{"/fr-par/test.txt"},
			wantLfs: noopFS,
			wantP:   "/test.txt",
			wantErr: false,
		},
		{
			name: "resolve valid share, no file",
			fields: fields{shares: map[string]billy.Filesystem{
				"/fr-par": noopFS,
			}},

			args:    args{"/fr-par/"},
			wantLfs: noopFS,
			wantP:   "/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MapperFS{
				NoOP:   tt.fields.NoOP,
				shares: tt.fields.shares,
			}
			gotLfs, gotP, err := m.Resolve(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapperFS.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLfs, tt.wantLfs) {
				t.Errorf("MapperFS.Resolve() gotLfs = %v, want %v", gotLfs, tt.wantLfs)
			}
			if gotP != tt.wantP {
				t.Errorf("MapperFS.Resolve() gotP = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}
