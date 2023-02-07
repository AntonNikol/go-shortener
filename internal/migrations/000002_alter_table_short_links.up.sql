CREATE UNIQUE INDEX IF NOT EXISTS full_url_unique
    ON short_links (full_url)