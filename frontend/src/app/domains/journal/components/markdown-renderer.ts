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

const MOOD_RECORD_LINK_PATTERN = /^\[\[mood:(\d+)(?:\|([^\]]+))?\]\]/;
let markedConfigured = false;

const moodRecordLinkExtension = {
  name: "moodRecordLink",
  level: "inline" as const,
  start(src: string): number | undefined {
    const index = src.indexOf("[[mood:");
    return index >= 0 ? index : undefined;
  },
  tokenizer(src: string):
    | {
        type: string;
        raw: string;
        moodRecordId: number;
        label: string;
      }
    | undefined {
    const match = MOOD_RECORD_LINK_PATTERN.exec(src);
    if (!match) {
      return undefined;
    }

    const moodRecordId = Number(match[1]);
    const label = match[2]?.trim() || `Mood record #${match[1]}`;

    return {
      type: "moodRecordLink",
      raw: match[0],
      moodRecordId,
      label,
    };
  },
  renderer(token: { moodRecordId: number; label: string }): string {
    return `<a href="/mood-records/${token.moodRecordId}" data-mood-record-link="true">${escapeHtml(token.label)}</a>`;
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
  host: {
    "(click)": "handleClick($event)",
  },
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

    const moodRecordID = parseMoodRecordID(anchor.getAttribute("href"));
    if (!moodRecordID) {
      return;
    }

    event.preventDefault();
    this.router.navigate(["/mood-records", moodRecordID]);
  }
}

function parseMoodRecordID(href: string | null): number | null {
  if (!href) {
    return null;
  }

  try {
    const baseUrl =
      typeof window === "undefined"
        ? "http://localhost"
        : window.location.origin;
    const url = new URL(href, baseUrl);
    const match = /^\/mood-records\/(\d+)\/?$/.exec(url.pathname);
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
    extensions: [moodRecordLinkExtension],
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
