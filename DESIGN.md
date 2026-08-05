---
name: Groupie Tracker
description: Dark neon venue for exploring artists and concert dates.
colors:
  stage-cyan: "#22d3ee"
  stage-cyan-soft: "rgba(34, 211, 238, 0.10)"
  stage-cyan-line: "rgba(34, 211, 238, 0.40)"
  amp-purple: "#a855f7"
  house-pink: "#ec4899"
  house-black: "#030712"
  surface-glass: "rgba(255, 255, 255, 0.05)"
  surface-deep: "#111827"
  border-faint: "rgba(255, 255, 255, 0.05)"
  border-soft: "rgba(255, 255, 255, 0.10)"
  border-strong: "#1f2937"
  stage-white: "#ffffff"
  text-body: "#d1d5db"
  text-muted: "#9ca3af"
  text-dim: "#6b7280"
typography:
  display:
    fontFamily: "Space Grotesk, sans-serif"
    fontSize: "clamp(2.25rem, 6vw, 4.5rem)"
    fontWeight: 700
    lineHeight: 1.1
    letterSpacing: "-0.03em"
  headline:
    fontFamily: "Space Grotesk, sans-serif"
    fontSize: "clamp(1.5rem, 3vw, 1.875rem)"
    fontWeight: 600
    lineHeight: 1.2
  title:
    fontFamily: "Space Grotesk, sans-serif"
    fontSize: "1rem"
    fontWeight: 600
    lineHeight: 1.3
  body:
    fontFamily: "DM Sans, sans-serif"
    fontSize: "clamp(0.875rem, 1.5vw, 1.125rem)"
    fontWeight: 400
    lineHeight: 1.6
    letterSpacing: "0.01em"
  label:
    fontFamily: "DM Sans, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 600
    letterSpacing: "0.04em"
    textTransform: "uppercase"
rounded:
  sm: "8px"
  md: "12px"
  lg: "16px"
  full: "9999px"
spacing:
  xs: "8px"
  sm: "12px"
  md: "16px"
  lg: "24px"
  xl: "32px"
  container: "1280px"
components:
  button-primary:
    backgroundColor: "{colors.stage-cyan-soft}"
    textColor: "{colors.stage-cyan}"
    typography: "{typography.title}"
    rounded: "{rounded.lg}"
    padding: "16px 32px"
  button-primary-hover:
    backgroundColor: "rgba(34, 211, 238, 0.20)"
    textColor: "{colors.stage-cyan}"
    rounded: "{rounded.lg}"
    padding: "16px 32px"
  button-solid:
    backgroundColor: "{colors.stage-white}"
    textColor: "{colors.house-black}"
    typography: "{typography.title}"
    rounded: "{rounded.md}"
    padding: "12px 24px"
  artist-card:
    backgroundColor: "{colors.surface-glass}"
    textColor: "{colors.stage-white}"
    rounded: "{rounded.md}"
    padding: "12px 16px"
  chip:
    backgroundColor: "rgba(17, 24, 39, 0.50)"
    textColor: "{colors.text-body}"
    typography: "{typography.body}"
    rounded: "{rounded.sm}"
    padding: "6px 12px"
  nav-link-active:
    backgroundColor: "{colors.surface-glass}"
    textColor: "{colors.stage-cyan}"
    typography: "{typography.body}"
    rounded: "{rounded.sm}"
    padding: "8px 16px"
---

# Design System: Groupie Tracker

## Overview

**Creative North Star: "The Midnight Venue"**

Groupie Tracker is a dim concert venue after the doors open and the house lights stay low. A visitor walks into a black chamber lit only by faint neon glow pools overhead, sees every act that plays here hanging on the walls, reads who is in the building tonight, and steps up to the one bright door that opens onto the floor. Every screen is a room in that venue: the landing is the entry hall with the acts on the wall, the artist grid is the floor where the acts are lined up, and the tour-dates page is the schedule board by the stage.

The system is image-first and chrome-last. Artist photography carries the emotional load; interface furniture is translucent glass with hairline borders that recede into the black. Density is loose and breathable — wide containers, generous vertical rhythm, single-column narratives. Color is spent almost entirely on three neon accents that behave like stage lights: they guide an action or mark live data, never purely decorate. There is no light source brighter than the text on the door; the interface only glows under the cursor.

The world is genuinely dark, not "dark mode." The base is near-black (`#030712`), surfaces are translucent glass, and the only full-opacity fills are the white door button, the white headline, and the artists' photographs. Hover is where the room wakes up: glass edges tint neon, card halos bloom, and images scale. Motion is authored one moment at a time — an entrance stagger here, a drifting wall there — never a running animation show.

**Key Characteristics:**
- A black chamber with faint neon glow pools; imagery carries the light.
- Translucent glass surfaces with 1px hairlines; no resting shadows.
- Three neon accents with strict jobs: cyan = the action, purple = secondary data, pink = tertiary hover.
- Space Grotesk makes the statements; DM Sans is the quiet house font.
- One authored motion moment per screen; reduced-motion is honored everywhere.
- Square poster-frames for artist imagery; hover is where things glow.

## Colors

A black chamber lit by three neon stage lights and faint ambient glow pools.

### Primary
- **Stage Cyan** (`#22d3ee`): the action color. The one bright door on the landing, active nav link, CTA buttons, scrollbar thumb, live pulse dots, and the primary stat in the stats band. When cyan appears, it means "do this" or "this is live."

### Secondary
- **Amp Purple** (`#a855f7`): secondary data and secondary accents. The concert-date count, the tour-timeline gradient start, the brand tile gradient, member-tag hovers, the small dot beside artist names in the grid, and the footer GitHub hover.

### Tertiary
- **House Pink** (`#ec4899`): the softest, rarest accent. Reserved for back-link hovers, date-chip hovers on the tour page, and the tour-timeline gradient end.

### Neutral
- **House Black** (`#030712`): the page background, underpinned by three faint radial glow pools — cyan high center, purple on the right, pink on the left — at 4-6% opacity. Never pure black.
- **Glass White** (`rgba(255, 255, 255, 0.05)`): the translucent surface for nav bar, cards, hero badge, icon tiles, and timeline event tiles.
- **Deep Pit** (`#111827`): solid chip/tag backgrounds, at 50% opacity (`gray-900/50`).
- **Hairline White** (`rgba(255, 255, 255, 0.05-0.10)`): 1px borders on glass surfaces.
- **Iron** (`#1f2937`): chip and image borders.
- **Stage White** (`#ffffff`): headlines, card names, stat numbers, and the solid door button text.
- **Body Ash** (`#d1d5db`): body copy and secondary text.
- **Muted Ash** (`#9ca3af`): labels, meta lines, icon tints, and secondary stat labels.
- **Dim Ash** (`#6b7280`): footer-only text; never used for content-bearing copy.

### Named Rules
**The Neon Rarity Rule.** Neon accents are attention, not decoration. Each accent has exactly one job — cyan signals the action, purple signals secondary live data, pink signals a tertiary hover — and they should never pile up on one screen. If a control needs more color than one accent, it's the wrong control.

**The Glass Border Rule.** Every surface is translucent glass over the black room, defined by a hairline border and what sits behind it. No surface is opaque, and no surface is borderless on a busy background.

## Typography

**Display Font:** Space Grotesk (500/600/700), fallback `sans-serif`.
**Body Font:** DM Sans (400/500/600), fallback `sans-serif`.

**Character:** Space Grotesk is the venue's printed-poster type — geometric, open, confident, with wide, friendly letterforms that read as announcement. DM Sans is the quiet house font: neutral, legible, humanist, doing the talking without drawing attention. The pairing is announcement versus housekeeping, and they never swap roles.

### Hierarchy
- **Display** (700, `clamp(2.25rem, 6vw, 4.5rem)`, line-height 1.1, tracking -0.03em): the first-viewport statement on the landing — "The whole scene is on the wall." Appears once per surface at most.
- **Headline** (600, `clamp(1.5rem, 3vw, 1.875rem)`, line-height 1.2): page-level titles on the artists grid and section titles.
- **Title** (600, 1rem, line-height 1.3): artist card names, button labels, brand wordmark.
- **Body** (400/500, `clamp(0.875rem, 1.5vw, 1.125rem)`, line-height 1.6, tracking 0.01em): all copy. Long-form blocks cap around 65-75ch (`max-w-3xl`).
- **Label** (600, 0.75rem, tracking 0.04em, uppercase): eyebrow lines, stat labels, hero badge text ("52 acts · one room").

### Named Rules
**The Grotesk Over the Glass Rule.** Space Grotesk only ever carries a statement — headlines, card names, stat numerals, the door. DM Sans carries everything else. Never set a paragraph, a caption, or a meta line in Space Grotesk; never set a headline in DM Sans.

**The Static Tracking Rule.** Track numerals tightly (`-0.03em`) so stat numbers read as a single mass, uppercase labels wide (`+0.04em`), and body copy just open enough (`+0.01em`). Tracking is part of the voice; don't flatten it.

## Layout

A single-column flow inside a `max-w-7xl` container (1280px) with `px-4 sm:px-6 lg:px-8` gutters. Content blocks center horizontally; the only full-bleed exceptions are the landing wall and the nav bar.

- **Artist grid:** 2 columns on mobile → 3 at 640px → 4 at 1024px → 5 at 1280px. Square image up top, name + year below.
- **Landing wall:** 4 columns on mobile → 6 at 640px → 8 at 1024px, full-bleed with the image grid bleeding off every edge (`inset: -15% 0`).
- **Stats band:** two cells divided by a `gap-px` hairline over a glass backing — the divider IS the gap.
- **Tour timeline:** a single left-rail column; event tiles sit right of a gradient rail.
- **Vertical rhythm:** page heroes `pt-16 sm:pt-24`, sections `pb-12 sm:pb-20`, interior padding `py-8`/`py-12`. More space above a heading than below it.

Density is loose. Cards breathe; grids use `gap-3` to `gap-4`; there is no masonry, no multi-column prose, no text crushing against a card edge.

## Elevation & Depth

The system is **tonal layering with glow as a state response — never resting shadows.** Depth comes from three stacked layers: the black room, the translucent glass surfaces, and the light that only appears under the cursor.

- The room itself is lit by the three body glow pools (cyan/purple/pink radials at 4-6% over `#030712`).
- Glass surfaces sit one hairline above the room; anything further forward is conveyed by glow, not shadow.
- The only depth at rest is the landing's vignette, which dims the wall around the centered statement.

### Named Rules
**The No-Resting-Shadow Rule.** No surface casts a shadow while idle. On hover or focus, a glow takes its place: cards bloom a cyan→purple→pink gradient halo (`card-glow`), artist images gain a cyan border and scale, the door's cyan border intensifies. If a surface needs to feel raised while nothing is touching it, it is missing its hairline or its glass fill.

## Shapes

Gently curved corners everywhere, sized by role: **8px** (`rounded-lg`) for chips, tags, nav links, and the brand tile; **12px** (`rounded-xl`) for cards, buttons, icon tiles, and images; **16px** (`rounded-2xl`) for the primary door and the stats band; full pills for badges and pulse dots.

Borders are 1px hairlines — `white/5` to `white/10` on glass, `gray-800` on chips and images. The recurring silhouette is the **square poster-frame**: artist photos render `aspect-square` with `object-cover`, so every act hangs in the same flat frame regardless of source aspect ratio.

## Components

### Buttons
- **Shape:** 12px radius for solid and ghost buttons (16px for the primary door); labels in Space Grotesk 600.
- **Primary — The Neon Door:** glass cyan — `bg` stage-cyan at 10%, 1px border at 40%, text stage-cyan, padding 16px 32px. Hover raises the fill to 20% and the border to 70%; the arrow translates 4px right. Focus shows a cyan ring with a `gray-950` offset. During an HTMX swap the label is replaced by a spinner.
- **Solid — Explore:** white fill, House Black text, padding 12px 24px, hover `gray-100`. The one opaque, high-contrast action; used for the primary jump into the grid.
- **Ghost:** the same glass-cyan recipe at smaller scale (12px 24px, 12px radius, 30% border) for secondary actions like the tour-dates CTA.

### Chips
- **Style:** `gray-900/50` fill, `gray-800` hairline, 8px radius, 6px 12px padding, Body Ash text.
- **State:** hover lifts the border and text into an accent — cyan/purple/pink cycling per third on member tags, pink on tour-date chips. Used for member lists and date pills; the pink hover is a deliberate tertiary Easter egg.

### Cards / Containers
- **Corner Style:** 12px radius.
- **Background:** Glass White (`white/5`).
- **Shadow Strategy:** none at rest; `card-glow` halo (linear-gradient 135° cyan→purple→pink at 30% alpha, `inset: -1px`) fades in on hover.
- **Border:** 1px hairline `white/5`.
- **Internal Padding:** 12-16px (`p-3 sm:p-4`).
- **Artist card behavior:** square image scales to 1.1 on hover (0.5s ease-out), card lifts 2px, name tints cyan. `view-transition-name` per artist id so cards morph into detail pages.

### Inputs / Fields
- Search and sort controls currently render as unstyled browser defaults (a known gap outside the landing pass). When styled, they should inherit the glass fill, `gray-800` hairline, 8px radius, and cyan focus treatment of the rest of the system rather than introducing a new field language.

### Navigation
- Fixed translucent glass bar. Brand is a cyan→purple gradient tile (8px radius) with a Space Grotesk 700 wordmark that tints cyan on hover. The active link is cyan text on glass; idle links are Muted Ash that brighten to white with a glass fill on hover. Mobile collapses to a hamburger opening a stacked glass sheet.

### Signature Components
- **The Venue Wall** (landing): a full-bleed grid of every artist's photo at 14% opacity and 75% saturation, drifting slowly (48s ease-in-out alternate) like the room is breathing, under a radial+linear vignette. Hover focuses the tile under the cursor — opacity to 50%, scale to 1.1, saturation to 1 — and the wall dims. Under `prefers-reduced-motion`, the drift is dropped and only the entrance remains.
- **The House Board** (landing): the live stats band at the bottom of the viewport — artist count and concert-date count in display numerals, uppercase label lines, glass cells separated by a hairline gap.
- **The Tour Timeline:** a gradient rail (`cyan → purple → pink`), a 4px cyan dot per event, glass event tiles that tint their border cyan on hover.

## Do's and Don'ts

### Do:
- **Do** keep artist imagery in square poster-frames (`aspect-square`, `object-cover`) and let it carry the screen.
- **Do** use Stage Cyan for exactly one action per screen; it means "do this."
- **Do** set every statement in Space Grotesk and every paragraph in DM Sans.
- **Do** build surfaces from translucent glass (`white/5`) plus a 1px hairline.
- **Do** author one motion moment per screen and kill it under `prefers-reduced-motion`.
- **Do** track stat numerals at `-0.03em` and uppercase labels at `+0.04em`.
- **Do** keep the room black (`#030712`) and let the three glow pools and the imagery provide the light.

### Don't:
- **Don't** use resting shadows to raise a surface; depth is borders, glass, and hover glow.
- **Don't** stack two neon accents on one control — pick the cyan action or the quiet ghost, not both.
- **Don't** set body copy in Space Grotesk or a headline in DM Sans.
- **Don't** make a surface opaque unless it's the white door button or an artist photograph.
- **Don't** add motion that plays unattended; no running carousels, no looping float on every screen, no entrance that repeats.
- **Don't** crop artist imagery to a ratio other than square in cards, or brighten the wall so much that the landing stops reading as one dim room.
