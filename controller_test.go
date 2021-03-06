// Copyright 2019 asana Author. All Rights Reserved.
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

package asana

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/goasana/asana/context"
)

func TestGetInt(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", "40")
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetInt("age")
	if val != 40 {
		t.Errorf("TestGetInt expect 40,get %T,%v", val, val)
	}
}

func TestGetInt8(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", "40")
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetInt8("age")
	if val != 40 {
		t.Errorf("TestGetInt8 expect 40,get %T,%v", val, val)
	}
	//ResponseWriter: int8
}

func TestGetInt16(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", "40")
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetInt16("age")
	if val != 40 {
		t.Errorf("TestGetInt16 expect 40,get %T,%v", val, val)
	}
}

func TestGetInt32(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", "40")
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetInt32("age")
	if val != 40 {
		t.Errorf("TestGetInt32 expect 40,get %T,%v", val, val)
	}
}

func TestGetInt64(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", "40")
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetInt64("age")
	if val != 40 {
		t.Errorf("TestGeetInt64 expect 40,get %T,%v", val, val)
	}
}

func TestGetUint8(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", strconv.FormatUint(math.MaxUint8, 10))
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetUint8("age")
	if val != math.MaxUint8 {
		t.Errorf("TestGetUint8 expect %v,get %T,%v", math.MaxUint8, val, val)
	}
}

func TestGetUint16(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", strconv.FormatUint(math.MaxUint16, 10))
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetUint16("age")
	if val != math.MaxUint16 {
		t.Errorf("TestGetUint16 expect %v,get %T,%v", math.MaxUint16, val, val)
	}
}

func TestGetUint32(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", strconv.FormatUint(math.MaxUint32, 10))
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetUint32("age")
	if val != math.MaxUint32 {
		t.Errorf("TestGetUint32 expect %v,get %T,%v", math.MaxUint32, val, val)
	}
}

func TestGetUint64(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetParam("age", strconv.FormatUint(math.MaxUint64, 10))
	ctrlr := Controller{Context: ctx}
	val, _ := ctrlr.GetUint64("age")
	if val != math.MaxUint64 {
		t.Errorf("TestGetUint64 expect %v,get %T,%v", uint64(math.MaxUint64), val, val)
	}
}

func TestAdditionalViewPaths(t *testing.T) {
	dir1 := "_asanaTmp"
	dir2 := "_asanaTmp2"
	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)

	dir1file := "file1.tpl"
	dir2file := "file2.tpl"

	genFile := func(dir string, name string, content string) {
		_ = os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0777)
		if f, err := os.Create(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		} else {
			defer f.Close()
			_, _ = f.WriteString(content)
			_ = f.Close()
		}

	}
	genFile(dir1, dir1file, `<div>{{.Content}}</div>`)
	genFile(dir2, dir2file, `<html>{{.Content}}</html>`)

	_ = AddViewPath(dir1)
	_ = AddViewPath(dir2)

	ctrl := Controller{
		Context: context.NewContext(),
		TplName:  "file1.tpl",
		ViewPath: dir1,
	}
	ctrl.Response().SetData(map[interface{}]interface{}{
		"Content": "value2",
	})
	if result, err := ctrl.RenderString(); err != nil {
		t.Fatal(err)
	} else {
		if result != "<div>value2</div>" {
			t.Fatalf("TestAdditionalViewPaths expect %s got %s", "<div>value2</div>", result)
		}
	}

	func() {
		ctrl.TplName = "file2.tpl"
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("TestAdditionalViewPaths expected error")
			}
		}()
		_, _ = ctrl.RenderString()
	}()

	ctrl.TplName = "file2.tpl"
	ctrl.ViewPath = dir2
	_, _ = ctrl.RenderString()
}
