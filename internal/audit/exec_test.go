package audit

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var (
	shells     = map[string]bool{"sh": true, "bash": true, "zsh": true, "cmd": true, "cmd.exe": true, "powershell": true, "pwsh": true}
	shellFlags = map[string]bool{"-c": true, "/c": true, "/C": true, "-Command": true}
)

// scanExecCalls parses every non-test .go file under roots and reports
// each exec.Command/exec.CommandContext call that runs a shell with a
// command-string flag, or builds an argument by string concatenation or
// fmt.Sprintf.
func scanExecCalls(t *testing.T, roots ...string) []string {
	t.Helper()
	var findings []string
	fset := token.NewFileSet()
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return err
			}
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok || !isExecCommand(call.Fun) {
					return true
				}
				pos := fset.Position(call.Pos())
				var sawShell, sawFlag bool
				for _, arg := range call.Args {
					if s, ok := stringLit(arg); ok {
						sawShell = sawShell || shells[s]
						sawFlag = sawFlag || shellFlags[s]
					}
					if isInterpolated(arg) {
						findings = append(findings, pos.String()+": argument built by string interpolation")
					}
				}
				if sawShell && sawFlag {
					findings = append(findings, pos.String()+": runs a shell with a command-string flag")
				}
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return findings
}

func isExecCommand(fun ast.Expr) bool {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "exec" && (sel.Sel.Name == "Command" || sel.Sel.Name == "CommandContext")
}

func stringLit(e ast.Expr) (string, bool) {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(lit.Value)
	return s, err == nil
}

func isInterpolated(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.BinaryExpr:
		return v.Op == token.ADD
	case *ast.CallExpr:
		if sel, ok := v.Fun.(*ast.SelectorExpr); ok {
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "fmt" && strings.HasPrefix(sel.Sel.Name, "Sprint") {
				return true
			}
		}
	}
	return false
}

func TestNoShellOrInterpolatedExecCalls(t *testing.T) {
	if findings := scanExecCalls(t, filepath.Join("..", "..", "internal"), filepath.Join("..", "..", "cmd")); len(findings) > 0 {
		t.Errorf("unsafe exec calls:\n%s", strings.Join(findings, "\n"))
	}
}

// The scanner must actually catch what it claims to; otherwise the test
// above passes vacuously.
func TestExecScannerCatchesUnsafeCalls(t *testing.T) {
	dir := t.TempDir()
	bad := `package bad

import (
	"fmt"
	"os/exec"
)

func a(x string) { exec.Command("sh", "-c", "echo hi") }
func b(x string) { exec.Command("docker", "rm "+x) }
func c(x string) { exec.Command("docker", fmt.Sprintf("rm %s", x)) }
func d(x string) { exec.Command("cmd", "/c", "start", x) }
`
	if err := os.WriteFile(filepath.Join(dir, "bad.go"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := scanExecCalls(t, dir); len(got) != 4 {
		t.Errorf("scanner found %d issues in the fixture, want 4:\n%s", len(got), strings.Join(got, "\n"))
	}
}
