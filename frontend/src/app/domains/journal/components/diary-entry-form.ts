import {
  Component,
  computed,
  effect,
  inject,
  input,
  output,
} from "@angular/core";
import { toSignal } from "@angular/core/rxjs-interop";
import {
  AbstractControl,
  FormBuilder,
  ReactiveFormsModule,
  ValidationErrors,
  Validators,
} from "@angular/forms";
import { startWith } from "rxjs";
import { MoodRecordApi } from "../services/mood-record-api";
import { DiaryEntry, EditDiaryEntryRequest } from "../models/diary-entry";
import { MoodRecord } from "../models/mood-record";
import { MarkdownRenderer } from "./markdown-renderer";

@Component({
  selector: "app-diary-entry-form",
  standalone: true,
  imports: [ReactiveFormsModule, MarkdownRenderer],
  templateUrl: "./diary-entry-form.html",
  styleUrl: "./diary-entry-form.scss",
})
export class DiaryEntryForm {
  private readonly formBuilder = inject(FormBuilder);
  private readonly moodRecordApi = inject(MoodRecordApi);

  readonly disabled = input(false);
  readonly entry = input<DiaryEntry | null>(null);
  readonly submitLabel = input("Save entry");

  readonly recentMoodRecordsPage = toSignal(this.moodRecordApi.getPaged(8), {
    initialValue: null,
  });
  readonly recentMoodRecords = computed(
    () => this.recentMoodRecordsPage()?.items ?? [],
  );

  constructor() {
    effect(() => {
      const entry = this.entry();
      if (entry) {
        this.form.patchValue({
          title: entry.title ?? "",
          occurredAt: toDatetimeLocalValue(entry.occurredAt),
          markdown: entry.markdown,
        });
        return;
      }

      this.form.reset({
        title: "",
        occurredAt: toDatetimeLocalValue(new Date()),
        markdown: "",
      });
    });
  }

  readonly form = this.formBuilder.group({
    title: this.formBuilder.control("", { nonNullable: true }),
    occurredAt: this.formBuilder.control("", {
      nonNullable: true,
      validators: [Validators.required],
    }),
    markdown: this.formBuilder.control("", {
      nonNullable: true,
      validators: [trimmedRequired],
    }),
  });

  readonly markdownValue = toSignal(
    this.form.controls.markdown.valueChanges.pipe(
      startWith(this.form.controls.markdown.value),
    ),
    { initialValue: this.form.controls.markdown.value },
  );

  readonly entrySubmit = output<EditDiaryEntryRequest>();

  insertMoodRecordLink(entry: MoodRecord, textarea: HTMLTextAreaElement): void {
    const control = this.form.controls.markdown;
    const value = control.value;
    const selectionStart = textarea.selectionStart ?? value.length;
    const selectionEnd = textarea.selectionEnd ?? value.length;
    const selectedText = value.slice(selectionStart, selectionEnd).trim();
    const label = normalizeMoodRecordLinkLabel(
      selectedText || buildMoodRecordLinkLabel(entry),
    );
    const link = `[[mood:${entry.id}|${label}]]`;

    const prefix =
      selectionStart > 0 && !/\s/.test(value.charAt(selectionStart - 1))
        ? " "
        : "";
    const suffix =
      selectionEnd < value.length && !/\s/.test(value.charAt(selectionEnd))
        ? " "
        : "";

    const nextValue =
      value.slice(0, selectionStart) +
      prefix +
      link +
      suffix +
      value.slice(selectionEnd);
    const caretPosition = selectionStart + prefix.length + link.length;

    control.setValue(nextValue);
    control.markAsDirty();

    queueMicrotask(() => {
      textarea.focus();
      textarea.setSelectionRange(caretPosition, caretPosition);
    });
  }

  handleSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { title, occurredAt, markdown } = this.form.getRawValue();
    const trimmedTitle = title.trim();
    const trimmedMarkdown = markdown.trim();

    this.entrySubmit.emit({
      ...(trimmedTitle ? { title: trimmedTitle } : {}),
      markdown: trimmedMarkdown,
      occurred_at: new Date(occurredAt).toISOString(),
    });
  }
}

function buildMoodRecordLinkLabel(entry: MoodRecord): string {
  const dateLabel = new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  }).format(entry.createdAt);
  const moodLabel = entry.note || entry.feeling;
  return [entry.emoji, moodLabel, dateLabel].filter(Boolean).join(" ");
}

function normalizeMoodRecordLinkLabel(label: string): string {
  return label
    .replaceAll("|", " ")
    .replaceAll("]]", " ")
    .replace(/\s+/g, " ")
    .trim();
}

function toDatetimeLocalValue(date: Date): string {
  const year = `${date.getFullYear()}`;
  const month = `${date.getMonth() + 1}`.padStart(2, "0");
  const day = `${date.getDate()}`.padStart(2, "0");
  const hours = `${date.getHours()}`.padStart(2, "0");
  const minutes = `${date.getMinutes()}`.padStart(2, "0");
  return `${year}-${month}-${day}T${hours}:${minutes}`;
}

function trimmedRequired(
  control: AbstractControl<string | null>,
): ValidationErrors | null {
  return control.value?.trim() ? null : { required: true };
}
