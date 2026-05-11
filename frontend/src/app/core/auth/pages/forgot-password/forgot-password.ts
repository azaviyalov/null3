import { Component, inject, signal } from "@angular/core";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { RouterModule } from "@angular/router";
import { Auth } from "../../services/auth";
import { ForgotPasswordRequest } from "../../models/password";
import { ForgotPasswordResponse } from "../../models/password";

@Component({
  selector: "app-forgot-password",
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    RouterModule,
  ],
  templateUrl: "./forgot-password.html",
  styleUrl: "./forgot-password.scss",
})
export class ForgotPassword {
  private readonly auth = inject(Auth);
  private readonly fb = inject(FormBuilder);

  readonly form = this.fb.group({
    email: ["", [Validators.required, Validators.email]],
  });

  readonly error = signal<string | null>(null);
  readonly isSubmitting = signal(false);
  readonly response = signal<ForgotPasswordResponse | null>(null);

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.error.set(null);
    this.response.set(null);
    this.isSubmitting.set(true);

    const req: ForgotPasswordRequest = {
      email: this.form.value.email!,
    };

    this.auth.requestPasswordReset(req).subscribe({
      next: (response) => {
        this.response.set(response);
        this.isSubmitting.set(false);
      },
      error: (error) => this.handleError(error),
    });
  }

  private handleError(error: HttpErrorResponse): void {
    this.isSubmitting.set(false);
    if (error.status === 0) {
      this.error.set("Network error. Please check your connection.");
      return;
    }
    this.error.set("Failed to request a password reset.");
  }
}
