import {
  Component,
  effect,
  inject,
  input,
  output,
  signal,
} from "@angular/core";
import { MatInputModule } from "@angular/material/input";
import { EditEntryRequest, Entry } from "../../models/entry";
import { MatFormFieldModule } from "@angular/material/form-field";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";

@Component({
  selector: "app-entry-form",
  imports: [
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    ReactiveFormsModule,
  ],
  templateUrl: "./entry-form.html",
  styleUrl: "./entry-form.scss",
})
export class EntryForm {
  private readonly formBuilder = inject(FormBuilder);

  readonly disabled = input(false);
  readonly entry = input<Entry | null>(null);

  readonly errorMessage = signal<string | null>(null);

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
      validators: [Validators.required],
    }),
    note: this.formBuilder.control("", { nonNullable: true }),
  });

  readonly entrySubmit = output<EditEntryRequest>();

  handleSubmit(): void {
    if (this.form.invalid) {
      this.errorMessage.set("Please fill in all required fields.");
      return;
    }

    this.errorMessage.set(null);

    const { feeling, note } = this.form.value;
    const payload = { feeling: feeling!.trim(), note: note?.trim() };
    this.entrySubmit.emit(payload);
  }
}
