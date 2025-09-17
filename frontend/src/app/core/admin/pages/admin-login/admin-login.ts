import { Component, inject, signal } from "@angular/core";
import { Router } from "@angular/router";
import { ReactiveFormsModule, FormBuilder, Validators } from "@angular/forms";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatButtonModule } from "@angular/material/button";
import { RouterModule } from "@angular/router";
import { AdminAuth } from "../../services/admin-auth";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  selector: "app-admin-login",
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    RouterModule,
  ],
  templateUrl: "./admin-login.html",
  styleUrl: "./admin-login.scss",
})
export class AdminLogin {
  private readonly adminAuth = inject(AdminAuth);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  readonly form = this.fb.group({
    username: ["", Validators.required],
    password: ["", Validators.required],
  });

  readonly error = signal<string | null>(null);
  readonly isSubmitting = signal(false);

  login(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.isSubmitting.set(true);
    this.error.set(null);

    const req = {
      username: this.form.value.username!,
      password: this.form.value.password!,
    };

    this.adminAuth.login(req).subscribe({
      next: () => {
        this.isSubmitting.set(false);
        this.router.navigate(["/admin"]);
      },
      error: (err) => this.handleError(err),
    });
  }

  private handleError(error: HttpErrorResponse): void {
    this.isSubmitting.set(false);
    let message = "Failed to login";
    if (error.status === 0) {
      message = "Network error. Please check your connection.";
    } else if (error.status === 401) {
      message = "Invalid admin credentials.";
    } else if (error.status === 403) {
      message = "Access forbidden.";
    } else if (error.status === 500) {
      message = "Server error. Please try again later.";
    }
    this.error.set(message);
  }
}