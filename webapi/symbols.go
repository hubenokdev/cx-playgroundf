package webapi

import (
	"unicode"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
)

// ProgramMetaResp is a program meta data response.
type ProgramMetaResp struct {
	UsedHeapMemory types.Pointer `json:"used_heap_memory"`
	FreeHeapMemory types.Pointer `json:"free_heap_memory"`
	StackSize      types.Pointer `json:"stack_size"`
	CallStackSize  int           `json:"call_stack_size"`
}

func extractProgramMeta(pg *ast.CXProgram) ProgramMetaResp {
	return ProgramMetaResp{
		UsedHeapMemory: pg.Heap.Pointer - pg.Heap.StartsAt,
		FreeHeapMemory: pg.Heap.Size - pg.Heap.StartsAt,
		StackSize:      pg.Stack.Size,
		CallStackSize:  len(pg.CallStack),
	}
}

// ExportedSymbol represents a single exported cx symbol.
type ExportedSymbol struct {
	Name      string      `json:"name"`
	Signature interface{} `json:"signature,omitempty"`
	Type      types.Code  `json:"type"`
	TypeName  string      `json:"type_name"`
}

// ExportedSymbolsResp is the exported symbols response.
type ExportedSymbolsResp struct {
	Functions []ExportedSymbol `json:"functions"`
	Structs   []ExportedSymbol `json:"structs"`
	Globals   []ExportedSymbol `json:"globals"`
}

func extractExportedSymbols(prgrm *ast.CXProgram, pkg *ast.CXPackage) ExportedSymbolsResp {
	resp := ExportedSymbolsResp{
		Functions: make([]ExportedSymbol, 0, len(pkg.Functions)),
		Structs:   make([]ExportedSymbol, 0, len(pkg.Structs)),
		Globals:   make([]ExportedSymbol, 0, len(pkg.Globals.Fields)),
	}

	for _, fIdx := range pkg.Functions {
		f := prgrm.GetFunctionFromArray(fIdx)

		if isExported(f.Name) {
			resp.Functions = append(resp.Functions, displayCXFunction(prgrm, pkg, f))
		}
	}

	for _, sIdx := range pkg.Structs {
		strct := &prgrm.CXStructs[sIdx]
		if isExported(strct.Name) {
			resp.Structs = append(resp.Structs, displayCXStruct(prgrm, strct))
		}
	}

	for _, gIdx := range pkg.Globals.Fields {
		typeSigGlobal := prgrm.GetCXTypeSignatureFromArray(gIdx)

		if isExported(typeSigGlobal.Name) {
			resp.Globals = append(resp.Globals, displayCXGlobal(prgrm, typeSigGlobal))
		}
	}

	return resp
}

func displayCXFunction(prgrm *ast.CXProgram, pkg *ast.CXPackage, f *ast.CXFunction) ExportedSymbol {
	return ExportedSymbol{
		Name:      f.Name,
		Signature: ast.SignatureStringOfFunction(prgrm, pkg, f),
		Type:      types.FUNC,
		TypeName:  types.FUNC.Name(),
	}
}

func displayCXStruct(prgrm *ast.CXProgram, s *ast.CXStruct) ExportedSymbol {
	return ExportedSymbol{
		Name:      s.Name,
		Signature: ast.SignatureStringOfStruct(prgrm, s),
		Type:      types.STRUCT,
		TypeName:  "struct",
	}
}

func displayCXGlobal(prgrm *ast.CXProgram, glbl *ast.CXTypeSignature) ExportedSymbol {
	glblType := types.Code(1)

	switch glbl.Type {
	case ast.TYPE_ATOMIC:
		glblType = types.Code(glbl.Meta)
	case ast.TYPE_POINTER_ATOMIC:
		glblType = types.Code(glbl.Meta)
	case ast.TYPE_ARRAY_ATOMIC:
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(glbl.Meta)
		glblType = types.Code(arrDetails.Type)
	case ast.TYPE_POINTER_ARRAY_ATOMIC:
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(glbl.Meta)
		glblType = types.Code(arrDetails.Type)
	case ast.TYPE_SLICE_ATOMIC:
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(glbl.Meta)
		glblType = types.Code(arrDetails.Type)
	case ast.TYPE_CXARGUMENT_DEPRECATE:
		arg := prgrm.GetCXArgFromArray(ast.CXArgumentIndex(glbl.Meta))
		glblType = arg.Type
	}
	return ExportedSymbol{
		Name:      glbl.Name,
		Signature: nil,
		Type:      glblType,
		TypeName:  glblType.Name(),
	}
}

func isExported(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper(rune(s[0]))
}
