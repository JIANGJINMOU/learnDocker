package images

import (
	"archive/tar"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"example.com/containeredu/internal/paths"
)

func TestImportDockerSaveTar(t *testing.T) {
	tmp := t.TempDir()
	tarPath := filepath.Join(tmp, "img.tar")
	f, err := os.Create(tarPath)
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewWriter(f)
	// create layer tar inside temp layout
	layerDir := filepath.Join(tmp, "layer")
	if err := os.MkdirAll(layerDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(layerDir, "hello.txt"), []byte("world"), 0o644); err != nil {
		t.Fatal(err)
	}
	layerTar := filepath.Join(tmp, "layer.tar")
	lf, err := os.Create(layerTar)
	if err != nil {
		t.Fatal(err)
	}
	lw := tar.NewWriter(lf)
	// add file to layer tar
	addFileToTar(t, lw, "hello.txt", []byte("world"))
	lw.Close()
	lf.Close()
	// write layer tar into main tar under name
	layerName := "abcdef/layer.tar"
	writeFileToTar(t, tr, layerName, layerTar)
	// write manifest.json
	manifest := []ManifestEntry{{
		Config:   "config.json",
		RepoTags: []string{"test:latest"},
		Layers:   []string{layerName},
	}}
	mb, _ := json.Marshal(manifest)
	addFileToTar(t, tr, "manifest.json", mb)
	tr.Close()
	f.Close()

	name := "testimage"
	if err := ImportDockerSaveTar(tarPath, name); err != nil {
		t.Fatalf("import error: %v", err)
	}
	// check extraction
	imgRoot := filepath.Join(paths.ImagesRoot(), name)
	_, err = os.Stat(filepath.Join(imgRoot, "metadata.json"))
	if err != nil {
		t.Fatalf("metadata not found: %v", err)
	}
	// layer content exists
	entries, _ := os.ReadDir(filepath.Join(imgRoot, "layers"))
	if len(entries) == 0 {
		t.Fatalf("no layers extracted")
	}
	// check file
	target := filepath.Join(imgRoot, "layers", entries[0].Name(), "hello.txt")
	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("hello.txt not found: %v", err)
	}
	if string(b) != "world" {
		t.Fatalf("unexpected content: %q", string(b))
	}
}

func TestImportDockerSaveTarMissingManifest(t *testing.T) {
	tmp := t.TempDir()
	tarPath := filepath.Join(tmp, "bad.tar")
	f, err := os.Create(tarPath)
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewWriter(f)
	addFileToTar(t, tr, "somefile", []byte("x"))
	tr.Close()
	f.Close()
	err = ImportDockerSaveTar(tarPath, "x")
	if err == nil {
		t.Fatalf("expected error for missing manifest")
	}
}

func addFileToTar(t *testing.T, tw *tar.Writer, name string, content []byte) {
	t.Helper()
	hdr := &tar.Header{
		Name: name,
		Mode: 0o644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
}

func writeFileToTar(t *testing.T, tw *tar.Writer, name string, src string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	addFileToTar(t, tw, name, data)
}
