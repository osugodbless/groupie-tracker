# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

**Primary: Music fans** browsing artist and concert data — visitors who want to discover bands, explore their members, founding years, first albums, and find out where and when they perform.

**Secondary: Evaluators** (01 Edu / Zone 01 staff) who assess the technical implementation — Go stdlib backend, HTMX integration, and frontend design. The product serves both audiences without compromising either.

## Product Purpose

Groupie Tracker is a browsable artist and concert-data explorer. It fetches band/artist profiles and tour schedules from the Groupie Trackers public API and presents them in a fast, visually engaging single-page interface. Success means visitors can browse any artist, learn about them, and quickly see where they've performed — without clicking through heavy pages or encountering visual noise.

## Positioning

A data-rich artist explorer that loads everything through partial-page swaps (HTMX) and a zero-framework Go backend — meaning instant navigation without a JavaScript framework or full-page reloads. The combination of server-rendered Go templates and HTMX gives music fans a fluid browsing experience that the category's typical JSON-SPA or static-site alternatives do not.

## Operating Context

- Visitors arrive via a browser (desktop or mobile) and interact through clicks and hovers — no accounts, logins, or persistent state.
- All data comes from the external Groupie Trackers API (`https://groupietrackers.herokuapp.com/api`). If the API is unreachable at startup, the server fatally exits — there is no offline or cached fallback.
- Navigation is entirely GET-based. Artist detail and tour-date views are loaded via HTMX `hx-get` and swap into the existing page (no full-page navigation after the initial load).
- The app runs locally at `http://0.0.0.0:8080` during development and can be deployed behind any Go-compatible host.

## Capabilities and Constraints

**Confirmed capabilities:**
- Browse all artists in a responsive grid with images, names, and creation years
- View artist detail pages (members, first album, creation year)
- View tour dates organized by location in a timeline layout
- Aggregate statistics (total artists, total concert dates) displayed on the home page
- Mobile-responsive layout with a hamburger menu
- Dark-mode-only visual theme

**Confirmed constraints:**
- Zero third-party Go dependencies — backend uses only Go standard library (`net/http`, `html/template`, `encoding/json`, `sync`)
- Frontend dependencies loaded from CDN: Tailwind CSS (v3 via CDN script), HTMX 2.x, Google Fonts (Space Grotesk + DM Sans)
- API fetches at startup only — data is static for the lifetime of the server process
- 15-second HTTP timeout on external API calls
- Templates parsed at startup via `template.Must` — syntax errors crash immediately

**Undecided:**
- Search / filter functionality
- Sorting or categorization of artists
- Caching or offline fallback strategy
- Deployment target (if any beyond localhost)
- Accessibility standard (WCAG level, if any)

## Brand Commitments

- **Name:** Groupie Tracker (confirmed in code, templates, and README)
- **Visual identity:** Dark theme (bg-gray-950), neon accent palette (cyan `#22d3ee`, purple `#a855f7`, pink `#ec4899`), rounded cards, border-heavy glassmorphic surfaces
- **Typography:** Space Grotesk for headings, DM Sans for body text
- **Tone:** Clean, modern, music-oriented — uses emoji as decorative elements (🎵, 🎤, 📅, 💿, 👥, 🗺️, 📍)
- **Repository:** `github.com/osugodbless/groupie-tracker`

## Evidence on Hand

- Artist data from the Groupie Trackers API (52 artists, includes images, member lists, creation years, first albums)
- Tour date data for all artists with location-to-dates mapping
- No testimonials, case studies, or usage analytics exist

## Product Principles

1. **Fluid browsing, not page-hopping.** Every interaction should feel like the same page updating, not a new page loading. HTMX partial-swaps are the mechanism; keep them fast and predictable.
2. **Visual clarity over density.** Artist cards, detail sections, and the tour timeline should feel spacious and scannable. The dark theme provides atmosphere; the neon accents guide attention.
3. **Data completeness without overwhelm.** Show everything available about an artist, but layer it — overview first, details on click, tour dates on request — so visitors self-serve their depth.
4. **Zero-framework simplicity on the backend.** Go stdlib only. Resist the temptation to add routing, templating, or HTTP libraries. The constraint is the differentiator.
5. **Mobile-first responsive.** The grid, timeline, navigation, and detail layouts must work on small screens before large ones. Mobile visitors are primary.

## Accessibility & Inclusion

No product-specific accessibility standard was established. The current implementation has no explicit a11y work (no ARIA labels, no focus management, no skip links, no color-contrast verification beyond what Tailwind's default palette provides).
