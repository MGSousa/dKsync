package dksync

import (
	"testing"
)

type fields struct {
	Tag   		string
	Debug 		bool
	toUpdate  	Args
}

func TestApiSync_ChangeProcessor(t *testing.T) {
	tests := []struct {
		name   		string
		fields  	fields
	}{
		{
			name: "ChangeProcessors",
			fields: fields {
				Tag: "nodesync",
				Debug: true,
				toUpdate: Args {
					Field: "Processors",
					Value: "{\"mail\":{\"subject\":\"dKron Watcher\", \"debug\":\"true\"}}",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sync := &ApiSync {
				Tag:   tt.fields.Tag,
				Debug: tt.fields.Debug,
				Args: Args {
					Field: tt.fields.toUpdate.Field,
					Value: tt.fields.toUpdate.Value,
					Del: tt.fields.toUpdate.Del,
				},
			}
			sync.Processor()
		})
	}
}

func TestApiSync_RemoveProcessor(t *testing.T) {
	tests := []struct {
		name   		string
		fields  	fields
	}{
		{
			name: "RemoveProcessors",
			fields: fields {
				Tag: "nodesync",
				Debug: true,
				toUpdate: Args {
					Field: "Processors",
					Value: "mail",
					Del:   true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sync := &ApiSync {
				Tag:   tt.fields.Tag,
				Debug: tt.fields.Debug,
				Args: Args {
					Field: tt.fields.toUpdate.Field,
					Value: tt.fields.toUpdate.Value,
					Del: tt.fields.toUpdate.Del,
				},
			}
			sync.Processor()
		})
	}
}

func TestApiSync_OtherOptions(t *testing.T) {
	tests := []struct {
		name   		string
		fields  	fields
	}{
		{
			name: "ChangeRetries",
			fields: fields {
				Tag: "nodesync",
				Debug: true,
				toUpdate: Args {
					Field: "Retries",
					Value: "2",
				},
			},
		},
		{
			name: "ChangeOwner",
			fields: fields {
				Tag: "nodesync",
				Debug: true,
				toUpdate: Args {
					Field: "Owner",
					Value: "KK-Test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sync := &ApiSync {
				Tag:   tt.fields.Tag,
				Debug: tt.fields.Debug,
				Args: Args {
					Field: tt.fields.toUpdate.Field,
					Value: tt.fields.toUpdate.Value,
					Del: tt.fields.toUpdate.Del,
				},
			}
			sync.Processor()
		})
	}
}
