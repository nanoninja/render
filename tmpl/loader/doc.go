// Package loader provides different template loading strategies for the tmpl package.
// It offers various implementations for loading templates from different sources such as
// the filesystem, embedded files, or directly from strings in memory.
//
// The package supports multiple loader types:
//
//   - FSLoader: Loads templates from the local filesystem
//   - EmbedLoader: Loads templates from Go's embed.FS
//   - StringLoader: Loads templates defined as strings in code
//   - CompositeLoader: Combines multiple loaders with priority ordering
//
// Each loader implements the tmpl.Loader interface, providing a consistent way
// to discover and load templates regardless of their source.
//
// Basic usage with filesystem:
//
//	src := loader.NewFS(tmpl.LoaderConfig{
//	    Root: "templates",
//	    Extension: ".html",
//	})
//
// Using string-based templates:
//
//	templates := map[string]string{
//	    "welcome.html": "Hello {{ .Name }}",
//	    "layout.html": "<body>{{ template \"content\" . }}</body>",
//	}
//	src := loader.NewString(templates, tmpl.LoaderConfig{
//	    Extension: ".html",
//	})
//
// Combining multiple loaders:
//
//	composite := loader.NewComposite(
//	    []tmpl.Loader{customLoader, defaultLoader},
//	    tmpl.LoaderConfig{Extension: ".html"},
//	)
//
// The package handles common concerns like:
//   - Template discovery and loading
//   - Path manipulation and security
//   - Extension filtering
//   - Template prioritization (with CompositeLoader)
package loader
