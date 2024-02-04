package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_gen(t *testing.T) {
	buf := &bytes.Buffer{}
	fmt.Println(os.Getwd())
	err := gen(buf, "testdata/user.go")
	require.NoError(t, err)
	require.Equal(t, buf.String(), "package testdata\n\nimport (\n    sqlx \"database/sql\"\n    \"myProject/orm\"\n)\n\n\nconst (\n    UserName = \"Name\"\n    UserAge = \"Age\"\n    UserNickName = \"NickName\"\n    UserPicture = \"Picture\"\n)\n\nconst (\n    UserDetailAddress = \"Address\"\n)\n\n\n\nfunc UserNameLT(val string) orm.Predicate {\n    return orm.C(\"Name\").LT(val)\n}\n\nfunc UserNameGT(val string) orm.Predicate {\n    return orm.C(\"Name\").GT(val)\n}\n\nfunc UserNameEQ(val string) orm.Predicate {\n    return orm.C(\"Name\").EQ(val)\n}\n\n\nfunc UserAgeLT(val *int) orm.Predicate {\n    return orm.C(\"Age\").LT(val)\n}\n\nfunc UserAgeGT(val *int) orm.Predicate {\n    return orm.C(\"Age\").GT(val)\n}\n\nfunc UserAgeEQ(val *int) orm.Predicate {\n    return orm.C(\"Age\").EQ(val)\n}\n\n\nfunc UserNickNameLT(val *sqlx.NullString) orm.Predicate {\n    return orm.C(\"NickName\").LT(val)\n}\n\nfunc UserNickNameGT(val *sqlx.NullString) orm.Predicate {\n    return orm.C(\"NickName\").GT(val)\n}\n\nfunc UserNickNameEQ(val *sqlx.NullString) orm.Predicate {\n    return orm.C(\"NickName\").EQ(val)\n}\n\n\nfunc UserPictureLT(val []byte) orm.Predicate {\n    return orm.C(\"Picture\").LT(val)\n}\n\nfunc UserPictureGT(val []byte) orm.Predicate {\n    return orm.C(\"Picture\").GT(val)\n}\n\nfunc UserPictureEQ(val []byte) orm.Predicate {\n    return orm.C(\"Picture\").EQ(val)\n}\n\n\nfunc UserDetailAddressLT(val string) orm.Predicate {\n    return orm.C(\"Address\").LT(val)\n}\n\nfunc UserDetailAddressGT(val string) orm.Predicate {\n    return orm.C(\"Address\").GT(val)\n}\n\nfunc UserDetailAddressEQ(val string) orm.Predicate {\n    return orm.C(\"Address\").EQ(val)\n}\n\n\n")
}

func Test_GenFile(t *testing.T) {
	f, err := os.Create("testdata/user.gen.go")
	require.NoError(t, err)
	err = gen(f, "testdata/user.go")
	require.NoError(t, err)
}
