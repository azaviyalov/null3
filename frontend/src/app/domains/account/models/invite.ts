export interface InviteValidationResponse {
  expires_at: string;
}

export interface InviteRegistrationRequest {
  login: string;
  email: string;
  password: string;
}
