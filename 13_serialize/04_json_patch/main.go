package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func main() {
	// json patch
	// 原始 JSON 文档
	original := []byte(`{
        "name": "John",
        "age": 30,
        "city": "New York"
    }`)

	// JSON Patch 文档
	patch := []byte(`[
        { "op": "replace", "path": "/name", "value": "Danny" },
        { "op": "remove", "path": "/age" },
        { "op": "add", "path": "/country", "value": "China" }
    ]`)

	// 创建补丁对象
	patchObj, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		panic(err)
	}

	// 应用补丁
	patched, err := patchObj.Apply(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Patched document: %s\n", patched)
	/*
		{
		    "name": "Danny",
		    "city": "New York",
		    "country": "China"
		}
	*/

	// merge patch
	src := []byte(`
{
    "a": "b",
    "c": {
        "d": "e",
        "f": "g"
    }
}`)
	mergeJson := []byte(`
{
    "a": "z",
    "c": {
        "f": null
    }
}`)
	rsp, err := jsonpatch.MergePatch(src, mergeJson)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Patched document: %s\n", rsp)
	/*
		{
		    "a": "z",
		    "c": {
		        "d": "e"
		    }
		}
	*/

}
