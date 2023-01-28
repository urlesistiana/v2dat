package unpack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urlesistiana/v2dat/v2data"
	"go.uber.org/zap"
)

func newGeoSiteCmd() *cobra.Command {
	args := new(UnpackArgs)
	c := &cobra.Command{
		Use:   "geosite",
		Args:  cobra.ExactArgs(1),
		Short: "Unpack geosite file to text files.",
		Run: func(cmd *cobra.Command, a []string) {
			args.file = a[0]
			if err := unpackGeoSite(args); err != nil {
				logger.Fatal("failed to unpack geosite", zap.Error(err))
			}
		},
		DisableFlagsInUseLine: true,
	}
	c.Flags().BoolVarP(&args.with_type_prefix, "with_type_prefix", "p", false, "with type prefix geosite")
	c.Flags().StringVarP(&args.outDir, "out", "o", "", "output dir")
	c.Flags().StringArrayVarP(&args.filters, "filter", "f", nil, "unpack given tag and attrs")
	return c
}

func unpackGeoSite(args *UnpackArgs) error {
	fmt.Println(args.with_type_prefix)
	filePath, suffixes, outDir := args.file, args.filters, args.outDir
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	geoSiteList, err := v2data.LoadGeoSiteList(b)
	if err != nil {
		return err
	}

	entries := make(map[string][]*v2data.Domain)
	for _, geoSite := range geoSiteList.GetEntry() {
		tag := strings.ToLower(geoSite.GetCountryCode())
		entries[tag] = geoSite.GetDomain()
	}

	save := func(suffix string, data []*v2data.Domain) error {
		file := unpackPath(fileName(filePath), suffix, args.with_type_prefix)
		if len(outDir) > 0 {
			file = filepath.Join(outDir, file)
		}
		logger.Info(
			"unpacking entry",
			zap.String("tag", suffix),
			zap.Int("length", len(data)),
			zap.String("file", file),
		)
		return convertV2DomainToTextFile(data, file)
	}

	if len(suffixes) > 0 {
		for _, suffix := range suffixes {
			tag, attrs := splitAttrs(suffix)
			entry, ok := entries[tag]
			if !ok {
				return fmt.Errorf("cannot find entry %s", tag)
			}
			entry = filterAttrs(entry, attrs)
			if err := save(suffix, entry); err != nil {
				return fmt.Errorf("failed to save %s, %w", suffix, err)
			}
		}
	} else { // If tag is omitted, unpack all tags.
		for tag, domains := range entries {
			if err := save(tag, domains); err != nil {
				return fmt.Errorf("failed to save %s, %w", tag, err)
			}
		}
	}
	return nil
}

func convertV2DomainToTextFile(domain []*v2data.Domain, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return convertV2DomainToText(domain, f)
}

func convertV2DomainToText(domain []*v2data.Domain, w io.Writer) error {
	bw := bufio.NewWriter(w)
	for _, r := range domain {
		var prefix string
		switch r.Type {
		case v2data.Domain_Plain:
			prefix = "keyword:"
		case v2data.Domain_Regex:
			prefix = "regexp:"
		case v2data.Domain_Domain:
			prefix = ""
		case v2data.Domain_Full:
			prefix = "full:"
		default:
			return fmt.Errorf("invalid domain type %d", r.Type)
		}
		if _, err := bw.WriteString(prefix); err != nil {
			return err
		}
		if _, err := bw.WriteString(r.Value); err != nil {
			return err
		}
		if _, err := bw.WriteRune('\n'); err != nil {
			return err
		}
	}
	return bw.Flush()
}
