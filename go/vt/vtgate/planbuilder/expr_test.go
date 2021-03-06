/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package planbuilder

import (
	"testing"

	"vitess.io/vitess/go/vt/sqlparser"
)

func TestValEqual(t *testing.T) {
	c1 := &column{}
	c2 := &column{}
	testcases := []struct {
		in1, in2 sqlparser.Expr
		out      bool
	}{{
		in1: &sqlparser.ColName{Metadata: c1, Name: sqlparser.NewColIdent("c1")},
		in2: &sqlparser.ColName{Metadata: c1, Name: sqlparser.NewColIdent("c1")},
		out: true,
	}, {
		// Objects that have the same name need not be the same because
		// they might have appeared in different scopes and could have
		// resolved to different columns.
		in1: &sqlparser.ColName{Metadata: c1, Name: sqlparser.NewColIdent("c1")},
		in2: &sqlparser.ColName{Metadata: c2, Name: sqlparser.NewColIdent("c1")},
		out: false,
	}, {
		in1: newValArg(":aa"),
		in2: &sqlparser.ColName{Metadata: c1, Name: sqlparser.NewColIdent("c1")},
		out: false,
	}, {
		in1: newValArg(":aa"),
		in2: newValArg(":aa"),
		out: true,
	}, {
		in1: newValArg(":aa"),
		in2: newValArg(":bb"),
	}, {
		in1: newStrLiteral("aa"),
		in2: newStrLiteral("aa"),
		out: true,
	}, {
		in1: newStrLiteral("11"),
		in2: newHexLiteral("3131"),
		out: true,
	}, {
		in1: newHexLiteral("3131"),
		in2: newStrLiteral("11"),
		out: true,
	}, {
		in1: newHexLiteral("3131"),
		in2: newHexLiteral("3131"),
		out: true,
	}, {
		in1: newHexLiteral("3131"),
		in2: newHexLiteral("3132"),
		out: false,
	}, {
		in1: newHexLiteral("313"),
		in2: newHexLiteral("3132"),
		out: false,
	}, {
		in1: newHexLiteral("3132"),
		in2: newHexLiteral("313"),
		out: false,
	}, {
		in1: newIntLiteral("313"),
		in2: newHexLiteral("3132"),
		out: false,
	}, {
		in1: newHexLiteral("3132"),
		in2: newIntLiteral("313"),
		out: false,
	}, {
		in1: newIntLiteral("313"),
		in2: newIntLiteral("313"),
		out: true,
	}, {
		in1: newIntLiteral("313"),
		in2: newIntLiteral("314"),
		out: false,
	}}
	for _, tc := range testcases {
		out := valEqual(tc.in1, tc.in2)
		if out != tc.out {
			t.Errorf("valEqual(%#v, %#v): %v, want %v", tc.in1, tc.in2, out, tc.out)
		}
	}
}

func newStrLiteral(in string) *sqlparser.Literal {
	return sqlparser.NewStrLiteral([]byte(in))
}

func newIntLiteral(in string) *sqlparser.Literal {
	return sqlparser.NewIntLiteral([]byte(in))
}

func newHexLiteral(in string) *sqlparser.Literal {
	return sqlparser.NewHexLiteral([]byte(in))
}

func newValArg(in string) sqlparser.Expr {
	return sqlparser.NewArgument([]byte(in))
}
