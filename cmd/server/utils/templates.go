package utils

import (
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"entity/interview/cmd/server/assets"
)

const templatesRoot = "templates"

var templateList map[string]*template.Template

func FindTemplate(w http.ResponseWriter, name string) *template.Template {
	tmpl, exists := templateList[name]
	if !exists {
		http.Error(w, "Template not found", http.StatusNotFound)
		return nil
	}
	return tmpl
}

func InitialiseTemplatesFromDir(templatefs fs.FS) error {
	templateList = make(map[string]*template.Template)

	entries, err := fs.ReadDir(templatefs, templatesRoot)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() {
			slog.Warn("Ignoring directory in templateList", "dir", name)
			continue
		}

		if name == "layout.tmpl" {
			slog.Info("Skipping layout.tmpl, will be parsed with subtemplates")
			continue
		}

		if strings.HasPrefix(name, "layout-") {
			base := path.Join(templatesRoot, "layout.tmpl")
			child := path.Join(templatesRoot, name)

			// Parse base + child together from the provided FS
			tmpl, err := template.ParseFS(templatefs, base, child)
			if err != nil {
				slog.Error("error parsing layout child", "file", name, "err", err)
				continue
			}

			key := strings.TrimPrefix(name, "layout-")
			templateList[key] = tmpl
			slog.Info("added root template", "file", name, "key", key)
			continue
		}

		// Standalone template
		file := path.Join(templatesRoot, name)
		tmpl, err := template.ParseFS(templatefs, file)
		if err != nil {
			slog.Error("Error parsing template", "file", name, "err", err)
			continue
		}
		templateList[name] = tmpl
	}

	return nil
}

func InitialiseTemplates() error {
	return InitialiseTemplatesFromDir(assets.TemplateFS)
	//return InitialiseTemplatesFromDir(os.DirFS("."))
}
