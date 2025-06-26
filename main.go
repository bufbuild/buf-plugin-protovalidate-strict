// Copyright 2024-2025 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/check/checkutil"
	"buf.build/go/bufplugin/descriptor"
	"buf.build/go/protovalidate"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
)

func main() {
	check.Main(
		&check.Spec{
			Rules: []*check.RuleSpec{
				{
					ID:      "PROTOVALIDATE_STRICT",
					Default: true,
					Purpose: "Checks that Protovalidate annotations on all existing messages and fields are unchanged.",
					Type:    check.RuleTypeBreaking,
					Handler: checkutil.NewFilePairRuleHandler(handleProtovalidateStrict),
				},
			},
		},
	)
}

func handleProtovalidateStrict(
	_ context.Context,
	responseWriter check.ResponseWriter,
	_ check.Request,
	fileDescriptor descriptor.FileDescriptor,
	againstFileDescriptor descriptor.FileDescriptor,
) error {
	return compareProtovalidateRules(
		responseWriter,
		fileDescriptor.ProtoreflectFileDescriptor().Messages(),
		againstFileDescriptor.ProtoreflectFileDescriptor().Messages(),
	)
}

func compareProtovalidateRules(
	responseWriter check.ResponseWriter,
	messages protoreflect.MessageDescriptors,
	againstMessages protoreflect.MessageDescriptors,
) error {
	for i := range messages.Len() {
		message := messages.Get(i)
		againstMessage := againstMessages.ByName(message.Name())
		if againstMessage == nil {
			continue
		}
		if err := compareProtovalidateMessageRules(responseWriter, message, againstMessage); err != nil {
			return err
		}
		if err := compareProtovalidateRules(responseWriter, message.Messages(), againstMessage.Messages()); err != nil {
			return err
		}
	}
	return nil
}

func compareProtovalidateMessageRules(
	responseWriter check.ResponseWriter,
	message protoreflect.MessageDescriptor,
	againstMessage protoreflect.MessageDescriptor,
) error {
	messageRules, err := protovalidate.ResolveMessageRules(message)
	if err != nil {
		return err
	}
	againstMessageRules, err := protovalidate.ResolveMessageRules(againstMessage)
	if err != nil {
		return err
	}
	if cmp.Diff(messageRules, againstMessageRules, protocmp.Transform()) != "" {
		responseWriter.AddAnnotation(
			check.WithMessagef(
				"Protovalidate message rules on message %q have changed.",
				message.Name(),
			),
			// Note: the annotation will use the message descriptor's location. However, in an
			// ideal world, we would get the location of the rules set on the message. In order
			// to do that, we would need to range across the field descriptors of
			// messageRules.ProtoReflect().Descriptor() and check for the [protoreflect.SourceLocation]
			// for each [protoreflect.SourcePath] based on the field number from the
			// [protoreflect.FileDescriptor.SourceLocations].
			check.WithDescriptor(message),
			check.WithAgainstDescriptor(againstMessage),
		)
	}
	return compareProtovalidateFieldRules(responseWriter, message.Fields(), againstMessage.Fields())
}

func compareProtovalidateFieldRules(
	responseWriter check.ResponseWriter,
	fields protoreflect.FieldDescriptors,
	againstFields protoreflect.FieldDescriptors,
) error {
	for i := range fields.Len() {
		field := fields.Get(i)
		againstField := againstFields.ByNumber(field.Number())
		if againstField == nil {
			continue
		}
		fieldRules, err := protovalidate.ResolveFieldRules(field)
		if err != nil {
			return err
		}
		againstFieldRules, err := protovalidate.ResolveFieldRules(againstField)
		if err != nil {
			return err
		}
		if cmp.Diff(fieldRules, againstFieldRules, protocmp.Transform()) != "" {
			// Note: this has the same source location limitation as message rules above.
			responseWriter.AddAnnotation(
				check.WithMessagef(
					"Protovalidate field rules on field %q have changed.",
					field.FullName(),
				),
				check.WithDescriptor(field),
				check.WithAgainstDescriptor(againstField),
			)
		}
	}
	return nil
}
