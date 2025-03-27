<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [json patch](#json-patch)
  - [JSON Patch](#json-patch)
  - [JSON Merge Patch](#json-merge-patch)
  - [golang 库: github.com/evanphx/json-patch](#golang-%E5%BA%93-githubcomevanphxjson-patch)
    - [功能与特性](#%E5%8A%9F%E8%83%BD%E4%B8%8E%E7%89%B9%E6%80%A7)
  - [JSON patch 在 istio 中的应用](#json-patch-%E5%9C%A8-istio-%E4%B8%AD%E7%9A%84%E5%BA%94%E7%94%A8)
  - [json merge patch 在 argo rollout 中应用](#json-merge-patch-%E5%9C%A8-argo-rollout-%E4%B8%AD%E5%BA%94%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# json patch


JSON Patch是一种用于描述对JSON文档所做的更改的格式（JSON Patch本身也是JSON结构）。
当只更改了一部分时，可用于避免发送整个文档。
可以与HTTP PATCH方法结合使用时，它允许以符合标准的方式对HTTP API进行部分更新,其 MIME 媒体类型为 "application/json-patch+json" 或则 “application/merge-patch+json”。


JSON格式文件改动的方案有很多，但是被IEFT官方收录的就两种JSON Patch和JSON Merge Patch。

- RFC 6902（JSON Patch）
- RFC 7396（JSON Merge Patch）

JSON Patch是在IETF的RFC 6902中指定的。

## JSON Patch
修改前
```json

{
  "baz": "qux",
  "foo": "bar"
}
```

pacth 
```json
[
  { "op": "replace", "path": "/baz", "value": "boo" },
  { "op": "add", "path": "/hello", "value": ["world"] },
  { "op": "remove", "path": "/foo" }
]
```

修改后
```json
{
  "baz": "boo",
  "hello": ["world"]
}
```

## JSON Merge Patch
Json Merge Patch 描述了如何修改目标JSON 文档的一种格式，即如果被提供的merge patch 中存在目标文档中不存在的成员，则该成员将会新加；
如果被提供的成员与目标成员都存在则做替换修改；在Merge Patch 中的null值意味着删除目标中该对象.

```json
{
     "a": "b",
     "c": {
       "d": "e",
       "f": "g"
     }
  }

```

```http request
# merge patch 用于改变a的值，删除f，
PATCH /target HTTP/1.1
Host: example.org
Content-Type: application/merge-patch+json

   {
     "a":"z",
     "c": {
       "f": null
     }
   }

```

结果
```json
{
    "a": "z",
    "c": {
        "d": "e"
    }
}
```



## golang 库: github.com/evanphx/json-patch
是一个JSON patch库，它提供了对JSON的 Patch 操作和Merge Patch操作

### 功能与特性
操作支持：支持添加（add）、移除（remove）、替换（replace）、移动（move）、复制（copy）和测试（test）操作。

易于集成：可以轻松与现有的 Go 项目集成，特别适合需要频繁修改 JSON 数据的应用场景。

高效：针对大多数常见操作进行了优化，确保在处理大规模 JSON 数据时依然高效。


## JSON patch 在 istio 中的应用
```go
// https://github.com/istio/istio/blob/1ad41e17ee31bbaf7acf131b61da894be8f22303/pkg/kube/inject/inject.go
func IntoObject(injector Injector, sidecarTemplate Templates, valuesConfig ValuesConfig,
	revision string, meshconfig *meshconfig.MeshConfig, in runtime.Object, warningHandler func(string),
) (any, error) {
	out := in.DeepCopyObject()
    // ..
	if injector != nil {
		patchBytes, err = injector.Inject(pod, namespace)
	}
	if err != nil {
		return nil, err
	}
	// TODO(Monkeyanator) istioctl injection still applies just the pod annotation since we don't have
	// the ProxyConfig CRs here.
	if pca, f := metadata.GetAnnotations()[annotation.ProxyConfig.Name]; f {
		var merr error
		meshconfig, merr = mesh.ApplyProxyConfig(pca, meshconfig)
		if merr != nil {
			return nil, merr
		}
	}

	if patchBytes == nil {
		if !injectRequired(IgnoredNamespaces.UnsortedList(), &Config{Policy: InjectionPolicyEnabled}, &pod.Spec, pod.ObjectMeta) {
			warningStr := fmt.Sprintf("===> Skipping injection because %q has sidecar injection disabled\n", fullName)
			if kind != "" {
				warningStr = fmt.Sprintf("===> Skipping injection because %s %q has sidecar injection disabled\n",
					kind, fullName)
			}
			warningHandler(warningStr)
			return out, nil
		}
		params := InjectionParameters{
			pod:        pod,
			deployMeta: deploymentMetadata,
			typeMeta:   typeMeta,
			// Todo replace with some template resolver abstraction
			templates:           sidecarTemplate,
			defaultTemplate:     []string{SidecarTemplateName},
			meshConfig:          meshconfig,
			proxyConfig:         meshconfig.GetDefaultConfig(),
			valuesConfig:        valuesConfig,
			revision:            revision,
			proxyEnvs:           map[string]string{},
			injectedAnnotations: nil,
		}
		patchBytes, err = injectPod(params)
	}
	if err != nil {
		return nil, err
	}
	// 应用 patch 返回patch 后的结果
	patched, err := applyJSONPatchToPod(pod, patchBytes)
	if err != nil {
		return nil, err
	}
	patchedObject, _, err := jsonSerializer.Decode(patched, nil, &corev1.Pod{})
	if err != nil {
		return nil, err
	}
	patchedPod := patchedObject.(*corev1.Pod)
	*metadata = patchedPod.ObjectMeta
	*podSpec = patchedPod.Spec
	return out, nil
}

func applyJSONPatchToPod(input *corev1.Pod, patch []byte) ([]byte, error) {
	objJS, err := runtime.Encode(jsonSerializer, input)
	if err != nil {
		return nil, err
	}

	p, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		return nil, err
	}

	patchedJSON, err := p.Apply(objJS)
	if err != nil {
		return nil, err
	}
	return patchedJSON, nil
}
```


## json merge patch 在 argo rollout 中应用
```go
// https://github.com/argoproj/argo-rollouts/blob/77a139861cef83cfc11aa4a6fd06eacac374afb1/rollout/trafficrouting/istio/istio.go

func updateDestinationRule(ctx context.Context, client dynamic.ResourceInterface, orig []byte, dRule, dRuleNew *DestinationRule) (bool, error) {
	dRuleBytes, err := json.Marshal(dRule)
	if err != nil {
		return false, err
	}
	dRuleNewBytes := destinationRuleReplaceExtraMarshal(dRuleNew)
	log.Debugf("dRuleNewBytes: %s", string(dRuleNewBytes))

	patch, err := jsonpatch.CreateMergePatch(dRuleBytes, dRuleNewBytes)
	if err != nil {
		return false, err
	}
	if string(patch) == "{}" {
		return false, nil
	}
	dRuleNewBytes, err = jsonpatch.MergePatch(orig, patch)
	if err != nil {
		return false, err
	}
	var newDRuleUn unstructured.Unstructured
	err = json.Unmarshal(dRuleNewBytes, &newDRuleUn.Object)
	if err != nil {
		return false, err
	}
	_, err = client.Update(ctx, &newDRuleUn, metav1.UpdateOptions{})
	if err != nil {
		return false, err
	}
	log.Infof("updating destinationrule: %s", string(patch))
	return true, nil
}
```


## 参考
- https://datatracker.ietf.org/doc/html/rfc6902
- https://datatracker.ietf.org/doc/html/rfc7396