import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { environment } from '../../environments/environment';
import { catchError, map, tap } from 'rxjs/operators';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export interface AuthResponse {
  status: string;
  message?: string;
  user_id?: string;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = environment.apiUrl;
  private authenticated = false;
  private userId: string | null = null;
  private cachedProfile: any = null;
  private cachedAuthCheck: boolean | null = null;
  private authCheckTimestamp: number = 0;
  private readonly httpOptions = { withCredentials: true };
  private readonly CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

  constructor(private http: HttpClient) {}

  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(
        `${this.apiUrl}/auth/login`,
        credentials,
        this.httpOptions,
      )
      .pipe(
        tap((response) => {
          if (response.status === 'success') {
            this.authenticated = true;
            this.cachedProfile = null;
          }
        }),
        map((response) => {
          if (response.status === 'success') {
            return response;
          }
          throw new Error('Login failed');
        }),
      );
  }

  register(userData: RegisterRequest): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(
        `${this.apiUrl}/auth/register`,
        userData,
        this.httpOptions,
      )
      .pipe(
        tap((response) => {
          if (response.status === 'success') {
            this.authenticated = true;
            this.cachedProfile = null;
          }
        }),
        map((response) => {
          if (response.status === 'success') {
            return response;
          }
          throw new Error('Registration failed');
        }),
      );
  }

  logout(): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(`${this.apiUrl}/auth/logout`, {}, this.httpOptions)
      .pipe(
        tap(() => this.clearSession()),
        catchError((error) => {
          this.clearSession();
          throw error;
        }),
      );
  }

  /** Validates the HttpOnly session cookie with the auth service. */
  checkAuth(): Observable<boolean> {
    const now = Date.now();

    // Return cached result if still valid
    if (
      this.cachedAuthCheck !== null &&
      now - this.authCheckTimestamp < this.CACHE_DURATION
    ) {
      return of(this.cachedAuthCheck);
    }

    return this.http
      .get<AuthResponse>(`${this.apiUrl}/auth/me`, this.httpOptions)
      .pipe(
        map((response) => {
          const ok = response.status === 'success' && !!response.user_id;
          this.authenticated = ok;
          this.userId = ok ? (response.user_id ?? null) : null;
          this.cachedAuthCheck = ok;
          this.authCheckTimestamp = now;
          return ok;
        }),
        catchError(() => {
          this.clearSession();
          this.cachedAuthCheck = false;
          this.authCheckTimestamp = now;
          return of(false);
        }),
      );
  }

  isLoggedIn(): boolean {
    return this.authenticated;
  }

  getUserId(): string | null {
    return this.userId;
  }

  getProfile(): Observable<any> {
    if (this.cachedProfile) {
      return of(this.cachedProfile);
    }

    return this.http.get<any>(`${this.apiUrl}/users/me`, this.httpOptions).pipe(
      tap((profile) => {
        if (profile?.id) {
          this.userId = profile.id;
          this.authenticated = true;
          this.cachedProfile = profile;
        }
      }),
      catchError((error) => {
        console.error('Error fetching profile:', error);
        throw error;
      }),
    );
  }

  clearSession(): void {
    this.authenticated = false;
    this.userId = null;
    this.cachedProfile = null;
    this.cachedAuthCheck = null;
    this.authCheckTimestamp = 0;
  }

  isBrowser(): boolean {
    return typeof window !== 'undefined';
  }
}
