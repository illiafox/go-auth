package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"path/filepath"
	"reflect"
)

func Load(ts *Templates, root string) error {
	root = path.Dir(root)

	var tmpl Template

	tType := reflect.TypeOf(tmpl)

	reflectType := reflect.TypeOf(ts).Elem()
	reflectValue := reflect.ValueOf(ts).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		typeField := reflectType.Field(i)

		if typeField.Type != tType {
			return fmt.Errorf("field '%s' has type '%s', not '%s'", typeField.Name, typeField.Type, tType)
		}

		folder, ok := typeField.Tag.Lookup("tmpl")
		if !ok || folder == "" {
			return fmt.Errorf("field '%s' does not have 'tmpl' struct tag", typeField.Name)
		}

		var files []string

		err := filepath.WalkDir(root+"/"+folder, func(file string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("open '%s': %w", file, err)
			}

			if d.IsDir() {
				return nil
			}

			if path.Ext(file) != ".html" {
				return fmt.Errorf("%s : extention %s is not supported ('*.http' only)", path.Dir(file), path.Ext(file))
			}

			files = append(files, file)

			return nil
		})

		if err != nil {
			return fmt.Errorf("walk dir: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("no files found in: %w", err)
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("parse template files: %w", err)
		}

		tmpl.t = t

		reflectValue.Field(i).Set(reflect.ValueOf(tmpl))
	}

	return nil
}
