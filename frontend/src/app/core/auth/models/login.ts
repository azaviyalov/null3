import { UserResponse } from "./user";

export interface LoginRequest {
  login: string;
  password: string;
}

export interface LoginResponse extends UserResponse {
  token: string;
}
