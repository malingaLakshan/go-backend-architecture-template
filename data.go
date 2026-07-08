The GET /sites/{siteId} check proves mock-server is not serving the file config.

I started mock-server with fail_config.json where:
mock_mode = "file"
mock_site_file = "configs/wrong_site_config.json"

But this command:
Invoke-RestMethod http://localhost:8080/sites/b3489888-aacf-4451-893c-d7d994240f93 | ConvertTo-Json -Depth 10

still returns the real Bentonville config, not:
{
  "id": "WRONG-SITE-ID",
  "name": "Wrong Site"
}

Please fix the config wiring.

Expected:
1. runMockServer must load configs/fail_config.json.
2. If mock_mode == "file":
   - read mock_site_file path from config
   - unmarshal it into site.SiteConfig
   - pass that config into mocktarget.StartServerWithConfig
3. Handler must store that config.
4. HandleGetSite must return Handler.siteConfig directly when it is not nil.
5. It must not call getMockSiteConfig fallback when Handler.siteConfig exists.

Please add debug logs:
- Loaded config file: <config path>
- Mock mode: file
- Mock site file: configs/wrong_site_config.json
- Serving configured mock site config: WRONG-SITE-ID

Only change mock-server config wiring. Do not change validate or play.