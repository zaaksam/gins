package command

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
	"github.com/zaaksam/gins/extend/cmdutil"
	"github.com/zaaksam/gins/extend/fileutil"
)

// genInstance 生成
var genInstance gen

type gen struct {
	*cli.Command

	isInit bool
	wd     string
}

type genModel struct {
	Name   string
	Fields []*genField
}

type genField struct {
	Name            string
	Type            string
	RawType         string
	XormName        string
	JSONName        string
	IsJSONTagString bool
}

// Gen 生成代码
func Gen() *cli.Command {
	genInstance.onceInit()

	return genInstance.Command
}

func (g *gen) onceInit() {
	if g.isInit {
		return
	}
	g.isInit = true

	g.wd, _ = os.Getwd()

	g.Command = &cli.Command{
		Name:  "gen",
		Usage: "[默认可省略] 生成 orm model 代码，示例：xmodel gen demo.go",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "r",
				Usage: "生成 orm.Register 函数调用",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "d",
				Usage: "调试代码生成",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			args := ctx.Args()
			if args.Len() == 0 {
				err = errors.New("model 文件名必须提供")
				return
			}
			modelFile := args.Get(0)

			if !strings.HasSuffix(modelFile, ".go") {
				err = fmt.Errorf("model 文件 %s 不是 .go 文件", modelFile)
				return
			}

			modelFilePath := filepath.Join(g.wd, modelFile)

			// 解析文件源代码
			pkgName, mds, err := g.parseFile(ctx, modelFilePath)
			if err != nil {
				return
			}

			// 生成 _orm.go 文件路径信息
			ormModelFile := strings.Replace(modelFile, ".go", "_orm.go", 1)
			ormModelFilePath := filepath.Join(g.wd, ormModelFile)

			// 从模板生成源代码
			err = g.parseTemplate(ctx, ormModelFilePath, pkgName, mds)
			if err != nil {
				return
			}

			// go fmt 格式化代码
			_, err = cmdutil.Exec("go", []string{"fmt", ormModelFile})
			return
		},
	}
}

func (g *gen) parseFile(ctx *cli.Context, modelFilePath string) (pkgName string, codeModels []*genModel, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, modelFilePath, nil, parser.AllErrors)
	if err != nil {
		return
	}

	if ctx.Bool("d") {
		// 查看 ast 结构
		ast.Print(fset, f)

		err = errors.New("因调试代码中止")
		return
	}

	pkgName = f.Name.Name

	codeModels = make([]*genModel, 0, 10)
	var codeModel *genModel

	// 解析规则需要参照 ast.Print 结果去处理
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// 是否 type 定义
		if genDecl.Tok.String() != "type" {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 是否 struct 结构
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 解析 struct 为 genModel 结构
			codeModel, err = g.parseASTStructType(typeSpec.Name.Name, structType)
			if err != nil {
				return
			}

			codeModels = append(codeModels, codeModel)
		}
	}

	if len(codeModels) == 0 {
		err = errors.New(modelFilePath + " 没有解析到 orm.Model 结构体")
	}
	return
}

func (g *gen) parseASTStructType(name string, structType *ast.StructType) (codeModel *genModel, err error) {
	l := len(structType.Fields.List)
	codeFields := make([]*genField, 0, l)

	var (
		codeField *genField
		ok        bool
	)

	for i := 0; i < l; i++ {
		// 解析 ast.Field 为 genField 结构
		codeField, ok = g.parseASTField(structType.Fields.List[i])
		if !ok {
			continue
		}

		// 第一行需要为 orm.Model
		if i == 0 && codeField.Type != "orm.Model" {
			err = errors.New(name + " 结构顶部没有嵌套 orm.Model")
			return
		}

		// 字段类型不为 *orm.Field
		if !strings.HasPrefix(codeField.Type, "*orm.Field") {
			continue
		}

		codeFields = append(codeFields, codeField)
	}

	if len(codeFields) == 0 {
		err = errors.New(name + " 结构没有发现任意 *orm.Field 字段")
		return

	}

	codeModel = &genModel{
		Name:   name,
		Fields: codeFields,
	}
	return
}

func (*gen) parseASTField(field *ast.Field) (codeField *genField, ok bool) {
	var (
		selectorExpr *ast.SelectorExpr // 嵌套类型
		starExpr     *ast.StarExpr     // 指针类型
		indexExpr    *ast.IndexExpr    // 自定义类型
		ident        *ast.Ident        // 常规类型
	)

	codeField = &genField{}

	// 除了嵌套类型，其他都有 Names
	if len(field.Names) == 1 {
		codeField.Name = field.Names[0].Name
	}

	if field.Tag != nil {
		tag := field.Tag.Value

		// 获取 json tag 信息
		tags := strings.SplitN(tag, `json:"`, 2)
		if len(tags) == 2 {
			tags = strings.SplitN(tags[1], `"`, 2)
			if len(tags) == 2 {
				jsonTags := strings.Split(tags[0], ",")

				for i, l := 0, len(jsonTags); i < l; i++ {
					if i == 0 {
						codeField.JSONName = jsonTags[i]
						continue
					}

					if jsonTags[i] == "string" {
						codeField.IsJSONTagString = true
					}
				}
			}
		}

		// 获取 xorm tag 信息
		tags = strings.SplitN(tag, `xorm:"`, 2)
		if len(tags) == 2 {
			if strings.Contains(tags[1], "'") {
				tags = strings.Split(tags[1], "'")
				if len(tags) == 3 {
					codeField.XormName = tags[1]
				}
			} else {
				tags = strings.SplitN(tags[1], `"`, 2)
				if len(tags) == 2 {
					codeField.XormName = tags[0]
				}
			}
		}
	}

	// 嵌套类型检查 orm.Model
	selectorExpr, ok = field.Type.(*ast.SelectorExpr)
	if ok {
		ident, ok = selectorExpr.X.(*ast.Ident)
		if !ok {
			return
		}

		codeField.RawType = ident.Name + "." + selectorExpr.Sel.Name
		codeField.Type = codeField.RawType
		return
	}

	// 指针类型检查 *orm.Field[uint8]
	starExpr, ok = field.Type.(*ast.StarExpr)
	if ok {
		indexExpr, ok = starExpr.X.(*ast.IndexExpr)
		if !ok {
			return
		}
		selectorExpr, ok = indexExpr.X.(*ast.SelectorExpr)
		if !ok {
			return
		}
		ident, ok = selectorExpr.X.(*ast.Ident)
		if !ok {
			return
		}

		codeField.RawType = indexExpr.Index.(*ast.Ident).Name
		codeField.Type = "*" + ident.Name + "." + selectorExpr.Sel.Name + "[" + codeField.RawType + "]"
		return
	}

	// 自定义类型 orm.Field[uint8]
	// indexExpr, ok = field.Type.(*ast.IndexExpr)
	// if ok {
	// 	selectorExpr, ok = indexExpr.X.(*ast.SelectorExpr)
	// 	if !ok {
	// 		return
	// 	}
	// 	ident, ok = selectorExpr.X.(*ast.Ident)
	// 	if !ok {
	// 		return
	// 	}

	// 	codeField.Type = ident.Name + "." + selectorExpr.Sel.Name + "[" + indexExpr.Index.(*ast.Ident).Name + "]"
	// 	return
	// }

	// 常规类型
	// ident, ok = field.Type.(*ast.Ident)
	// if ok {
	// 	codeField.Type = ident.Name
	// }
	return
}

func (*gen) parseTemplate(ctx *cli.Context, ormModelFilePath, pkgName string, mds []*genModel) (err error) {
	var tmpl *template.Template
	tmpl, err = template.New("orm").Parse(genTemplate)
	if err != nil {
		return
	}

	file, err := fileutil.NewFile(ormModelFilePath)
	if err != nil {
		return
	}

	data := make(map[string]interface{})
	data["pkg"] = pkgName
	data["isReg"] = ctx.Bool("r")
	data["models"] = mds

	err = tmpl.Execute(file, data)
	return
}
