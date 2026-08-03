# Surface Brief: Home Page

## Job and Audience

**Audience:** Music fans arriving at the site to browse and discover artists. They come with curiosity — no account, no search query, just a desire to see who's here and explore.

**Mode:** Persuade (the page earns further exploration) / Operate (browsing the grid is the primary task).

## Outcome and Proof

The visitor must quickly understand this is a music artist platform, feel drawn to browse the grid, and click into an artist with confidence. Success: visitors scan cards, pause on ones that interest them, and click.

**Unique truth:** Unlike a generic gallery, every card represents a real music act with a distinct visual identity (their image), a formation year, and tour activity — the grid itself tells a story of active artists.

## Selected Direction

**Visual authority:** A modern music-streaming / Bandcamp / Spotify-inspired dark aesthetic — the grid feels like browsing album art in a record store, not a database listing. The existing dark-neon identity (gray-950, neon cyan/purple/pink) is the right foundation; it needs intensity and polish, not replacement.

**Structural thesis:** The hero and stats are a brief landing, then the grid is the page. Cards must feel like album covers — the image dominates, information recedes into subtle chroma. Hover reveals depth (gradient overlay, slight scale, accent border).

**Sequence:** Hero → stats → grid. No pagination or load-more (all artists visible on one scroll).

**Focal moment:** The grid at full width on desktop — a wall of album-art-sized cards that makes the visitor want to browse.

## Scope and Boundaries

- Frontend only: no changes to `cmd/server.go`, `internal/config/`, `internal/handlers/`, `internal/routes/`
- Only `templates/index.html`, `templates/base.html`, and `static/style.css` may be modified
- No new features: no search, no filtering, no sorting, no pagination
- No new JavaScript or CDN dependencies beyond existing (Tailwind via CDN script, HTMX 2.x, Google Fonts)
- Maintain all current HTMX behavior (card click → swap artist detail into view)
- Maintain mobile responsiveness (current hamburger menu, grid breakpoints)

## States and Ranges

- **Data range:** 52 artists currently; cards must look good at any count between 1 and ~100
- **Loading state:** Skeleton loader per card (already implemented via `hx-indicator`)
- **Empty state:** "No artists found" message (already exists)
- **Error state:** If API fails at startup, the server crashes — no runtime error handling needed
- **Image failure:** Cards should degrade gracefully if an artist image URL fails to load

## Interaction and Layout

- **Grid:** responsive 1→2→3→4 columns (sm→md→lg→xl)
- **Card:** square aspect ratio, image fills the card, info overlays or sits below
- **Hover:** image scale-up + gradient overlay + accent border glow
- **Stats section:** two compact stat cards (total artists, total concert dates)
- **Typography:** Space Grotesk headings, DM Sans body — already configured
- **Color:** dark foundation (bg-gray-950, gray-900 cards), neon accents for interactive states and data highlights

## Constraints

- Tailwind CSS loaded from CDN script (no build step, no purge — all classes must be in the HTML)
- No custom fonts beyond Space Grotesk + DM Sans (already loaded from Google Fonts)
- HTMX for all dynamic behavior — no vanilla JS beyond the mobile menu toggle
- No backend changes — `config.ArtistByID` data structure stays as-is
- Template must remain valid Go `html/template` syntax
