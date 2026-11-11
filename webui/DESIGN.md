# Relique web UI design system

The Relique web UI is a tool for checking backup jobs and acting on them. Every screen should answer "is my data safe, and what needs attention?" quickly. This document sets the visual rules. Tokens live in [`src/index.css`](src/index.css), and the reusable building blocks live in [`src/components/`](src/components/).

Only a light theme exists for now. A dark theme will reuse the same token names.

## Principles

1. **Readability first.** Text is the interface. Choose the font size, weight and color before adding any container, border or badge.
2. **Minimal and functional.** Every visual element must carry information. No decorative shadows, gradients, illustrations or oversized icons.
3. **Typography over chrome.** Use headings to group content and color to give meaning. Use bold to show hierarchy. Blocks are separated by one flat device, a tinted tile, and never by borders, shadows or badges.
4. **Purposeful whitespace.** Spacing is tight by default. Add space only to separate groups of information that differ in meaning.

## Tokens

### Color

| Token | Value | Tailwind / daisyUI | Use |
|---|---|---|---|
| base-100 | `#FFFFFF` | `bg-base-100` | Page background |
| base-200 | `#F7F6F3` | `bg-base-200` | Tiles (section background), sidebar |
| base-300 | `#E6E3DD` | `border-base-300` | The **only** border and divider color |
| base-content | `#1F1E1B` | `text-base-content` | Primary text |
| muted | `#6B6760` | `text-muted` | Labels, metadata, secondary text, empty values |
| primary | `#E8A317` | `bg-primary`, `btn-primary` | Amber accent. Use it as a **fill only**: primary button, active nav marker, logo |
| primary-content | `#1F1E1B` | `text-primary-content` | Text on amber. Amber is too light for white text |
| primary-strong | `#8F5400` | `text-link` | Links and any accent-colored **text** |
| accent | `#FDF3DC` | `bg-accent` | Tinted background for a selected or highlighted item, if needed |
| info | `#2563A8` | `text-info` | Status: running / informational |
| success | `#23753D` | `text-success` | Status: success / reachable |
| warning | `#C2410C` | `text-warning` | Status: incomplete / warning. Deliberately orange-red so it never reads as the amber accent |
| error | `#C42B25` | `text-error` | Status: error / unreachable / destructive |

Every text color above has a contrast ratio of at least 4.5:1 on `base-100`, `base-200` and `accent` (WCAG AA). Recheck the ratio whenever a value changes.

**Rules**
- Amber means "interactive or current". It never signals a status.
- Status colors only signal status, and always come with a text label or icon. Color is never the only signal.
- A count of zero is shown in `text-muted`, not in its status color. Only non-zero values draw attention.
- Show booleans as plain "Yes" or "No" (with "No" muted) using `<YesNo>`. Do not use colored badges.
- Never use raw Tailwind palette colors (`slate-*`, `gray-*`, `red-500`, `bg-white`, etc.). Use tokens only.

### Typography

There are three families, all bundled through `@fontsource` (no CDN), so the UI also works on offline servers:

- **Bricolage Grotesque** (variable, optical size axis), available as `font-display`. Used for headings, stat numbers and the wordmark. It gives the UI its character, so keep it to display sizes and never use it for body text, labels or tables.
- **IBM Plex Sans** (400/500/600), the default `font-sans`. Used for all UI text.
- **IBM Plex Mono** (400/500), used through `code`. Reserved for technical data.

| Role | Classes | Example |
|---|---|---|
| Page title (`h1`) | `font-display text-3xl font-bold tracking-tight` | The client name, "Backup job a4f0e578" |
| Section heading (`h2`) | `font-display text-xl font-semibold tracking-tight` | "Connection", "Latest jobs" |
| Sub-heading (`h3`) | `font-display text-base font-semibold` | A module name within the Modules section |
| Body | 15px (`text-[0.9375rem]`, set on `html`) | Values, table cells |
| Metadata / secondary | `text-muted` | Section descriptions, property names, timestamps in headers |
| Label | `text-label` (xs, medium, uppercase, wide tracking, muted) | Table column headers |
| Technical data | `code` (mono, 13px, inherits color) | Paths, addresses, IDs, versions, ports |
| Numbers | `tabular-nums` | Sizes, counts, dates, durations |
| Big stat | `font-display text-5xl font-semibold tabular-nums` (via `<Stat>`) | Dashboard counters |

**Rules**
- Each page has exactly one `h1`, and it names the object shown, never a generic label like "General info".
- On detail and edit pages, the `h1` is the object type in regular weight and `text-muted`, followed by the name in bold, with no colon: "Client **local**", "Edit client **local**", "Module **generic**", "Repository **local**".
- Show hierarchy with size and weight, not with boxes. To make something stand out, bold it or color it before reaching for a component.
- Do not use italics for empty or placeholder text. Use `<EmptyValue>` (muted "None").

### Spacing

All spacing uses the 4px Tailwind scale.

| Context | Value |
|---|---|
| Page padding | `px-6 py-6 lg:px-10` (in `layout/main.tsx`) |
| Content max width | `max-w-[100rem]` (1600px), centered. Content uses the full width below that |
| Page header → first tile | `pb-6` on `PageHeader` + the 16px gap (40px in total) |
| Between tiles | `gap-4` (16px), set by `Main` (`flex flex-col gap-4`). Grids of tiles also use `gap-4` |
| Tile padding | `p-5 lg:p-6` |
| Inside a section (heading → content) | `space-y-4` (set by `Section`) |
| Property list rows | `gap-y-2`, label column 7–11rem |
| Table cells | `px-3 py-2` |
| Groups inside a tile (e.g. modules) | Grid with `gap-x-10 gap-y-8` (or `gap-y-10` for forms) |

The default is tight. Add space only to separate groups of information that differ in meaning, never just to fill the screen.

### Shape

The radius is `0.25rem` on fields and buttons and `0.375rem` on boxes (tiles, modals). Borders are 1px. Depth and noise are 0. There are **no shadows and no gradients**.

## Layout

- **Sidebar** (`layout/sidebar.tsx`): fixed, `w-56`, `bg-base-200` with a right border. The active item has a semibold label, a 2px amber left marker and a light fill. Icons are Remix *line* variants in `text-muted`.
  The logo block shows the wordmark with the Relique version below it (`code text-xs text-muted`). App-level info like the version lives here, not in page headers.
- **Main** (`layout/main.tsx`): stacks the page content in a `flex flex-col gap-4` column, up to 1600px wide.
- **Page structure**: `PageHeader` (directly on the white page) → one or more `Section` tiles → `PropertyList` / tables / stats inside. Every block of content lives in a tile, including list tables (an untitled `<Section>`). The dashboard is the exception to the header rule: it shows only the breadcrumb (no title) and then its tiles. Place short tiles side by side with `grid md:grid-cols-2 gap-4` (or `lg:`) so wide screens are used.
- **Use the width.** Spread stats evenly with a column grid (`grid-cols-*`) instead of left-aligned flex rows. Lay forms out in two columns on `lg` (max `max-w-5xl`), not as a narrow single column.

## Components

| Component | Use it for | Don't |
|---|---|---|
| `PageHeader` | Breadcrumb, `h1`, meta line (status, ID, version), page actions | Put more than one primary button in `actions` |
| `Section` | A tile (`bg-base-200`, rounded, no border) with an `h2`, an optional description shown inline after a ·, an actions slot, then content. `variant="plain"` drops the tile, for rare nested use | Nest tiles inside tiles (use `h3` for sub-groups). Add borders or shadows to tiles |
| `PropertyList` | Key/value details (`<dl>`). Supports `mono` values, a loading skeleton, and empty → "None" | Use a `<table>` for key/value pairs |
| `ClientProperties`, `ModuleProperties`, `RepositoryProperties` | Standard property sets for domain objects | Duplicate their fields by hand |
| `Table` (via `*List` components) | Lists with filtering, sorting and pagination, placed inside a tile. Rows have a hairline separator and a translucent hover (it works on any background). Cells don't wrap; the table scrolls horizontally instead | Add zebra stripes, outer borders or colored headers |
| `Stat` | A big number + label. `tone` applies only when the value is > 0. `to` makes it a link | Add icons or boxes around stats |
| `Status` | A status dot + colored label (`variant="dot"` for dense table cells, with the label kept as a tooltip and for screen readers) | Use color without a label. Use badges for status |
| `YesNo` | Boolean values | Red/green "true"/"false" |
| `EmptyValue` | A muted "None" (or a custom label) for missing data | Leave cells blank or use italics |
| `TextPlaceholder` | Skeleton line while loading | Show spinners for whole sections |

### Actions

- **Primary** (`btn btn-sm btn-primary`): amber fill with dark text. At most one per view (Create, Save).
- **Secondary** (`btn btn-sm`): every other action (Edit, Cancel, Reload).
- **Destructive in a toolbar** (`btn btn-sm btn-ghost text-error`). The solid `btn-error` is reserved for the confirm button inside the confirmation dialog.
- **Links** (`text-link`): amber-brown text with an underline on hover. Use them for navigation to related objects (client names, job IDs).
- Icons inside buttons are optional and use line variants (`ri-*-line`).

### Forms

- Use two columns on `lg` (`grid lg:grid-cols-2 gap-x-6`, max `max-w-5xl`). Pair short related fields inside a column with `grid grid-cols-2 gap-x-6`. Repeated sub-forms (modules) go in a 2-column grid.
- Disabled fields add `disabled:bg-base-300/60` so they stay visible on a tile.
- Use the daisyUI `fieldset` + `label` + `input w-full`. Technical fields (paths, addresses, users, ports) add `font-mono`.
- Show validation errors inline below the field as `text-error`, and only once the user has edited the field. Add `input-error` to the field.
- Put help text below the field using the muted `label` class.

### Feedback

- Toasts (`Utils.notify`) use react-toastify's light theme, with a semibold title and a muted body.
- Use confirmation dialogs (`modal`) for destructive actions. The title states the action and the object, and the body states the consequence.

## Don't

- Tiles inside tiles, bordered cards, or content placed loose on the page outside a tile (except the page header).
- Shadows, gradients, glows or decorative icons (e.g. `text-6xl` icons beside stats).
- Badges for names (clients, modules, variants). Use plain text, joined with commas for lists.
- Raw Tailwind palette colors or hex values in components. Add a token instead.
- Status shown only through color.
- Amber used for status or for body text.
