import {
  Component,
  effect,
  inject,
  input,
  output,
  signal,
} from "@angular/core";
import { EditMoodEntryRequest, MoodEntry } from "../models/mood-entry";
import {
  AbstractControl,
  FormBuilder,
  ReactiveFormsModule,
  ValidationErrors,
} from "@angular/forms";

interface EmojiOption {
  readonly value: string;
  readonly label: string;
}

const EMOJI_OPTIONS: readonly EmojiOption[] = [
  { value: "😀", label: "Joyful" },
  { value: "🙂", label: "Good" },
  { value: "😌", label: "Calm" },
  { value: "🥰", label: "Loved" },
  { value: "🤔", label: "Reflective" },
  { value: "😴", label: "Tired" },
  { value: "🥲", label: "Tender" },
  { value: "😟", label: "Anxious" },
  { value: "😔", label: "Low" },
  { value: "😤", label: "Frustrated" },
  { value: "😡", label: "Angry" },
  { value: "😵‍💫", label: "Overwhelmed" },
];

@Component({
  selector: "app-mood-entry-form",
  imports: [ReactiveFormsModule],
  templateUrl: "./mood-entry-form.html",
  styleUrl: "./mood-entry-form.scss",
})
export class MoodEntryForm {
  private readonly formBuilder = inject(FormBuilder);

  readonly disabled = input(false);
  readonly entry = input<MoodEntry | null>(null);
  readonly submitLabel = input("Save entry");
  readonly noteExpanded = signal(false);
  readonly emojiOptions = EMOJI_OPTIONS;

  constructor() {
    effect(() => {
      const entry = this.entry();
      if (entry) {
        this.form.patchValue({
          feeling: entry.feeling,
          emoji: entry.emoji ?? "",
          note: entry.note ?? "",
        });
        this.noteExpanded.set(!!entry.note?.trim());
        return;
      }

      this.form.reset({
        feeling: "",
        emoji: "",
        note: "",
      });
      this.noteExpanded.set(false);
    });
  }

  readonly form = this.formBuilder.group({
    feeling: this.formBuilder.control("", {
      nonNullable: true,
      validators: [trimmedRequired],
    }),
    emoji: this.formBuilder.control("", { nonNullable: true }),
    note: this.formBuilder.control("", { nonNullable: true }),
  });

  readonly entrySubmit = output<EditMoodEntryRequest>();

  toggleNote(): void {
    this.noteExpanded.update((expanded) => !expanded);
  }

  selectEmoji(emoji: string): void {
    this.form.controls.emoji.setValue(emoji);
  }

  clearEmoji(): void {
    this.form.controls.emoji.setValue("");
  }

  isEmojiSelected(emoji: string): boolean {
    return this.form.controls.emoji.value === emoji;
  }

  handleSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { feeling, emoji, note } = this.form.value;
    const trimmedEmoji = emoji?.trim();
    const trimmedNote = note?.trim();
    const payload: EditMoodEntryRequest = {
      feeling: feeling!.trim(),
      ...(trimmedEmoji ? { emoji: trimmedEmoji } : {}),
      ...(trimmedNote ? { note: trimmedNote } : {}),
    };
    this.entrySubmit.emit(payload);
  }
}

function trimmedRequired(
  control: AbstractControl<string | null>,
): ValidationErrors | null {
  return control.value?.trim() ? null : { required: true };
}
