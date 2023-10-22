<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [cobra在k8s中的应用](#cobra%E5%9C%A8k8s%E4%B8%AD%E7%9A%84%E5%BA%94%E7%94%A8)
  - [kubectl 的启动工程](#kubectl-%E7%9A%84%E5%90%AF%E5%8A%A8%E5%B7%A5%E7%A8%8B)
    - [1. 主函数路径： cmd/kubectl/kubectl.go](#1-%E4%B8%BB%E5%87%BD%E6%95%B0%E8%B7%AF%E5%BE%84-cmdkubectlkubectlgo)
    - [2. 包具体实现： pkg/kubectl/cmd/cmd.go](#2-%E5%8C%85%E5%85%B7%E4%BD%93%E5%AE%9E%E7%8E%B0-pkgkubectlcmdcmdgo)
    - [3. cobra具体定义](#3-cobra%E5%85%B7%E4%BD%93%E5%AE%9A%E4%B9%89)
    - [4.开始添加子命令](#4%E5%BC%80%E5%A7%8B%E6%B7%BB%E5%8A%A0%E5%AD%90%E5%91%BD%E4%BB%A4)
      - [5.拿kubectl apply命令讲解](#5%E6%8B%BFkubectl-apply%E5%91%BD%E4%BB%A4%E8%AE%B2%E8%A7%A3)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# cobra在k8s中的应用

## kubectl 的启动工程

### 1. 主函数路径： cmd/kubectl/kubectl.go
```go

func main() {
	command := cmd.NewDefaultKubectlCommand()
	    
	// ...
	
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
```
### 2. 包具体实现： pkg/kubectl/cmd/cmd.go
```go
// NewDefaultKubectlCommand creates the `kubectl` command with default arguments
func NewDefaultKubectlCommand() *cobra.Command {
	return NewDefaultKubectlCommandWithArgs(NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes), os.Args, os.Stdin, os.Stdout, os.Stderr)
}

// NewDefaultKubectlCommandWithArgs creates the `kubectl` command with arguments
func NewDefaultKubectlCommandWithArgs(pluginHandler PluginHandler, args []string, in io.Reader, out, errout io.Writer) *cobra.Command {
	cmd := NewKubectlCommand(in, out, errout)
    // ...
	return cmd
}
```

### 3. cobra具体定义
```go
func NewKubectlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	warningHandler := rest.NewWarningWriter(err, rest.WarningWriterOptions{Deduplicate: true, Color: term.AllowsColorOutput(err)})
	warningsAsErrors := false

	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "kubectl",
		Short: i18n.T("kubectl controls the Kubernetes cluster manager"),
		Long: templates.LongDesc(`
      kubectl controls the Kubernetes cluster manager.

      Find more information at:
            https://kubernetes.io/docs/reference/kubectl/overview/`),
		Run: runHelp,
		// Hook before and after Run initialize and write profiles to disk,
		// respectively.
		PersistentPreRunE: func(*cobra.Command, []string) error {
			rest.SetDefaultWarningHandler(warningHandler)
			// 这里是做pprof性能分析，跳转到对应代码可以看到，我们可以用参数 --profile xxx 来采集性能指标，默认保存在当前目录下的profile.pprof中
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			//  保存pprof性能分析指标
			if err := flushProfiling(); err != nil {
				return err
			}
			// 打印warning条数
			if warningsAsErrors {
				count := warningHandler.WarningCount()
				switch count {
				case 0:
					// no warnings
				case 1:
					return fmt.Errorf("%d warning received", count)
				default:
					return fmt.Errorf("%d warnings received", count)
				}
			}
			return nil
		},
		//bash自动补齐功能
		BashCompletionFunction: bashCompletionFunc,
	}

	flags := cmds.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	addProfilingFlags(flags)

	flags.BoolVar(&warningsAsErrors, "warnings-as-errors", warningsAsErrors, "Treat warnings received from the server as errors and exit with a non-zero exit code")

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	matchVersionKubeConfigFlags.AddFlags(cmds.PersistentFlags())

	cmds.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// 实例化Factory接口，工厂模式
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	// Sending in 'nil' for the getLanguageFn() results in using
	// the LANG environment variable.
	//
	// TODO: Consider adding a flag or file preference for setting
	// the language, instead of just loading from the LANG env. variable.
	i18n.LoadTranslations("kubectl", nil)

	// From this point and forward we get warnings on flags that contain "_" separators
	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}

	// kubectl定义了7类命令，结合Message和各个子命令的package名来看
	groups := templates.CommandGroups{
		{
			// 1. 初级命令，包括 create/expose/run/set
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				create.NewCmdCreate(f, ioStreams),
				expose.NewCmdExposeService(f, ioStreams),
				run.NewCmdRun(f, ioStreams),
				set.NewCmdSet(f, ioStreams),
			},
		},
		{
			// 2. 中级命令，包括explain/get/edit/delete
			Message: "Basic Commands (Intermediate):",
			Commands: []*cobra.Command{
				explain.NewCmdExplain("kubectl", f, ioStreams),
				get.NewCmdGet("kubectl", f, ioStreams),
				edit.NewCmdEdit(f, ioStreams),
				delete.NewCmdDelete(f, ioStreams),
			},
		},
		{
			// 3. 部署命令，包括 rollout/scale/autoscale
			Message: "Deploy Commands:",
			Commands: []*cobra.Command{
				rollout.NewCmdRollout(f, ioStreams),
				scale.NewCmdScale(f, ioStreams),
				autoscale.NewCmdAutoscale(f, ioStreams),
			},
		},
		{
			// 4. 集群管理命令，包括 cerfificate/cluster-info/top/cordon/drain/taint
			Message: "Cluster Management Commands:",
			Commands: []*cobra.Command{
				certificates.NewCmdCertificate(f, ioStreams),
				clusterinfo.NewCmdClusterInfo(f, ioStreams),
				top.NewCmdTop(f, ioStreams),
				drain.NewCmdCordon(f, ioStreams),
				drain.NewCmdUncordon(f, ioStreams),
				drain.NewCmdDrain(f, ioStreams),
				taint.NewCmdTaint(f, ioStreams),
			},
		},
		{
			// 5. 故障排查和调试，包括 describe/logs/attach/exec/port-forward/proxy/cp/auth
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				describe.NewCmdDescribe("kubectl", f, ioStreams),
				logs.NewCmdLogs(f, ioStreams),
				attach.NewCmdAttach(f, ioStreams),
				cmdexec.NewCmdExec(f, ioStreams),
				portforward.NewCmdPortForward(f, ioStreams),
				proxy.NewCmdProxy(f, ioStreams),
				cp.NewCmdCp(f, ioStreams),
				auth.NewCmdAuth(f, ioStreams),
			},
		},
		{
			// 6. 高级命令，包括diff/apply/patch/replace/wait/convert/kustomize
			Message: "Advanced Commands:",
			Commands: []*cobra.Command{
				diff.NewCmdDiff(f, ioStreams),
				apply.NewCmdApply("kubectl", f, ioStreams),
				patch.NewCmdPatch(f, ioStreams),
				replace.NewCmdReplace(f, ioStreams),
				wait.NewCmdWait(f, ioStreams),
				convert.NewCmdConvert(f, ioStreams),
				kustomize.NewCmdKustomize(ioStreams),
			},
		},
		{
			// 7. 设置命令，包括label，annotate，completion
			Message: "Settings Commands:",
			Commands: []*cobra.Command{
				label.NewCmdLabel(f, ioStreams),
				annotate.NewCmdAnnotate("kubectl", f, ioStreams),
				completion.NewCmdCompletion(ioStreams.Out, ""),
			},
		},
	}
	// 开始添加子命令
	groups.Add(cmds)

	filters := []string{"options"}

	// Hide the "alpha" subcommand if there are no alpha commands in this build.
	alpha := cmdpkg.NewCmdAlpha(f, ioStreams)
	if !alpha.HasSubCommands() {
		filters = append(filters, alpha.Name())
	}

	templates.ActsAsRootCommand(cmds, filters, groups...)

	// 代码补全相关
	for name, completion := range bashCompletionFlags {
		if cmds.Flag(name) != nil {
			if cmds.Flag(name).Annotations == nil {
				cmds.Flag(name).Annotations = map[string][]string{}
			}
			cmds.Flag(name).Annotations[cobra.BashCompCustom] = append(
				cmds.Flag(name).Annotations[cobra.BashCompCustom],
				completion,
			)
		}
	}

	// 添加其余子命令，包括 alpha/config/plugin/version/api-versions/api-resources/options
	cmds.AddCommand(alpha)
	cmds.AddCommand(cmdconfig.NewCmdConfig(f, clientcmd.NewDefaultPathOptions(), ioStreams))
	cmds.AddCommand(plugin.NewCmdPlugin(f, ioStreams))
	cmds.AddCommand(version.NewCmdVersion(f, ioStreams))
	cmds.AddCommand(apiresources.NewCmdAPIVersions(f, ioStreams))
	cmds.AddCommand(apiresources.NewCmdAPIResources(f, ioStreams))
	cmds.AddCommand(options.NewCmdOptions(ioStreams.Out))

	return cmds
}
```

### 4.开始添加子命令

```go
func (g CommandGroups) Add(c *cobra.Command) {
	for _, group := range g {
		c.AddCommand(group.Commands...)
	}
}

```
#### 5.拿kubectl apply命令讲解
```go
//  apply.NewCmdApply("kubectl", f, ioStreams) 

type Factory interface {
    genericclioptions.RESTClientGetter
    
    // DynamicClient returns a dynamic client ready for use
    DynamicClient() (dynamic.Interface, error)
    
    // KubernetesClientSet gives you back an external clientset
    KubernetesClientSet() (*kubernetes.Clientset, error)
    
    // Returns a RESTClient for accessing Kubernetes resources or an error.
    RESTClient() (*restclient.RESTClient, error)
    
    // NewBuilder returns an object that assists in loading objects from both disk and the server
    // and which implements the common patterns for CLI interactions with generic resources.
    NewBuilder() *resource.Builder
    
    // Returns a RESTClient for working with the specified RESTMapping or an error. This is intended
    // for working with arbitrary resources and is not guaranteed to point to a Kubernetes APIServer.
    ClientForMapping(mapping *meta.RESTMapping) (resource.RESTClient, error)
    // Returns a RESTClient for working with Unstructured objects.
    UnstructuredClientForMapping(mapping *meta.RESTMapping) (resource.RESTClient, error)
    
    // Returns a schema that can validate objects stored on disk.
    Validator(validate bool) (validation.Schema, error)
    // OpenAPISchema returns the schema openapi schema definition
    OpenAPISchema() (openapi.Resources, error)
}
```
解释f是个工厂
```go
// NewCmdApply creates the `apply` command
func NewCmdApply(baseName string, f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewApplyOptions(ioStreams)

	// Store baseName for use in printing warnings / messages involving the base command name.
	// This is useful for downstream command that wrap this one.
	o.cmdBaseName = baseName

	cmd := &cobra.Command{
		Use:                   "apply (-f FILENAME | -k DIRECTORY)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Apply a configuration to a resource by filename or stdin"),
		Long:                  applyLong,
		Example:               applyExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(validateArgs(cmd, args))
			cmdutil.CheckErr(validatePruneAll(o.Prune, o.All, o.Selector))
			cmdutil.CheckErr(o.Run())
		},
	}

	// bind flag structs
	o.DeleteFlags.AddFlags(cmd)
	o.RecordFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.Overwrite, "overwrite", o.Overwrite, "Automatically resolve conflicts between the modified and live configuration by using values from the modified configuration")
	cmd.Flags().BoolVar(&o.Prune, "prune", o.Prune, "Automatically delete resource objects, including the uninitialized ones, that do not appear in the configs and are created by either apply or create --save-config. Should be used with either -l or --all.")
	cmdutil.AddValidateFlags(cmd)
	cmd.Flags().StringVarP(&o.Selector, "selector", "l", o.Selector, "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Select all resources in the namespace of the specified resource types.")
	cmd.Flags().StringArrayVar(&o.PruneWhitelist, "prune-whitelist", o.PruneWhitelist, "Overwrite the default whitelist with <group/version/kind> for --prune")
	cmd.Flags().BoolVar(&o.OpenAPIPatch, "openapi-patch", o.OpenAPIPatch, "If true, use openapi to calculate diff when the openapi presents and the resource can be found in the openapi spec. Otherwise, fall back to use baked-in types.")
	cmdutil.AddDryRunFlag(cmd)
	cmdutil.AddServerSideApplyFlags(cmd)
	cmdutil.AddFieldManagerFlagVar(cmd, &o.FieldManager, FieldManagerClientSideApply)

	// apply subcommands
	cmd.AddCommand(NewCmdApplyViewLastApplied(f, ioStreams))
	cmd.AddCommand(NewCmdApplySetLastApplied(f, ioStreams))
	cmd.AddCommand(NewCmdApplyEditLastApplied(f, ioStreams))

	return cmd
}
```
主要看run函数
```go
func (o *ApplyOptions) Run() error {
	// ..
	// 获取资源对象
	infos, err := o.GetObjects()
    // ...
	// Iterate through all objects, applying each one.
	for _, info := range infos {
		if err := o.applyOneObject(info); err != nil {
			errs = append(errs, err)
		}
	}
    // ...

	return nil
}

func (o *ApplyOptions) GetObjects() ([]*resource.Info, error) {
	var err error = nil
	if !o.objectsCached {
		// include the uninitialized objects by default if --prune is true
		// unless explicitly set --include-uninitialized=false
		r := o.Builder.
			Unstructured().
			Schema(o.Validator).
			ContinueOnError().
			NamespaceParam(o.Namespace).DefaultNamespace().
			FilenameParam(o.EnforceNamespace, &o.DeleteOptions.FilenameOptions).
			LabelSelectorParam(o.Selector).
			Flatten().
			Do()
		o.objects, err = r.Infos()
		o.objectsCached = true
	}
	return o.objects, err
}
```

先看最后的Do(),主要返回了一系列的visitors
```go
func (b *Builder) Do() *Result {
	r := b.visitorResult()
    // ...
	return r
}


// 根据一系列条件去生成不同的visitor，如path,selectormnems
func (b *Builder) visitorResult() *Result {
	if len(b.errs) > 0 {
		return &Result{err: utilerrors.NewAggregate(b.errs)}
	}

	if b.selectAll {
		selector := labels.Everything().String()
		b.labelSelector = &selector
	}

	// visit items specified by paths
	if len(b.paths) != 0 {
		return b.visitByPaths()
	}

	// visit selectors
	if b.labelSelector != nil || b.fieldSelector != nil {
		return b.visitBySelector()
	}

	// visit items specified by resource and name
	if len(b.resourceTuples) != 0 {
		return b.visitByResource()
	}

	// visit items specified by name
	if len(b.names) != 0 {
		return b.visitByName()
	}

	if len(b.resources) != 0 {
		for _, r := range b.resources {
			_, err := b.mappingFor(r)
			if err != nil {
				return &Result{err: err}
			}
		}
		return &Result{err: fmt.Errorf("resource(s) were provided, but no name, label selector, or --all flag specified")}
	}
	return &Result{err: missingResourceError}
}
```
Visitor接口
```go
// Visitor lets clients walk a list of resources.
type Visitor interface {
	Visit(VisitorFunc) error
}

type VisitorFunc func(*Info, error) error

// 最终解析处理的资源信息
type Info struct {
	// Client will only be present if this builder was not local
	Client RESTClient
	// Mapping will only be present if this builder was not local
	Mapping *meta.RESTMapping

	// Namespace will be set if the object is namespaced and has a specified value.
	Namespace string
	Name      string

	// Optional, Source is the filename or URL to template file (.json or .yaml),
	// or stdin to use to handle the resource
	Source string
	// Optional, this is the most recent value returned by the server if available. It will
	// typically be in unstructured or internal forms, depending on how the Builder was
	// defined. If retrieved from the server, the Builder expects the mapping client to
	// decide the final form. Use the AsVersioned, AsUnstructured, and AsInternal helpers
	// to alter the object versions.
	Object runtime.Object
	// Optional, this is the most recent resource version the server knows about for
	// this type of resource. It may not match the resource version of the object,
	// but if set it should be equal to or newer than the resource version of the
	// object (however the server defines resource version).
	ResourceVersion string
}
```

具体介绍fileVisitor
```go
// FileVisitor is wrapping around a StreamVisitor, to handle open/close files
type FileVisitor struct {
	Path string
	*StreamVisitor
}

// Visit in a FileVisitor is just taking care of opening/closing files
func (v *FileVisitor) Visit(fn VisitorFunc) error {
	var f *os.File
	if v.Path == constSTDINstr {
		f = os.Stdin
	} else {
		// 打开文件
		var err error
		f, err = os.Open(v.Path)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	// TODO: Consider adding a flag to force to UTF16, apparently some
	// Windows tools don't write the BOM
	utf16bom := unicode.BOMOverride(unicode.UTF8.NewDecoder())
	v.StreamVisitor.Reader = transform.NewReader(f, utf16bom)

	return v.StreamVisitor.Visit(fn)
}
```
真正调用的是StreamVisitor的Visit方法
```go
// NewStreamVisitor is a helper function that is useful when we want to change the fields of the struct but keep calls the same.
func NewStreamVisitor(r io.Reader, mapper *mapper, source string, schema ContentValidator) *StreamVisitor {
	return &StreamVisitor{
		Reader: r,
		mapper: mapper,   // 映射获取相关的资源
		Source: source,
		Schema: schema,
	}
}

// Visit implements Visitor over a stream. StreamVisitor is able to distinct multiple resources in one stream.
func (v *StreamVisitor) Visit(fn VisitorFunc) error {
	// 使用yaml解析
	d := yaml.NewYAMLOrJSONDecoder(v.Reader, 4096)
	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error parsing %s: %v", v.Source, err)
		}
		// TODO: This needs to be able to handle object in other encodings and schemas.
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		if err := ValidateSchema(ext.Raw, v.Schema); err != nil {
			return fmt.Errorf("error validating %q: %v", v.Source, err)
		}
		// 解析出来的资源信息
		info, err := v.infoForData(ext.Raw, v.Source)
		if err != nil {
			if fnErr := fn(info, err); fnErr != nil {
				return fnErr
			}
			continue
		}
		if err := fn(info, nil); err != nil {
			return err
		}
	}
}
```

mapper的初始化
```go
// Unstructured updates the builder so that it will request and send unstructured
// objects. Unstructured objects preserve all fields sent by the server in a map format
// based on the object's JSON structure which means no data is lost when the client
// reads and then writes an object. Use this mode in preference to Internal unless you
// are working with Go types directly.
func (b *Builder) Unstructured() *Builder {
	if b.mapper != nil {
		b.errs = append(b.errs, fmt.Errorf("another mapper was already selected, cannot use unstructured types"))
		return b
	}
	b.objectTyper = unstructuredscheme.NewUnstructuredObjectTyper()
	b.mapper = &mapper{
		localFn:      b.isLocal,
		restMapperFn: b.restMapperFn,
		clientFn:     b.getClient,
		decoder:      &metadataValidatingDecoder{unstructured.UnstructuredJSONScheme},
	}

	return b
}

```



看完Do，再看visitor对应入口的builder

1、path

       func (b *Builder) Path(paths ...string) *Builder {
2、selector

       func (b *Builder) SelectorParam(s string) *Builder {
       func (b *Builder) Selector(selector labels.Selector) *Builder {
       func (b *Builder) SelectAllParam(selectAll bool) *Builder

3、namespace

       func (b *Builder) NamespaceParam(namespace string) *Builder {
       func (b *Builder) DefaultNamespace() *Builder {
       func (b *Builder) AllNamespaces(allNamespace bool) *Builder {
4、resource

       func (b *Builder) ResourceNames(resource string, names ...string) *Builder {
       func (b *Builder) ResourceTypes(types ...string) *Builder {
       func (b *Builder) ResourceTypeOrNameArgs(allowEmptySelector bool, args ...string)        *Builder

5、url

       func (b *Builder) URL(urls ...*url.URL) *Builder {
6、stream

       func (b *Builder) Stream(r io.Reader, name string) *Builder {
7、stdin

       func (b *Builder) Stdin() *Builder {}'


拿FilenameParam(o.EnforceNamespace, &o.DeleteOptions.FilenameOptions)做案例
```go
r := o.Builder.
        Unstructured().
        Schema(o.Validator).
        ContinueOnError().
        NamespaceParam(o.Namespace).DefaultNamespace().
        FilenameParam(o.EnforceNamespace, &o.DeleteOptions.FilenameOptions).
        LabelSelectorParam(o.Selector).
        Flatten().
        Do()
```
```go
func (b *Builder) FilenameParam(enforceNamespace bool, filenameOptions *FilenameOptions) *Builder {
    // ...
	paths := filenameOptions.Filenames
	for _, s := range paths {
		switch {
		// 控制流终端
		case s == "-":
			b.Stdin()
        //url 中获取
		case strings.Index(s, "http://") == 0 || strings.Index(s, "https://") == 0:
			url, err := url.Parse(s)
			if err != nil {
				b.errs = append(b.errs, fmt.Errorf("the URL passed to filename %q is not valid: %v", s, err))
				continue
			}
			b.URL(defaultHttpGetAttempts, url)
        // 默认文件当中
		default:
			if !recursive {
				b.singleItemImplied = true
			}
			b.Path(recursive, s)
		}
	}
    // ...

	return b
}
```

开始创建fileVisitor
```go
// Path accepts a set of paths that may be files, directories (all can containing
// one or more resources). Creates a FileVisitor for each file and then each
// FileVisitor is streaming the content to a StreamVisitor. If ContinueOnError() is set
// prior to this method being called, objects on the path that are unrecognized will be
// ignored (but logged at V(2)).
func (b *Builder) Path(recursive bool, paths ...string) *Builder {
	for _, p := range paths {
        // ...
		visitors, err := ExpandPathsToFileVisitors(b.mapper, p, recursive, FileExtensions, b.schema)
		if err != nil {
			b.errs = append(b.errs, fmt.Errorf("error reading %q: %v", p, err))
		}
		if len(visitors) > 1 {
			b.dir = true
		}

		b.paths = append(b.paths, visitors...)
	}
    // ... 
	return b
}


func ExpandPathsToFileVisitors(mapper *mapper, paths string, recursive bool, extensions []string, schema ContentValidator) ([]Visitor, error) {
	var visitors []Visitor
	err := filepath.Walk(paths, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if path != paths && !recursive {
				return filepath.SkipDir
			}
			return nil
		}
		// Don't check extension if the filepath was passed explicitly
		if path != paths && ignoreFile(path, extensions) {
			return nil
		}
		// 开始创建FileVisitor
		visitor := &FileVisitor{
			Path:          path,
			StreamVisitor: NewStreamVisitor(nil, mapper, path, schema),
		}

		visitors = append(visitors, visitor)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return visitors, nil
}
```


