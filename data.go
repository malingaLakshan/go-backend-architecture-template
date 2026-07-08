// Config-based mode: serve the pre-loaded site config when set.
// Do not compare requested siteId here.
// QA may intentionally provide a mismatched config to test validation failure.
var config *site.SiteConfig
if h.siteConfig != nil {
	fmt.Printf("[INFO] Serving configured mock site config: %s\n", h.siteConfig.SiteID)
	config = h.siteConfig
} else {
	config = getMockSiteConfig(siteID)
	if config == nil {
		http.Error(w, fmt.Sprintf("site not found: %s", siteID), http.StatusNotFound)
		return
	}
}