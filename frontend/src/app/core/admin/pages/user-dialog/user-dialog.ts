import { Component, inject, signal, OnInit } from "@angular/core";
import { FormBuilder, Validators, ReactiveFormsModule } from "@angular/forms";
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatButtonModule } from "@angular/material/button";
import { AdminAuth, User, CreateUserRequest, UpdateUserRequest } from "../../services/admin-auth";
import { HttpErrorResponse } from "@angular/common/http";

export interface UserDialogData {
  mode: "create" | "edit";
  user?: User;
}

@Component({
  selector: "app-user-dialog",
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
  ],
  templateUrl: "./user-dialog.html",
  styleUrl: "./user-dialog.scss",
})
export class UserDialogComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly adminAuth = inject(AdminAuth);
  private readonly dialogRef = inject(MatDialogRef<UserDialogComponent>);
  readonly data = inject<UserDialogData>(MAT_DIALOG_DATA);

  readonly error = signal<string | null>(null);
  readonly isSubmitting = signal(false);

  readonly form = this.fb.group({
    name: ["", Validators.required],
    email: ["", [Validators.required, Validators.email]],
    password: ["", this.data.mode === "create" ? [Validators.required, Validators.minLength(6)] : [Validators.minLength(6)]],
  });

  ngOnInit() {
    if (this.data.mode === "edit" && this.data.user) {
      this.form.patchValue({
        name: this.data.user.name,
        email: this.data.user.email,
        password: "", // Leave empty for edit mode
      });
    }
  }

  save() {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.isSubmitting.set(true);
    this.error.set(null);

    if (this.data.mode === "create") {
      this.createUser();
    } else {
      this.updateUser();
    }
  }

  private createUser() {
    const req: CreateUserRequest = {
      name: this.form.value.name!,
      email: this.form.value.email!,
      password: this.form.value.password!,
    };

    this.adminAuth.createUser(req).subscribe({
      next: () => {
        this.isSubmitting.set(false);
        this.dialogRef.close(true);
      },
      error: (err) => this.handleError(err),
    });
  }

  private updateUser() {
    if (!this.data.user) return;

    const req: UpdateUserRequest = {
      name: this.form.value.name!,
      email: this.form.value.email!,
    };

    // Only include password if it's not empty
    if (this.form.value.password) {
      req.password = this.form.value.password;
    }

    this.adminAuth.updateUser(this.data.user.id, req).subscribe({
      next: () => {
        this.isSubmitting.set(false);
        this.dialogRef.close(true);
      },
      error: (err) => this.handleError(err),
    });
  }

  private handleError(error: HttpErrorResponse): void {
    this.isSubmitting.set(false);
    let message = "Failed to save user";
    if (error.status === 0) {
      message = "Network error. Please check your connection.";
    } else if (error.status === 400) {
      message = "Invalid user data. Please check your input.";
    } else if (error.status === 409) {
      message = "A user with this email already exists.";
    } else if (error.status === 500) {
      message = "Server error. Please try again later.";
    }
    this.error.set(message);
  }

  cancel() {
    this.dialogRef.close(false);
  }
}