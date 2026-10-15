# Dwaar — marketing site

A self-contained, dependency-free landing page for [Dwaar](../README.md).

## Files

- `index.html` — the English page (default)
- `mr.html` — the Marathi page
- `styles.css` — styling (Catppuccin-inspired dark palette) and animations
- `script.js` — reveals, counters, live ticker, nav, copy buttons, demo tabs, theme toggle
- `favicon.svg` — logo mark

## Languages

`index.html` (English) and `mr.html` (मराठी) are two full pages that share the
same stylesheet and script. Each links to the other through the **EN / मराठी**
switcher in the nav, the footer, and `hreflang` alternate tags in the head.

To add another language, copy `mr.html`, translate the visible text, keep the
existing `id`/`class`/`data-*` attributes intact, and add a matching
`[data-copy]` `data-copied-label` plus `hreflang` links.

## Animations

All motion is progressive enhancement:

- Scroll reveals with automatic stagger, via `IntersectionObserver`
- Counting stat numbers and a live byte ticker in the hero terminal
- Drifting hero glow, blinking cursor, pulsing running status
- Nav shadow, scroll-progress bar, hover micro-interactions

Two conventions keep this robust:

- Elements are hidden **only while `:not(.is-visible)`**, so the reveal never
  sets a `transform` on a visible element and component hovers such as
  `.card:hover { transform: translateY(-3px) }` keep working.
- `script.js` also runs a viewport pass on load and on scroll, so content is
  revealed deterministically even if the observer has not fired yet, and
  anything above the viewport after a restored scroll is shown immediately.

Content is fully visible without JavaScript (reveal styles only apply once
`js` is added to `<html>`), and every effect is disabled under
`prefers-reduced-motion: reduce`.

## Preview locally

Any static file server works. For example:

```bash
cd site
python3 -m http.server 8080
# then open http://localhost:8080
```

## Deploy

The site has no build step, so it can be served as-is by any static host
(GitHub Pages, Netlify, Cloudflare Pages, S3, …).

For GitHub Pages, publish the `site/` directory from your deployment workflow
or point Pages at the `site` folder via a `gh-pages` branch.
