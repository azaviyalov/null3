import { Injectable, inject } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Observable } from "rxjs";
import { environment } from "../../../../environments/environment";
import {
  InviteRegistrationRequest,
  InviteValidationResponse,
} from "../models/invite";
import {
  ForgotPasswordRequest,
  ForgotPasswordResponse,
  MessageResponse,
  ResetPasswordRequest,
} from "../models/password";
import { UserResponse } from "../../session/models/user";

@Injectable({ providedIn: "root" })
export class AccountApi {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${environment.apiUrl}/auth`;

  registerWithInvite(
    token: string,
    data: InviteRegistrationRequest,
  ): Observable<UserResponse> {
    return this.http.post<UserResponse>(
      `${this.baseUrl}/invites/${token}/register`,
      data,
    );
  }

  getInvite(token: string): Observable<InviteValidationResponse> {
    return this.http.get<InviteValidationResponse>(
      `${this.baseUrl}/invites/${token}`,
    );
  }

  requestPasswordReset(
    data: ForgotPasswordRequest,
  ): Observable<ForgotPasswordResponse> {
    return this.http.post<ForgotPasswordResponse>(
      `${this.baseUrl}/forgot-password`,
      data,
    );
  }

  resetPassword(data: ResetPasswordRequest): Observable<MessageResponse> {
    return this.http.post<MessageResponse>(
      `${this.baseUrl}/reset-password`,
      data,
    );
  }
}
