package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type cursorFrontMatter struct {
	Description string `yaml:"description"`
	AlwaysApply bool   `yaml:"alwaysApply"`
}

type githubFrontMatter struct {
	ApplyTo     string   `yaml:"applyTo"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}

type agentMetadata struct {
	Description string   `yaml:"description"`
	AlwaysApply bool     `yaml:"alwaysApply"`
	Tags        []string `yaml:"tags,omitempty"`
}

type docFrontMatter struct {
	Title           string        `yaml:"title"`
	SidebarLabel    string        `yaml:"sidebar_label,omitempty"`
	SidebarPosition *int          `yaml:"sidebar_position,omitempty"`
	Tags            []string      `yaml:"tags,omitempty"`
	Agent           agentMetadata `yaml:"agent"`
}

func main() {
	var (
		cursorDir = flag.String("cursor", ".cursor/rules", "directory containing Cursor rule files")
		outputDir = flag.String("output", ".github/instructions/rules", "directory to write GitHub Copilot instruction files")
		applyTo   = flag.String("apply-to", "**", "glob pattern for applying GitHub instructions")
		docsDir   = flag.String("docs", "docs/content/rules", "directory to write canonical documentation files")
	)

	flag.Parse()

	if err := syncRules(*cursorDir, *outputDir, *applyTo, *docsDir); err != nil {
		fmt.Fprintf(os.Stderr, "sync failed: %v\n", err)
		os.Exit(1)
	}
}

func syncRules(cursorDir, outputDir, applyTo, docsDir string) error {
	entries, err := os.ReadDir(cursorDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("cursor directory %q not found", cursorDir)
		}
		return fmt.Errorf("read cursor dir: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		return fmt.Errorf("create docs dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".mdc" {
			continue
		}

		srcPath := filepath.Join(cursorDir, entry.Name())
		content, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", srcPath, err)
		}

		meta, body, err := splitFrontMatter(content)
		if err != nil {
			return fmt.Errorf("parse %s: %w", srcPath, err)
		}

		var cursorMeta cursorFrontMatter
		if err := yaml.Unmarshal(meta, &cursorMeta); err != nil {
			return fmt.Errorf("unmarshal front matter for %s: %w", srcPath, err)
		}

		ghMeta := githubFrontMatter{
			ApplyTo:     applyTo,
			Name:        strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())),
			Description: cursorMeta.Description,
		}

		if cursorMeta.AlwaysApply {
			ghMeta.Tags = append(ghMeta.Tags, "always-apply")
		}

		metaBytes, err := yaml.Marshal(&ghMeta)
		if err != nil {
			return fmt.Errorf("marshal GitHub front matter for %s: %w", srcPath, err)
		}

		outputFile := filepath.Join(outputDir, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".instructions.md")
		if err := writeInstruction(outputFile, metaBytes, body); err != nil {
			return fmt.Errorf("write %s: %w", outputFile, err)
		}

		title := deriveTitle(body)
		if title == "" {
			title = fallbackTitle(entry.Name())
		}

		docMeta := docFrontMatter{
			Title:        title,
			SidebarLabel: title,
			Tags:         []string{"ai-rules"},
			Agent: agentMetadata{
				Description: cursorMeta.Description,
				AlwaysApply: cursorMeta.AlwaysApply,
				Tags:        append([]string{}, ghMeta.Tags...),
			},
		}

		if pos := sidebarPositionFromName(entry.Name()); pos != nil {
			docMeta.SidebarPosition = pos
		}

		docPath := filepath.Join(docsDir, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".md")
		if err := writeDoc(docPath, docMeta, body); err != nil {
			return fmt.Errorf("write doc %s: %w", docPath, err)
		}
	}

	return nil
}

func splitFrontMatter(data []byte) ([]byte, []byte, error) {
	content := bytes.TrimLeft(data, "\ufeff\r\n")

	var newline []byte
	switch {
	case bytes.HasPrefix(content, []byte("---\r\n")):
		newline = []byte("\r\n")
		content = content[len("---\r\n"):]
	case bytes.HasPrefix(content, []byte("---\n")):
		newline = []byte("\n")
		content = content[len("---\n"):]
	default:
		return nil, nil, fmt.Errorf("missing front matter header")
	}

	delimiter := append(append([]byte{}, newline...), []byte("---")...)
	end := bytes.Index(content, delimiter)
	if end == -1 {
		return nil, nil, fmt.Errorf("closing front matter delimiter not found")
	}

	meta := bytes.TrimRight(content[:end], "\r\n")

	body := content[end+len(delimiter):]
	body = bytes.TrimLeft(body, "\r\n")

	return meta, body, nil
}

func writeInstruction(path string, meta, body []byte) error {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(meta)
	if !bytes.HasSuffix(meta, []byte("\n")) {
		buf.WriteByte('\n')
	}
	buf.WriteString("---\n\n")
	buf.Write(bytes.TrimLeft(body, "\r\n"))

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}

func writeDoc(path string, meta docFrontMatter, body []byte) error {
	metaBytes, err := yaml.Marshal(meta)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(metaBytes)
	if !bytes.HasSuffix(metaBytes, []byte("\n")) {
		buf.WriteByte('\n')
	}
	buf.WriteString("---\n\n")
	buf.Write(bytes.TrimLeft(body, "\r\n"))

	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func deriveTitle(body []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func fallbackTitle(filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.Split(base, "-")
	for i, part := range parts {
		part = strings.ReplaceAll(part, "_", " ")
		if part == "" {
			parts[i] = part
			continue
		}

		runes := []rune(part)
		runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
		parts[i] = string(runes)
	}
	return strings.Join(parts, " ")
}

func sidebarPositionFromName(name string) *int {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	parts := strings.SplitN(base, "-", 2)
	if len(parts) == 0 {
		return nil
	}

	val, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}

	return &val
}
