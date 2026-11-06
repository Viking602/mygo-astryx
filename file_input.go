package astryx

import (
	"fmt"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type UploadFile struct {
	Path, Name, MIME string
	Size             int64
}
type FileInputState struct {
	Files    []UploadFile
	Error    string
	Busy     bool
	Progress float64
	pending  bool
}
type FileInputOptions struct {
	Label, Accept      string
	Multiple, Disabled bool
	MaxFiles           int
	MaxSize            int64
	Browse             func() ([]string, error)
	OnChange           func([]UploadFile)
}

func acceptsFile(name, contentType, accept string) bool {
	if strings.TrimSpace(accept) == "" {
		return true
	}
	ext := strings.ToLower(filepath.Ext(name))
	contentType = strings.ToLower(contentType)
	if mediaType, _, err := mime.ParseMediaType(contentType); err == nil {
		contentType = mediaType
	}
	for _, pattern := range strings.Split(strings.ToLower(accept), ",") {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		if pattern == ext || strings.HasPrefix(pattern, ".") && strings.HasSuffix(strings.ToLower(name), pattern) || pattern == contentType || strings.HasSuffix(pattern, "/*") && strings.HasPrefix(contentType, strings.TrimSuffix(pattern, "*")) {
			return true
		}
	}
	return false
}
func AddUploadPaths(state *FileInputState, paths []string, o FileInputOptions) bool {
	if state.Busy || o.Disabled {
		return false
	}
	original := slices.Clone(state.Files)
	var problems []string
	if !o.Multiple && len(paths) > 1 {
		problems = append(problems, "仅允许选择一个文件")
		paths = paths[:1]
	}
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			problems = append(problems, filepath.Base(path)+"：无法读取")
			continue
		}
		if !info.Mode().IsRegular() {
			problems = append(problems, info.Name()+"：请选择普通文件")
			continue
		}
		typeName := mime.TypeByExtension(filepath.Ext(info.Name()))
		if !acceptsFile(info.Name(), typeName, o.Accept) {
			problems = append(problems, info.Name()+"：文件类型不允许")
			continue
		}
		if o.MaxSize > 0 && info.Size() > o.MaxSize {
			problems = append(problems, info.Name()+"：文件过大")
			continue
		}
		if slices.ContainsFunc(state.Files, func(f UploadFile) bool { return f.Path == path }) {
			continue
		}
		if !o.Multiple {
			state.Files = nil
		} else if o.MaxFiles > 0 && len(state.Files) >= o.MaxFiles {
			problems = append(problems, "已达到文件数量上限")
			break
		}
		state.Files = append(state.Files, UploadFile{Path: path, Name: info.Name(), Size: info.Size(), MIME: typeName})
	}
	state.Error = strings.Join(problems, "；")
	changed := !slices.Equal(original, state.Files)
	if changed && o.OnChange != nil {
		o.OnChange(slices.Clone(state.Files))
	}
	return changed
}
func (k *UI) FileInput(state *FileInputState, options ...FileInputOptions) ui.Element {
	o := option(options)
	if o.Label == "" {
		o.Label = "选择文件"
	}
	disabled := o.Disabled || state.Busy
	root := ui.Column(k.c).Padding(18, 20).Gap(10).Radius(k.tokens.RadiusLarge).
		Border(1, k.tokens.Native.Border).BorderStyle(ui.BorderDashed).
		Background(k.tokens.Subtle.Alpha(0.45)).Role(ui.RoleGroup).Label(o.Label)
	if state.pending {
		root.SetChanged()
		state.pending = false
	}
	if !disabled && root.FileDragOver() {
		root.Background(k.tokens.InfoMuted).BorderColor(k.tokens.Info).BorderStyle(ui.BorderSolid)
	}
	if paths := func() []string {
		if disabled {
			return nil
		}
		return root.DroppedFiles()
	}(); len(paths) > 0 {
		if AddUploadPaths(state, paths, o) {
			root.SetChanged()
		}
	}
	root.Children(func() {
		ui.Row(k.c).AlignItems(ui.Center).Gap(12).Wrap().Children(func() {
			ui.Box(k.c).Size(36, 36).Radius(k.tokens.Radius).Center().Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).Children(func() {
				k.Icon(IconUpload).FontSize(16).TextColor(k.tokens.Muted)
			})
			ui.Column(k.c).Gap(2).Grow(1).MinWidth(160).Children(func() {
				k.Text(o.Label).FontWeight(600)
				k.Caption("拖放文件到这里，或使用系统文件选择器。")
			})
			k.Button("浏览文件", ButtonOptions{Variant: Outline, Disabled: disabled}).OnClick(func() {
				browse := o.Browse
				if browse == nil {
					browse = func() ([]string, error) {
						return mygo.Dialog.Open(mygo.OpenDialogOptions{Title: o.Label, Multiple: o.Multiple})
					}
				}
				paths, err := browse()
				if err != nil {
					state.Error = err.Error()
					return
				}
				if AddUploadPaths(state, paths, o) {
					state.pending = true
				}
			})
		})
		for _, file := range state.Files {
			ui.Row(k.c.Key(file.Path)).AlignItems(ui.Center).Padding(6, 10).Gap(10).Radius(k.tokens.RadiusSmall).
				Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).Children(func() {
				k.Caption(fmt.Sprintf("%s · %.1f KB", file.Name, float64(file.Size)/1024)).Grow(1)
				k.IconButton("移除文件"+file.Name, IconClose, ButtonOptions{Size: Small, Variant: Ghost, Disabled: disabled}).OnClick(func() {
					state.Files = slices.DeleteFunc(state.Files, func(f UploadFile) bool { return f.Path == file.Path })
					root.SetChanged()
					if o.OnChange != nil {
						o.OnChange(slices.Clone(state.Files))
					}
				})
			})
		}
		if state.Busy {
			k.Progress(state.Progress).FillWidth().Label("文件处理进度")
		}
		if state.Error != "" {
			root.Error(state.Error)
			k.Caption(state.Error).TextColor(k.tokens.Native.Danger)
		}
	})
	return root
}
