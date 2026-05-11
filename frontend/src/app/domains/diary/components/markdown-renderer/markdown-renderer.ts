import { CommonModule } from "@angular/common";
import {
  Component,
  ViewEncapsulation,
  computed,
  inject,
  input,
} from "@angular/core";
import { Router } from "@angular/router";
import { marked } from "marked";

const MOOD_ENTRY_LINK_PATTERN = /^\[\[mood:(\d+)(?:\|([^\]]+))?\]\]/;
let markedConfigured = false;

const moodEntryLinkExtension = {
  name: "moodEntryLink",
  level: "inline" as const,
  start(src: string): number | undefined {
    const index = src.indexOf("[[mood:");
    return index >= 0 ? index : undefined;
  },
  tokenizer(src: string):
    | {
        type: string;
        raw: string;
        moodEntryId: number;
        label: string;
      }
    | undefined {
    const match = MOOD_ENTRY_LINK_PATTERN.exec(src);
    if (!match) {
      return undefined;
    }

    const moodEntryId = Number(match[1]);
    const label = match[2]?.trim() || `Mood Entry #${match[1]}`;

    return {
      type: "moodEntryLink",
      raw: match[0],
      moodEntryId,
      label,
    };
  },
  renderer(token: { moodEntryId: number; label: string }): string {
    return `<a href="/mood/entries/${token.moodEntryId}" data-mood-entry-link="true">${escapeHtml(token.label)}</a>`;
  },
};

configureMarked();

@Component({
  selector: "app-markdown-renderer",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./markdown-renderer.html",
  styleUrl: "./markdown-renderer.scss",
  encapsulation: ViewEncapsulation.None,
})
export class MarkdownRenderer {
  private readonly router = inject(Router);

  readonly markdown = input("");
  readonly emptyMessage = input("Nothing written yet.");

  readonly html = computed(() => {
    const source = this.markdown().trim();
    if (!source) {
      return "";
    }

    return marked.parse(source) as string;
  });

  handleClick(event: MouseEvent): void {
    const target = event.target;
    if (!(target instanceof Element)) {
      return;
    }

    const anchor = target.closest("a");
    if (!(anchor instanceof HTMLAnchorElement)) {
      return;
    }

    const moodEntryID = parseMoodEntryID(anchor.getAttribute("href"));
    if (!moodEntryID) {
      return;
    }

    event.preventDefault();
    this.router.navigate(["/mood/entries", moodEntryID]);
  }
}

function parseMoodEntryID(href: string | null): number | null {
  if (!href) {
    return null;
  }

  try {
    const baseUrl =
      typeof window === "undefined"
        ? "http://localhost"
        : window.location.origin;
    const url = new URL(href, baseUrl);
    const match = /^\/mood\/entries\/(\d+)\/?$/.exec(url.pathname);
    if (!match) {
      return null;
    }

    return Number(match[1]);
  } catch {
    return null;
  }
}

function configureMarked(): void {
  if (markedConfigured) {
    return;
  }

  marked.use({
    breaks: true,
    gfm: true,
    extensions: [moodEntryLinkExtension],
  });
  markedConfigured = true;
}

function escapeHtml(value: string): string {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}
