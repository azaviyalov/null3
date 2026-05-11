import { Component, effect, inject, input, output } from "@angular/core";
import { EditEntryRequest, Entry } from "../../models/entry";
import {
  AbstractControl,
  FormBuilder,
  ReactiveFormsModule,
  ValidationErrors,
} from "@angular/forms";

@Component({
  selector: "app-entry-form",
  imports: [ReactiveFormsModule],
  templateUrl: "./entry-form.html",
  styleUrl: "./entry-form.scss",
})
export class EntryForm {
  private readonly formBuilder = inject(FormBuilder);

  readonly disabled = input(false);
  readonly entry = input<Entry | null>(null);
  readonly submitLabel = input("Save entry");

  constructor() {
    effect(() => {
      const entry = this.entry();
      if (entry) {
        this.form.patchValue({
          feeling: entry.feeling,
          note: entry.note,
        });
      }
    });
  }

  readonly form = this.formBuilder.group({
    feeling: this.formBuilder.control("", {
      nonNullable: true,
      validators: [trimmedRequired],
    }),
    note: this.formBuilder.control("", { nonNullable: true }),
  });

  readonly entrySubmit = output<EditEntryRequest>();

  handleSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { feeling, note } = this.form.value;
    const payload = { feeling: feeling!.trim(), note: note?.trim() };
    this.entrySubmit.emit(payload);
  }
}

function trimmedRequired(
  control: AbstractControl<string | null>,
): ValidationErrors | null {
  return control.value?.trim() ? null : { required: true };
}
