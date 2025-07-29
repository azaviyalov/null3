import {
  Component,
  inject,
  signal,
  AfterViewInit,
  ElementRef,
  viewChild,
} from "@angular/core";
import { Router } from "@angular/router";
import { ReactiveFormsModule, FormBuilder, Validators } from "@angular/forms";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatButtonModule } from "@angular/material/button";
import { Auth } from "../../services/auth";
import { LoginRequest } from "../../models/login";
import { HttpErrorResponse } from "@angular/common/http";

const HOME_ROUTE = "";

@Component({
  selector: "app-login",
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
  ],
  templateUrl: "./login.html",
  styleUrl: "./login.scss",
})
export class Login implements AfterViewInit {
  private readonly auth = inject(Auth);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  readonly loginInput = viewChild<ElementRef<HTMLInputElement>>("loginInput");
  readonly passwordInput =
    viewChild<ElementRef<HTMLInputElement>>("passwordInput");

  readonly form = this.fb.group({
    login: ["", Validators.required],
    password: ["", Validators.required],
  });

  readonly error = signal<string | null>(null);
  readonly isSubmitting = signal(false);

  ngAfterViewInit(): void {
    // Update validity, so you don't have to click login button twice
    this.form.controls["login"].updateValueAndValidity();
    this.form.controls["password"].updateValueAndValidity();
  }

  login(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.isSubmitting.set(true);
    this.error.set(null);

    const req: LoginRequest = {
      login: this.form.value.login!,
      password: this.form.value.password!,
    };

    this.auth.login(req).subscribe({
      next: () => {
        this.isSubmitting.set(false);
        this.router.navigate([HOME_ROUTE]);
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
      message = "Unauthorized: Incorrect login credentials.";
    } else if (error.status === 500) {
      message = "Server error. Please try again later.";
    }
    this.error.set(message);
  }
}
