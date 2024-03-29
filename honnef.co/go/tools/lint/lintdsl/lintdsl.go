// Package lintdsl provides helpers for implementing static analysis
// checks. Dot-importing this package is encouraged.
package lintdsl

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"honnef.co/go/tools/lint"
	"honnef.co/go/tools/ssa"
)

type packager interface {
	Package() *ssa.Package
}

func CallName(call *ssa.CallCommon) string {
	if call.IsInvoke() {
		return ""
	}
	switch v := call.Value.(type) {
	case *ssa.Function:
		fn, ok := v.Object().(*types.Func)
		if !ok {
			return ""
		}
		return fn.FullName()
	case *ssa.Builtin:
		return v.Name()
	}
	return ""
}

func IsCallTo(call *ssa.CallCommon, name string) bool { return CallName(call) == name }
func IsType(T types.Type, name string) bool           { return types.TypeString(T, nil) == name }

func FilterDebug(instr []ssa.Instruction) []ssa.Instruction {
	var out []ssa.Instruction
	for _, ins := range instr {
		if _, ok := ins.(*ssa.DebugRef); !ok {
			out = append(out, ins)
		}
	}
	return out
}

func IsExample(fn *ssa.Function) bool {
	if !strings.HasPrefix(fn.Name(), "Example") {
		return false
	}
	f := fn.Prog.Fset.File(fn.Pos())
	if f == nil {
		return false
	}
	return strings.HasSuffix(f.Name(), "_test.go")
}

func IsPointerLike(T types.Type) bool {
	switch T := T.Underlying().(type) {
	case *types.Interface, *types.Chan, *types.Map, *types.Pointer:
		return true
	case *types.Basic:
		return T.Kind() == types.UnsafePointer
	}
	return false
}

func IsGenerated(f *ast.File) bool {
	comments := f.Comments
	if len(comments) > 0 {
		comment := comments[0].Text()
		return strings.Contains(comment, "Code generated by") ||
			strings.Contains(comment, "DO NOT EDIT")
	}
	return false
}

func IsIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// isBlank returns whether id is the blank identifier "_".
// If id == nil, the answer is false.
func IsBlank(id ast.Expr) bool {
	ident, _ := id.(*ast.Ident)
	return ident != nil && ident.Name == "_"
}

func IsIntLiteral(expr ast.Expr, literal string) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.INT && lit.Value == literal
}

// Deprecated: use IsIntLiteral instead
func IsZero(expr ast.Expr) bool {
	return IsIntLiteral(expr, "0")
}

func TypeOf(j *lint.Job, expr ast.Expr) types.Type          { return j.Program.Info.TypeOf(expr) }
func IsOfType(j *lint.Job, expr ast.Expr, name string) bool { return IsType(TypeOf(j, expr), name) }

func ObjectOf(j *lint.Job, ident *ast.Ident) types.Object { return j.Program.Info.ObjectOf(ident) }

func IsInTest(j *lint.Job, node lint.Positioner) bool {
	// FIXME(dh): this doesn't work for global variables with
	// initializers
	f := j.Program.SSA.Fset.File(node.Pos())
	return f != nil && strings.HasSuffix(f.Name(), "_test.go")
}

func IsInMain(j *lint.Job, node lint.Positioner) bool {
	if node, ok := node.(packager); ok {
		return node.Package().Pkg.Name() == "main"
	}
	pkg := j.NodePackage(node)
	if pkg == nil {
		return false
	}
	return pkg.Pkg.Name() == "main"
}

func SelectorName(j *lint.Job, expr *ast.SelectorExpr) string {
	sel := j.Program.Info.Selections[expr]
	if sel == nil {
		if x, ok := expr.X.(*ast.Ident); ok {
			pkg, ok := j.Program.Info.ObjectOf(x).(*types.PkgName)
			if !ok {
				// This shouldn't happen
				return fmt.Sprintf("%s.%s", x.Name, expr.Sel.Name)
			}
			return fmt.Sprintf("%s.%s", pkg.Imported().Path(), expr.Sel.Name)
		}
		panic(fmt.Sprintf("unsupported selector: %v", expr))
	}
	return fmt.Sprintf("(%s).%s", sel.Recv(), sel.Obj().Name())
}

func IsNil(j *lint.Job, expr ast.Expr) bool {
	return j.Program.Info.Types[expr].IsNil()
}

func BoolConst(j *lint.Job, expr ast.Expr) bool {
	val := j.Program.Info.ObjectOf(expr.(*ast.Ident)).(*types.Const).Val()
	return constant.BoolVal(val)
}

func IsBoolConst(j *lint.Job, expr ast.Expr) bool {
	// We explicitly don't support typed bools because more often than
	// not, custom bool types are used as binary enums and the
	// explicit comparison is desired.

	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	obj := j.Program.Info.ObjectOf(ident)
	c, ok := obj.(*types.Const)
	if !ok {
		return false
	}
	basic, ok := c.Type().(*types.Basic)
	if !ok {
		return false
	}
	if basic.Kind() != types.UntypedBool && basic.Kind() != types.Bool {
		return false
	}
	return true
}

func ExprToInt(j *lint.Job, expr ast.Expr) (int64, bool) {
	tv := j.Program.Info.Types[expr]
	if tv.Value == nil {
		return 0, false
	}
	if tv.Value.Kind() != constant.Int {
		return 0, false
	}
	return constant.Int64Val(tv.Value)
}

func ExprToString(j *lint.Job, expr ast.Expr) (string, bool) {
	val := j.Program.Info.Types[expr].Value
	if val == nil {
		return "", false
	}
	if val.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(val), true
}

// Dereference returns a pointer's element type; otherwise it returns
// T.
func Dereference(T types.Type) types.Type {
	if p, ok := T.Underlying().(*types.Pointer); ok {
		return p.Elem()
	}
	return T
}

// DereferenceR returns a pointer's element type; otherwise it returns
// T. If the element type is itself a pointer, DereferenceR will be
// applied recursively.
func DereferenceR(T types.Type) types.Type {
	if p, ok := T.Underlying().(*types.Pointer); ok {
		return DereferenceR(p.Elem())
	}
	return T
}

func IsGoVersion(j *lint.Job, minor int) bool {
	return j.Program.GoVersion >= minor
}

func IsCallToAST(j *lint.Job, node ast.Node, name string) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	fn, ok := j.Program.Info.ObjectOf(sel.Sel).(*types.Func)
	return ok && fn.FullName() == name
}

func IsCallToAnyAST(j *lint.Job, node ast.Node, names ...string) bool {
	for _, name := range names {
		if IsCallToAST(j, node, name) {
			return true
		}
	}
	return false
}

func Render(j *lint.Job, x interface{}) string {
	fset := j.Program.SSA.Fset
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

func RenderArgs(j *lint.Job, args []ast.Expr) string {
	var ss []string
	for _, arg := range args {
		ss = append(ss, Render(j, arg))
	}
	return strings.Join(ss, ", ")
}

func Preamble(f *ast.File) string {
	cutoff := f.Package
	if f.Doc != nil {
		cutoff = f.Doc.Pos()
	}
	var out []string
	for _, cmt := range f.Comments {
		if cmt.Pos() >= cutoff {
			break
		}
		out = append(out, cmt.Text())
	}
	return strings.Join(out, "\n")
}

func Inspect(node ast.Node, fn func(node ast.Node) bool) {
	if node == nil {
		return
	}
	ast.Inspect(node, fn)
}
